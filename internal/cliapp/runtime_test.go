package cliapp

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ikuaidev/ikuai-cli/internal/session"
)

func TestSyncSessionSessionWinsOverEnv(t *testing.T) {
	sf := filepath.Join(t.TempDir(), "config.json")
	t.Setenv("IKUAI_CLI_CONFIG_FILE", sf)
	t.Setenv("IKUAI_CLI_BASE_URL", "https://env-router")
	t.Setenv("IKUAI_CLI_TOKEN", "env-token")

	if err := session.SaveLogin("https://session-router", "session-token"); err != nil {
		t.Fatalf("SaveLogin() error = %v", err)
	}

	r := New(&bytes.Buffer{}, &bytes.Buffer{})
	if err := r.SyncSession(); err != nil {
		t.Fatalf("SyncSession() error = %v", err)
	}

	if r.Session.BaseURL != "https://session-router" {
		t.Errorf("BaseURL = %q, want session value", r.Session.BaseURL)
	}
	if r.Session.Token != "session-token" {
		t.Errorf("Token = %q, want session value", r.Session.Token)
	}
	if r.CredSource != "token" {
		t.Errorf("CredSource = %q, want %q", r.CredSource, "token")
	}
}

func TestSyncSessionEnvFallbackWhenNoSessionToken(t *testing.T) {
	sf := filepath.Join(t.TempDir(), "config.json")
	t.Setenv("IKUAI_CLI_CONFIG_FILE", sf)
	t.Setenv("IKUAI_CLI_BASE_URL", "https://env-router/")
	t.Setenv("IKUAI_CLI_TOKEN", "env-token")

	r := New(&bytes.Buffer{}, &bytes.Buffer{})
	if err := r.SyncSession(); err != nil {
		t.Fatalf("SyncSession() error = %v", err)
	}

	if r.Session.BaseURL != "https://env-router" {
		t.Errorf("BaseURL = %q, want env value (trailing slash trimmed)", r.Session.BaseURL)
	}
	if r.Session.Token != "env-token" {
		t.Errorf("Token = %q, want env value", r.Session.Token)
	}
	if r.CredSource != "env" {
		t.Errorf("CredSource = %q, want %q", r.CredSource, "env")
	}
	if r.APIClient == nil {
		t.Error("APIClient should not be nil in env mode")
	}
}

func TestSyncSessionEnvNotPersistedToDisk(t *testing.T) {
	sf := filepath.Join(t.TempDir(), "config.json")
	t.Setenv("IKUAI_CLI_CONFIG_FILE", sf)
	t.Setenv("IKUAI_CLI_BASE_URL", "https://env-router")
	t.Setenv("IKUAI_CLI_TOKEN", "env-token")

	r := New(&bytes.Buffer{}, &bytes.Buffer{})
	if err := r.SyncSession(); err != nil {
		t.Fatalf("SyncSession() error = %v", err)
	}

	// Session file should either not exist or be empty — env must never be persisted.
	data, err := os.ReadFile(sf) //nolint:gosec // test-only: path from t.TempDir()
	if err != nil {
		if !os.IsNotExist(err) {
			t.Fatalf("ReadFile() unexpected error = %v", err)
		}
		return // file doesn't exist — good
	}
	// File exists but should not contain env values.
	if bytes.Contains(data, []byte("env-router")) {
		t.Fatalf("session file contains env BaseURL — env was persisted to disk: %s", data)
	}
	if bytes.Contains(data, []byte("env-token")) {
		t.Fatalf("session file contains env Token — env was persisted to disk: %s", data)
	}
}

func TestSyncSessionNeitherSessionNorEnv(t *testing.T) {
	t.Setenv("IKUAI_CLI_CONFIG_FILE", filepath.Join(t.TempDir(), "config.json"))
	t.Setenv("IKUAI_CLI_BASE_URL", "")
	t.Setenv("IKUAI_CLI_TOKEN", "")

	r := New(&bytes.Buffer{}, &bytes.Buffer{})
	if err := r.SyncSession(); err != nil {
		t.Fatalf("SyncSession() error = %v", err)
	}

	if r.CredSource != "none" {
		t.Errorf("CredSource = %q, want %q", r.CredSource, "none")
	}
	if r.APIClient != nil {
		t.Error("APIClient should be nil when no credentials")
	}
}

func TestSyncSessionEnvPartialNotActivated(t *testing.T) {
	t.Setenv("IKUAI_CLI_CONFIG_FILE", filepath.Join(t.TempDir(), "config.json"))

	// Only BASE_URL set, no TOKEN.
	t.Setenv("IKUAI_CLI_BASE_URL", "https://env-router")
	t.Setenv("IKUAI_CLI_TOKEN", "")

	r := New(&bytes.Buffer{}, &bytes.Buffer{})
	if err := r.SyncSession(); err != nil {
		t.Fatalf("SyncSession() error = %v", err)
	}

	if r.CredSource != "none" {
		t.Errorf("CredSource = %q, want %q (partial env should not activate)", r.CredSource, "none")
	}
}

func TestSyncSessionEnvWhitespaceOnlyTreatedAsEmpty(t *testing.T) {
	t.Setenv("IKUAI_CLI_CONFIG_FILE", filepath.Join(t.TempDir(), "config.json"))
	t.Setenv("IKUAI_CLI_BASE_URL", "  ")
	t.Setenv("IKUAI_CLI_TOKEN", "  ")

	r := New(&bytes.Buffer{}, &bytes.Buffer{})
	if err := r.SyncSession(); err != nil {
		t.Fatalf("SyncSession() error = %v", err)
	}

	if r.CredSource != "none" {
		t.Errorf("CredSource = %q, want %q (whitespace-only env treated as empty)", r.CredSource, "none")
	}
}

func TestSyncSessionSessionHasBaseURLButNoToken_EnvWins(t *testing.T) {
	sf := filepath.Join(t.TempDir(), "config.json")
	t.Setenv("IKUAI_CLI_CONFIG_FILE", sf)
	t.Setenv("IKUAI_CLI_BASE_URL", "https://env-router")
	t.Setenv("IKUAI_CLI_TOKEN", "env-token")

	// Session has BaseURL from set-url but no Token (never logged in).
	if err := session.SaveBaseURL("https://session-router"); err != nil {
		t.Fatalf("SaveBaseURL() error = %v", err)
	}

	r := New(&bytes.Buffer{}, &bytes.Buffer{})
	if err := r.SyncSession(); err != nil {
		t.Fatalf("SyncSession() error = %v", err)
	}

	// No session Token → env wins.
	if r.CredSource != "env" {
		t.Errorf("CredSource = %q, want %q", r.CredSource, "env")
	}
	if r.Session.BaseURL != "https://env-router" {
		t.Errorf("BaseURL = %q, want env value", r.Session.BaseURL)
	}
}
