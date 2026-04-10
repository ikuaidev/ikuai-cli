package wireless

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

func TestBlacklistListBuildsExpectedQueryParams(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-wifi"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/wireless/access-control/rules" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/wireless/access-control/rules")
			}
			q := req.URL.Query()
			if q.Get("page") != "3" {
				t.Fatalf("page = %q, want %q", q.Get("page"), "3")
			}
			if q.Get("page_size") != "18" {
				t.Fatalf("page_size = %q, want %q", q.Get("page_size"), "18")
			}
			if q.Get("filter") != "enabled==yes" {
				t.Fatalf("filter = %q, want %q", q.Get("filter"), "enabled==yes")
			}
			if q.Get("order") != "desc" {
				t.Fatalf("order = %q, want %q", q.Get("order"), "desc")
			}
			if q.Get("order_by") != "mac" {
				t.Fatalf("order_by = %q, want %q", q.Get("order_by"), "mac")
			}
			if got := req.Header.Get("Authorization"); got != "Bearer token-wifi" {
				t.Fatalf("Authorization = %q, want %q", got, "Bearer token-wifi")
			}
			return jsonResponse(`{"code":0,"data":{"items":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"blacklist", "list", "--page", "3", "--page-size", "18", "--filter", "enabled==yes", "--order", "desc", "--order-by", "mac"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := `{"items":[]}` + "\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestACSetSendsExpectedJSONBody(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-wifi"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodPut)
			}
			if req.URL.String() != "https://router.local/api/v4.0/network/ac/services" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/network/ac/services")
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
	cmd.SetArgs([]string{"ac", "set", "--data", `{"enable":"yes"}`})
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
