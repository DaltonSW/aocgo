package session

import (
	"errors"
	"os"
	"path/filepath"
)

func GetSessionToken() (string, error) {
	sessionToken, err := getTokenFromFile("")
	if sessionToken != "" {
		return sessionToken, err
	}

	sessionToken, err = getTokenFromEnv()
	if sessionToken != "" {
		return sessionToken, err
	}

	return "", errors.New("Unable to load AoC session token from file or environment variable")
}

// Making this a separate function so it's testable
func getTokenFromEnv() (string, error) {
	token := os.Getenv("AOC_SESSION_TOKEN")
	if token == "" {
		return "", errors.New("Couldn't load session token from environment variable!")
	}
	return token, nil
}

func getTokenFromFile(path string) (string, error) {
	var filePath string
	if path == "" {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		filePath = filepath.Join(userHomeDir, ".config", "aocutil", "session.token")

	} else {
		filePath = path
	}
	file, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(file), nil
}
