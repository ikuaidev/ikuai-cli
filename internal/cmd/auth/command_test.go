package auth

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/ikuaidev/ikuai-cli/internal/output"
	"github.com/ikuaidev/ikuai-cli/internal/session"
)

func TestSetURLTrimsTrailingSlash(t *testing.T) {
	t.Setenv("IKUAI_CLI_CONFIG_FILE", t.TempDir()+"/config.json")

	var out bytes.Buffer
	cmd := New(cliapp.New(&out, &out))
	cmd.SetArgs([]string{"set-url", "https://192.168.1.1/"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	s, err := session.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if s.BaseURL != "https://192.168.1.1" {
		t.Fatalf("BaseURL = %q, want %q", s.BaseURL, "https://192.168.1.1")
	}
}

func TestSetURLNormalizesToHTTPS(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		want       string
		normalized bool
		wantInput  string
	}{
		{
			name:       "http scheme",
			args:       []string{"set-url", "http://192.168.1.1/"},
			want:       "https://192.168.1.1",
			normalized: true,
			wantInput:  "http://192.168.1.1/",
		},
		{
			name:       "missing scheme",
			args:       []string{"set-url", "192.168.1.1"},
			want:       "https://192.168.1.1",
			normalized: true,
			wantInput:  "192.168.1.1",
		},
		{
			name:       "flag value",
			args:       []string{"set-url", "--url", "router.local/"},
			want:       "https://router.local",
			normalized: true,
			wantInput:  "router.local/",
		},
		{
			name:       "path query and fragment",
			args:       []string{"set-url", "https://192.168.1.1/login?from=cli#top"},
			want:       "https://192.168.1.1",
			normalized: true,
		},
		{
			name:       "already normalized",
			args:       []string{"set-url", "https://router.local"},
			want:       "https://router.local",
			normalized: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("IKUAI_CLI_CONFIG_FILE", t.TempDir()+"/config.json")

			var out bytes.Buffer
			app := cliapp.New(&out, &out)
			app.Format = output.JSON

			cmd := New(app)
			cmd.SetArgs(tt.args)
			if err := cmd.Execute(); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			s, err := session.Load()
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if s.BaseURL != tt.want {
				t.Fatalf("BaseURL = %q, want %q", s.BaseURL, tt.want)
			}

			var got struct {
				Message    string `json:"message"`
				BaseURL    string `json:"base_url"`
				Normalized bool   `json:"normalized"`
				InputURL   string `json:"input_url,omitempty"`
			}
			if err := json.Unmarshal(bytes.TrimSpace(out.Bytes()), &got); err != nil {
				t.Fatalf("output JSON error = %v; output = %q", err, out.String())
			}
			if got.Message != "Base URL saved" {
				t.Fatalf("message = %q, want %q", got.Message, "Base URL saved")
			}
			if got.BaseURL != tt.want {
				t.Fatalf("output base_url = %q, want %q", got.BaseURL, tt.want)
			}
			if got.Normalized != tt.normalized {
				t.Fatalf("normalized = %v, want %v", got.Normalized, tt.normalized)
			}
			if got.InputURL != tt.wantInput {
				t.Fatalf("input_url = %q, want %q", got.InputURL, tt.wantInput)
			}
			if tt.normalized && strings.ContainsAny(got.InputURL, "?#@") {
				t.Fatalf("input_url = %q, want sensitive URL parts omitted", got.InputURL)
			}
			if !tt.normalized && got.InputURL != "" {
				t.Fatalf("input_url = %q, want empty for unchanged input", got.InputURL)
			}
		})
	}
}

func TestSetURLRejectsInvalidRouterURL(t *testing.T) {
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
			t.Setenv("IKUAI_CLI_CONFIG_FILE", t.TempDir()+"/config.json")

			var out bytes.Buffer
			cmd := New(cliapp.New(&out, &out))
			cmd.SetArgs([]string{"set-url", tt.input})

			err := cmd.Execute()
			if err == nil {
				t.Fatal("Execute() error = nil, want error")
			}
			for _, want := range []string{
				"Invalid router URL",
				tt.wantErr,
				"ikuai-cli auth set-url 192.168.1.1",
				"ikuai-cli auth set-url https://192.168.1.1",
			} {
				if !strings.Contains(err.Error(), want) {
					t.Fatalf("error = %q, want to contain %q", err.Error(), want)
				}
			}

			s, loadErr := session.Load()
			if loadErr != nil {
				t.Fatalf("Load() error = %v", loadErr)
			}
			if s.BaseURL != "" {
				t.Fatalf("BaseURL = %q, want empty after invalid URL", s.BaseURL)
			}
		})
	}
}

func TestSetURLRequiresNonBlankURL(t *testing.T) {
	t.Setenv("IKUAI_CLI_CONFIG_FILE", t.TempDir()+"/config.json")

	var out bytes.Buffer
	cmd := New(cliapp.New(&out, &out))
	cmd.SetArgs([]string{"set-url", "   "})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when URL is blank, got nil")
	}
}

func TestSetTokenSavesToken(t *testing.T) {
	t.Setenv("IKUAI_CLI_CONFIG_FILE", t.TempDir()+"/config.json")

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON

	cmd := New(app)
	cmd.SetArgs([]string{"set-token", "my-test-token-123"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	s, err := session.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if s.Token != "my-test-token-123" {
		t.Fatalf("Token = %q, want %q", s.Token, "my-test-token-123")
	}

	got := strings.TrimSpace(out.String())
	want := `{"message":"Token saved"}`
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestSetTokenRequiresArg(t *testing.T) {
	t.Setenv("IKUAI_CLI_CONFIG_FILE", t.TempDir()+"/config.json")

	var out bytes.Buffer
	cmd := New(cliapp.New(&out, &out))
	cmd.SetArgs([]string{"set-token"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when no token arg provided, got nil")
	}
}

func TestClearClearsBaseURLAndToken(t *testing.T) {
	t.Setenv("IKUAI_CLI_CONFIG_FILE", t.TempDir()+"/config.json")

	if err := session.SaveLogin("https://192.168.1.1", "token-123"); err != nil {
		t.Fatalf("SaveLogin() error = %v", err)
	}

	var out bytes.Buffer
	cmd := New(cliapp.New(&out, &out))
	cmd.SetArgs([]string{"clear"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	s, err := session.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if s.BaseURL != "" {
		t.Errorf("BaseURL = %q, want empty", s.BaseURL)
	}
	if s.Token != "" {
		t.Errorf("Token = %q, want empty", s.Token)
	}

	got := strings.TrimSpace(out.String())
	if !strings.Contains(got, "Cleared") {
		t.Fatalf("output = %q, want to contain 'Cleared'", got)
	}
}
