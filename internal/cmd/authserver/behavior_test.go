package authserver

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ikuaidev/ikuai-cli/internal/api"
	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/ikuaidev/ikuai-cli/internal/output"
	"github.com/ikuaidev/ikuai-cli/internal/session"
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

func TestSetSendsExpectedJSONBody(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-authsrv"}
	seenGet := false
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() != "https://router.local/api/v4.0/auth/web/services" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/auth/web/services")
			}
			if req.Method == http.MethodGet {
				seenGet = true
				return jsonResponse(`{"code":0,"message":"Success","results":{"data":[{"id":1,"enabled":"no","max_time":0,"idle_time":60,"user_auth":1,"interface":"lan1"}]}}`), nil
			}
			if req.Method != http.MethodPut {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodPut)
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			want := `{"enabled":"yes","id":1,"idle_time":120,"interface":"lan1","max_time":0,"user_auth":1}`
			if string(body) != want {
				t.Fatalf("body = %q, want %q", string(body), want)
			}

			return jsonResponse(`{"code":0,"message":"saved","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"set", "--enabled", "yes", "--idle-time", "120"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !seenGet {
		t.Fatal("set should read current config before PUT")
	}

	got := out.String()
	want := "{\"message\":\"saved\"}\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
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
