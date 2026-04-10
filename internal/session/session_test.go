package session

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveBaseURLUsesOverridePathAndTrimsSlash(t *testing.T) {
	t.Setenv(configFileEnv, filepath.Join(t.TempDir(), "config.json"))

	if err := SaveBaseURL("https://192.168.1.1/"); err != nil {
		t.Fatalf("SaveBaseURL() error = %v", err)
	}

	s, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if s.BaseURL != "https://192.168.1.1" {
		t.Fatalf("BaseURL = %q, want %q", s.BaseURL, "https://192.168.1.1")
	}

	if _, err := os.Stat(os.Getenv(configFileEnv)); err != nil {
		t.Fatalf("session file not written to override path: %v", err)
	}
}

func TestClearWipesURLTokenPreservesSSH(t *testing.T) {
	t.Setenv(configFileEnv, filepath.Join(t.TempDir(), "config.json"))

	if err := SaveLogin("https://192.168.1.1", "token-xyz"); err != nil {
		t.Fatalf("SaveLogin() error = %v", err)
	}
	if err := SaveSSHCreds("root", "secret", 22); err != nil {
		t.Fatalf("SaveSSHCreds() error = %v", err)
	}

	if err := Clear(); err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	s, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if s.BaseURL != "" {
		t.Errorf("BaseURL = %q, want empty", s.BaseURL)
	}
	if s.Token != "" {
		t.Errorf("Token = %q, want empty", s.Token)
	}
	if s.SSHUser != "root" {
		t.Errorf("SSHUser = %q, want %q", s.SSHUser, "root")
	}
	if s.SSHPassword != "secret" {
		t.Errorf("SSHPassword = %q, want %q", s.SSHPassword, "secret")
	}
	if s.SSHPort != 22 {
		t.Errorf("SSHPort = %d, want %d", s.SSHPort, 22)
	}
}

func TestClearOnMissingSessionFileSucceeds(t *testing.T) {
	t.Setenv(configFileEnv, filepath.Join(t.TempDir(), "config.json"))

	if err := Clear(); err != nil {
		t.Fatalf("Clear() on missing file error = %v", err)
	}

	s, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if s.BaseURL != "" || s.Token != "" {
		t.Errorf("session not empty after Clear() on missing file: %+v", s)
	}
}

func TestClearIsIdempotent(t *testing.T) {
	t.Setenv(configFileEnv, filepath.Join(t.TempDir(), "config.json"))

	if err := SaveLogin("https://192.168.1.1", "token"); err != nil {
		t.Fatalf("SaveLogin() error = %v", err)
	}
	if err := Clear(); err != nil {
		t.Fatalf("first Clear() error = %v", err)
	}
	if err := Clear(); err != nil {
		t.Fatalf("second Clear() error = %v", err)
	}

	s, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if s.BaseURL != "" || s.Token != "" {
		t.Errorf("session not empty after double Clear: %+v", s)
	}
}
