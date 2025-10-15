package session

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddSetAndRemoveUser(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	active, err := AddUser("primary", "token-1", true)
	if err != nil {
		t.Fatalf("AddUser returned error: %v", err)
	}
	if active == nil || active.Label != "primary" || active.Token != "token-1" {
		t.Fatalf("AddUser returned unexpected active user: %#v", active)
	}

	list, err := ListUsers()
	if err != nil {
		t.Fatalf("ListUsers returned error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 user, got %d", len(list))
	}
	if !list[0].Active {
		t.Fatalf("expected newly added user to be active")
	}

	_, err = AddUser("secondary", "token-2", false)
	if err != nil {
		t.Fatalf("AddUser(secondary) returned error: %v", err)
	}

	active, err = SetActiveUser("secondary")
	if err != nil {
		t.Fatalf("SetActiveUser returned error: %v", err)
	}
	if active.Label != "secondary" {
		t.Fatalf("expected secondary to be active, got %s", active.Label)
	}

	active, err = RemoveUser("secondary")
	if err != nil {
		t.Fatalf("RemoveUser returned error: %v", err)
	}
	if active == nil || active.Label != "primary" {
		t.Fatalf("expected primary to remain active after removing secondary, got %#v", active)
	}

	active, err = RemoveUser("primary")
	if err != nil {
		t.Fatalf("RemoveUser last user returned error: %v", err)
	}
	if active != nil {
		t.Fatalf("expected active to be nil after removing last user, got %#v", active)
	}
}

func TestGetActiveUserFromEnv(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("AOC_SESSION_TOKEN", "env-token")
	t.Cleanup(func() {
		os.Unsetenv("AOC_SESSION_TOKEN")
	})

	active, err := GetActiveUser(false)
	if err != nil {
		t.Fatalf("GetActiveUser returned error: %v", err)
	}
	if active.Label != envUserLabel {
		t.Fatalf("expected label %q, got %q", envUserLabel, active.Label)
	}
	if active.Token != "env-token" {
		t.Fatalf("expected token env-token, got %s", active.Token)
	}
}

func TestLegacyTokenMigration(t *testing.T) {
	configRoot := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configRoot)

	legacyDir := filepath.Join(configRoot, "aocgo")
	if err := os.MkdirAll(legacyDir, 0o700); err != nil {
		t.Fatalf("unable to create legacy dir: %v", err)
	}

	legacyPath := filepath.Join(legacyDir, legacyTokenFileName)
	if err := os.WriteFile(legacyPath, []byte("legacy-token"), 0o600); err != nil {
		t.Fatalf("unable to write legacy token: %v", err)
	}

	active, err := GetActiveUser(false)
	if err != nil {
		t.Fatalf("GetActiveUser returned error after migration: %v", err)
	}
	if active.Label != "default" {
		t.Fatalf("expected default label after migration, got %s", active.Label)
	}
	if active.Token != "legacy-token" {
		t.Fatalf("expected token legacy-token, got %s", active.Token)
	}

	// Ensure data persisted in new store.
	list, err := ListUsers()
	if err != nil {
		t.Fatalf("ListUsers after migration returned error: %v", err)
	}
	if len(list) != 1 || list[0].Label != "default" {
		t.Fatalf("expected one migrated user named default, got %#v", list)
	}
}
