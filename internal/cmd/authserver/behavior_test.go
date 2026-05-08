package authserver

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ikuaidev/ikuai-cli/internal/api"
	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/ikuaidev/ikuai-cli/internal/output"
	"github.com/ikuaidev/ikuai-cli/internal/session"
	"github.com/spf13/cobra"
)

func TestGetRequestsExpectedEndpoint(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-authsrv"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.String() != "https://router.local/api/v4.0/auth/web/services" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/auth/web/services")
			}
			return jsonResponse(`{"code":0,"data":{"enable":"yes"}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"get"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := `{"enable":"yes"}` + "\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestGetUsesDefaultColumns(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.Table
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-authsrv"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			return jsonResponse(`{"code":0,"message":"Success","results":{"data":[{"id":1,"enabled":"no","interface":"lan1","idle_time":60,"max_time":0,"user_auth":1,"coupon_auth":0,"phone_auth":0,"static_pwd":1,"nopasswd":0,"https_redirect":0,"group_key":"secret"}]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"get"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "ENABLED") || !strings.Contains(got, "lan1") {
		t.Fatalf("auth-server get table should show default columns: %q", got)
	}
	if strings.Contains(got, "GROUP_KEY") || strings.Contains(got, "secret") {
		t.Fatalf("auth-server get table should not show non-default sensitive columns: %q", got)
	}
}

func TestAuthServerDoesNotExposeSetWithoutYAMLPut(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	cmd := New(app)
	if findSubcommand(cmd, "set") != nil {
		t.Fatal("auth-server set must not be exposed while YAML lacks PUT /api/v4.0/auth/web/services")
	}

	root, err := findRepoRoot()
	if err != nil {
		t.Skipf("repo root not found: %v", err)
	}
	yamlPath := filepath.Join(root, "web-api-yaml", "yaml", "auth", "auth-web-services.yaml")
	yamlBody, err := os.ReadFile(yamlPath) // #nosec G304 -- test reads a fixed repository fixture path.
	if err != nil {
		t.Skipf("YAML contract not available at %s: %v", yamlPath, err)
	}
	if authServerYAMLDeclaresPut(string(yamlBody)) {
		t.Fatal("YAML now declares PUT /api/v4.0/auth/web/services; update CLI and this guard test intentionally")
	}
}

func authServerYAMLDeclaresPut(body string) bool {
	inPath := false
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "  /api/v4.0/") {
			inPath = strings.TrimSuffix(strings.TrimSpace(line), ":") == "/api/v4.0/auth/web/services"
			continue
		}
		if inPath && strings.TrimSpace(line) == "put:" {
			return true
		}
	}
	return false
}

func findSubcommand(cmd *cobra.Command, name string) *cobra.Command {
	for _, child := range cmd.Commands() {
		if child.Name() == name {
			return child
		}
	}
	return nil
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}
