package session

import (
	"os"
	"path/filepath"
	"strings"
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

func TestNormalizeBaseURLCorrectsCommonInputs(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        string
		wantChanged bool
	}{
		{
			name:        "bare host",
			input:       "192.168.1.1",
			want:        "https://192.168.1.1",
			wantChanged: true,
		},
		{
			name:        "http scheme",
			input:       "http://192.168.1.1/",
			want:        "https://192.168.1.1",
			wantChanged: true,
		},
		{
			name:        "path query and fragment",
			input:       "https://192.168.1.1/login?from=cli#top",
			want:        "https://192.168.1.1",
			wantChanged: true,
		},
		{
			name:        "host with port",
			input:       "192.168.1.1:8443",
			want:        "https://192.168.1.1:8443",
			wantChanged: true,
		},
		{
			name:        "already normalized",
			input:       "https://router.local",
			want:        "https://router.local",
			wantChanged: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, changed, err := NormalizeBaseURL(tt.input)
			if err != nil {
				t.Fatalf("NormalizeBaseURL() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("NormalizeBaseURL() = %q, want %q", got, tt.want)
			}
			if changed != tt.wantChanged {
				t.Fatalf("changed = %v, want %v", changed, tt.wantChanged)
			}
		})
	}
}

func TestNormalizeBaseURLRejectsUnclearInputs(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name:    "unsupported scheme",
			input:   "ftp://192.168.1.1",
			wantErr: "only accepts http or https",
		},
		{
			name:    "missing host",
			input:   "https://",
			wantErr: "missing router host",
		},
		{
			name:    "userinfo",
			input:   "https://admin:secret@192.168.1.1",
			wantErr: "must not include username or password",
		},
		{
			name:    "spaces",
			input:   "not a url with spaces",
			wantErr: "invalid router host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := NormalizeBaseURL(tt.input)
			if err == nil {
				t.Fatal("NormalizeBaseURL() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %q, want to contain %q", err.Error(), tt.wantErr)
			}
		})
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
