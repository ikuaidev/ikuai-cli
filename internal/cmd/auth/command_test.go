package auth

import (
	"bytes"
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
		name string
		args []string
		want string
	}{
		{
			name: "http scheme",
			args: []string{"set-url", "http://192.168.1.1/"},
			want: "https://192.168.1.1",
		},
		{
			name: "missing scheme",
			args: []string{"set-url", "192.168.1.1"},
			want: "https://192.168.1.1",
		},
		{
			name: "flag value",
			args: []string{"set-url", "--url", "router.local/"},
			want: "https://router.local",
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

			got := strings.TrimSpace(out.String())
			wantOutput := `{"base_url":"` + tt.want + `","message":"Base URL saved"}`
			if got != wantOutput {
				t.Fatalf("output = %q, want %q", got, wantOutput)
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
