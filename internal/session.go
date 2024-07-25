package internal

import (
	"errors"
	"os"
	"path/filepath"
)

func GetSessionToken() (string, error) {
	sessionToken := getTokenFromFile()
	if sessionToken != "" {
		return sessionToken, nil
	}

	sessionToken = os.Getenv("AOC_SESSION_TOKEN")
	if sessionToken != "" {
		return sessionToken, nil
	}

	return "", errors.New("Unable to load AoC session token from file or environment variable")
}

func getTokenFromFile() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	filepath := filepath.Join(userHomeDir, ".config", "aocutil", "session.token")
	file, err := os.ReadFile(filepath)
	if err != nil {
		return ""
	}
	return string(file)
}
