package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go.dalton.dog/aocgo/internal/api"
	"go.dalton.dog/aocgo/internal/output"
)

const (
	usersFileName       = "users.json"
	legacyTokenFileName = "session.token"
	envUserLabel        = "environment"
)

// ActiveUser represents the currently selected user.
type ActiveUser struct {
	Label string
	Token string
}

// Summary provides hydrated details about a stored user, primarily for CLI use.
type Summary struct {
	Label  string
	Token  string
	Active bool
}

type userStore struct {
	Active string                `json:"active"`
	Users  map[string]storedUser `json:"users"`
}

type storedUser struct {
	Token string `json:"token"`
}

// AddUser registers a new user in the local store. When makeActive is true or when
// the store previously had no active user, the provided label becomes the active user.
func AddUser(label, token string, makeActive bool) (*ActiveUser, error) {
	label = strings.TrimSpace(label)
	token = strings.TrimSpace(token)

	if label == "" {
		return nil, errors.New("user label cannot be empty")
	}
	if token == "" {
		return nil, errors.New("session token cannot be empty")
	}

	store, err := loadStore()
	if err != nil {
		return nil, err
	}

	if store.Users == nil {
		store.Users = make(map[string]storedUser)
	}

	if _, exists := store.Users[label]; exists {
		return nil, fmt.Errorf("user %q already exists", label)
	}

	store.Users[label] = storedUser{Token: token}

	if store.Active == "" || makeActive {
		store.Active = label
	}

	if err := saveStore(store); err != nil {
		return nil, err
	}

	return getActiveFromStore(store, false)
}

// ListUsers returns all users with active-state annotations.
func ListUsers() ([]Summary, error) {
	store, err := loadStore()
	if err != nil {
		return nil, err
	}

	summaries := make([]Summary, 0, len(store.Users))
	for label, data := range store.Users {
		summaries = append(summaries, Summary{
			Label:  label,
			Token:  data.Token,
			Active: label == store.Active,
		})
	}

	sort.SliceStable(summaries, func(i, j int) bool {
		return strings.ToLower(summaries[i].Label) < strings.ToLower(summaries[j].Label)
	})

	return summaries, nil
}

// SetActiveUser switches the active user to the provided label.
func SetActiveUser(label string) (*ActiveUser, error) {
	label = strings.TrimSpace(label)
	if label == "" {
		return nil, errors.New("user label cannot be empty")
	}

	store, err := loadStore()
	if err != nil {
		return nil, err
	}

	if _, exists := store.Users[label]; !exists {
		return nil, fmt.Errorf("user %q not found", label)
	}

	store.Active = label

	if err := saveStore(store); err != nil {
		return nil, err
	}

	return getActiveFromStore(store, false)
}

// RemoveUser deletes the given user from the store. If the removed user was active,
// the next available user becomes active (alphabetical order). If no users remain,
// no active user is returned.
func RemoveUser(label string) (*ActiveUser, error) {
	label = strings.TrimSpace(label)
	if label == "" {
		return nil, errors.New("user label cannot be empty")
	}

	store, err := loadStore()
	if err != nil {
		return nil, err
	}

	if _, exists := store.Users[label]; !exists {
		return nil, fmt.Errorf("user %q not found", label)
	}

	delete(store.Users, label)

	if len(store.Users) == 0 {
		store.Active = ""
		if err := saveStore(store); err != nil {
			return nil, err
		}
		return nil, nil
	}

	if store.Active == label {
		labels := make([]string, 0, len(store.Users))
		for k := range store.Users {
			labels = append(labels, k)
		}
		sort.SliceStable(labels, func(i, j int) bool {
			return strings.ToLower(labels[i]) < strings.ToLower(labels[j])
		})
		store.Active = labels[0]
	}

	if err := saveStore(store); err != nil {
		return nil, err
	}

	return getActiveFromStore(store, false)
}

// GetActiveUser returns the currently active user credentials.
func GetActiveUser(healthLog bool) (*ActiveUser, error) {
	store, err := loadStore()
	if err != nil {
		return nil, err
	}

	active, err := getActiveFromStore(store, healthLog)
	if err == nil && active != nil {
		return active, nil
	}

	if token, err := getTokenFromEnv(); err == nil {
		if healthLog {
			output.Info("Found session token in environment variable.", "label", envUserLabel, "token", api.MaskSecret(token))
		}
		return &ActiveUser{
			Label: envUserLabel,
			Token: token,
		}, nil
	}

	return nil, errors.New("unable to determine an active user; add one via `aocli user add`")
}

// GetSessionToken maintains backwards compatibility for legacy callers that only need the token.
func GetSessionToken(healthLog bool) (string, error) {
	active, err := GetActiveUser(healthLog)
	if err != nil {
		return "", err
	}
	return active.Token, nil
}

// GetUser returns stored credentials for the given label.
func GetUser(label string) (*ActiveUser, error) {
	label = strings.TrimSpace(label)
	if label == "" {
		return nil, errors.New("user label cannot be empty")
	}

	store, err := loadStore()
	if err != nil {
		return nil, err
	}

	user, ok := store.Users[label]
	if !ok {
		return nil, fmt.Errorf("user %q not found", label)
	}

	return &ActiveUser{Label: label, Token: user.Token}, nil
}

func getActiveFromStore(store *userStore, healthLog bool) (*ActiveUser, error) {
	if store == nil || store.Users == nil || store.Active == "" {
		return nil, errors.New("no active user configured")
	}

	user, ok := store.Users[store.Active]
	if !ok || strings.TrimSpace(user.Token) == "" {
		return nil, errors.New("active user has no valid session token")
	}

	if healthLog {
		output.Info("Found active user in config file.", "label", store.Active, "token", api.MaskSecret(user.Token))
	}

	return &ActiveUser{
		Label: store.Active,
		Token: strings.TrimSpace(user.Token),
	}, nil
}

func loadStore() (*userStore, error) {
	storePath, err := ensureConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(storePath)
	if errors.Is(err, os.ErrNotExist) {
		store := &userStore{
			Users: make(map[string]storedUser),
		}
		if migrated, migErr := migrateLegacyToken(store); migErr != nil {
			return nil, migErr
		} else if migrated {
			if err := saveStore(store); err != nil {
				return nil, err
			}
		}
		return store, nil
	}
	if err != nil {
		return nil, err
	}

	store := &userStore{}
	if err := json.Unmarshal(data, store); err != nil {
		return nil, err
	}

	if store.Users == nil {
		store.Users = make(map[string]storedUser)
	}

	return store, nil
}

func saveStore(store *userStore) error {
	if store == nil {
		return errors.New("nil store provided")
	}

	storePath, err := ensureConfigFilePath()
	if err != nil {
		return err
	}

	storeData, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(storePath, storeData, 0o600)
}

func ensureConfigFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appConfigDir := filepath.Join(configDir, "aocgo")

	if err := os.MkdirAll(appConfigDir, 0o700); err != nil {
		return "", err
	}

	return filepath.Join(appConfigDir, usersFileName), nil
}

func migrateLegacyToken(store *userStore) (bool, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return false, err
	}

	legacyPath := filepath.Join(configDir, "aocgo", legacyTokenFileName)

	data, err := os.ReadFile(legacyPath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	token := strings.TrimSpace(string(data))
	if token == "" {
		return false, nil
	}

	if store.Users == nil {
		store.Users = make(map[string]storedUser)
	}

	const defaultLabel = "default"

	store.Users[defaultLabel] = storedUser{Token: token}
	store.Active = defaultLabel

	return true, nil
}

// Making this a separate function so it's testable
func getTokenFromEnv() (string, error) {
	token := os.Getenv("AOC_SESSION_TOKEN")
	if token == "" {
		return "", errors.New("couldn't load session token from environment variable")
	}
	return strings.TrimSpace(token), nil
}
