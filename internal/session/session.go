// Package session manages the ~/.ikuai-cli/config.json file.
package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const configFileEnv = "IKUAI_CLI_CONFIG_FILE"

// Session holds CLI configuration and credentials.
type Session struct {
	BaseURL     string `json:"base_url"`
	Token       string `json:"token"`
	SSHUser     string `json:"ssh_user,omitempty"`
	SSHPassword string `json:"ssh_password,omitempty"`
	SSHPort     int    `json:"ssh_port,omitempty"`
}

func configFile() string {
	if path := strings.TrimSpace(os.Getenv(configFileEnv)); path != "" {
		return path
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ikuai-cli", "config.json")
}

// Load reads the session file. Returns an empty Session on missing/corrupt file.
func Load() (*Session, error) {
	data, err := os.ReadFile(configFile())
	if err != nil {
		if os.IsNotExist(err) {
			return &Session{}, nil
		}
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return &Session{}, nil
	}
	return &s, nil
}

func save(s *Session) error {
	path := configFile()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// SaveBaseURL stores the base URL (trailing slash stripped).
func SaveBaseURL(url string) error {
	s, _ := Load()
	if s == nil {
		s = &Session{}
	}
	s.BaseURL = strings.TrimRight(url, "/")
	return save(s)
}

// SaveToken stores the JWT token.
func SaveToken(token string) error {
	s, _ := Load()
	if s == nil {
		s = &Session{}
	}
	s.Token = token
	return save(s)
}

// SaveLogin stores the base URL and token in one call.
func SaveLogin(url, token string) error {
	s, _ := Load()
	if s == nil {
		s = &Session{}
	}
	s.BaseURL = strings.TrimRight(url, "/")
	s.Token = token
	return save(s)
}

// Clear removes the base URL and token from the session.
// SSH credentials are intentionally preserved — use a dedicated command to
// clear those. After Clear(), the session file equates to "never logged in".
func Clear() error {
	s, _ := Load()
	if s == nil {
		s = &Session{}
	}
	s.BaseURL = ""
	s.Token = ""
	return save(s)
}

// SaveSSHCreds stores SSH credentials.
func SaveSSHCreds(user, password string, port int) error {
	s, _ := Load()
	if s == nil {
		s = &Session{}
	}
	s.SSHUser = user
	s.SSHPassword = password
	s.SSHPort = port
	return save(s)
}
