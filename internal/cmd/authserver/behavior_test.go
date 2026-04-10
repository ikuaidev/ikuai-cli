package authserver

import (
	"bytes"
	"io"
	"net/http"
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

func TestSetSendsExpectedJSONBody(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-authsrv"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodPut)
			}
			if req.URL.String() != "https://router.local/api/v4.0/auth/web/services" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/auth/web/services")
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			if string(body) != `{"enable":"yes"}` {
				t.Fatalf("body = %q, want %q", string(body), `{"enable":"yes"}`)
			}

			return jsonResponse(`{"code":0,"message":"saved","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"set", "--data", `{"enable":"yes"}`})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
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
