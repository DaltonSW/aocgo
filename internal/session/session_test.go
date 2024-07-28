package session

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetSessionTokenFromEnv(t *testing.T) {
	// Setup: Create a temporary home directory and backup environment variable
	backupEnv := os.Getenv("AOC_SESSION_TOKEN")
	defer os.Setenv("AOC_SESSION_TOKEN", backupEnv)

	var err error

	// Test: Get token from environment variable
	expectedEnvToken := "test_token_from_env"
	os.Setenv("AOC_SESSION_TOKEN", expectedEnvToken)
	token, err := getTokenFromEnv()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if token != expectedEnvToken {
		t.Fatalf("Expected token %v, got %v", expectedEnvToken, token)
	}
	os.Unsetenv("AOC_SESSION_TOKEN")

	// Test: No token available
	token, err = getTokenFromEnv()
	if token != "" {
		t.Fatalf("Expected empty token, got %v", token)
	}
}

func TestGetTokenFromFile(t *testing.T) {
	// Setup: Create a temporary config directory and session token file
	homeDir := t.TempDir()
	configDir := filepath.Join(homeDir, ".config", "aocutil")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Unable to create config directory: %v", err)
	}
	tokenFile := filepath.Join(configDir, "session.token")

	// Test: Read token from file
	expectedToken := "test_token"
	if err := os.WriteFile(tokenFile, []byte(expectedToken), 0644); err != nil {
		t.Fatalf("Unable to write session token file: %v", err)
	}
	token, err := getTokenFromFile(tokenFile)
	if token != expectedToken {
		t.Fatalf("Expected token %v, got %v", expectedToken, token)
	}
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Cleanup: Remove the session token file
	if err := os.Remove(tokenFile); err != nil {
		t.Fatalf("Unable to remove session token file: %v", err)
	}

	// Test: Token file does not exist
	token, err = getTokenFromFile(tokenFile)
	if token != "" {
		t.Fatalf("Expected empty token, got %v", token)
	}
}
