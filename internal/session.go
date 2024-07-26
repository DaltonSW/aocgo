package internal

import (
	"errors"
	"os"
	"path/filepath"
)

func GetSessionToken() (string, error) {
	sessionToken := getTokenFromFile("")
	if sessionToken != "" {
		return sessionToken, nil
	}

	sessionToken = os.Getenv("AOC_SESSION_TOKEN")
	if sessionToken != "" {
		return sessionToken, nil
	}

	return "", errors.New("Unable to load AoC session token from file or environment variable")
}

func getTokenFromFile(path string) string {
	var filePath string
	if path == "" {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		filePath = filepath.Join(userHomeDir, ".config", "aocutil", "session.token")

	} else {
		filePath = path
	}
	file, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return string(file)
}
