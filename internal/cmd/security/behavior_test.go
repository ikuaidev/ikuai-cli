package security

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

func TestACLListBuildsExpectedQueryParams(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-sec"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/security/acl-rules" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/security/acl-rules")
			}
			q := req.URL.Query()
			if q.Get("page") != "7" {
				t.Fatalf("page = %q, want %q", q.Get("page"), "7")
			}
			if q.Get("page_size") != "12" {
				t.Fatalf("page_size = %q, want %q", q.Get("page_size"), "12")
			}
			if q.Get("filter") != "enabled==yes" {
				t.Fatalf("filter = %q, want %q", q.Get("filter"), "enabled==yes")
			}
			if q.Get("order") != "desc" {
				t.Fatalf("order = %q, want %q", q.Get("order"), "desc")
			}
			if q.Get("order_by") != "id" {
				t.Fatalf("order_by = %q, want %q", q.Get("order_by"), "id")
			}
			if got := req.Header.Get("Authorization"); got != "Bearer token-sec" {
				t.Fatalf("Authorization = %q, want %q", got, "Bearer token-sec")
			}
			return jsonResponse(`{"code":0,"data":{"items":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"acl", "list", "--page", "7", "--page-size", "12", "--filter", "enabled==yes", "--order", "desc", "--order-by", "id"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := `{"items":[]}` + "\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestMACSetModeSendsExpectedJSONBody(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-sec"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodPut)
			}
			if req.URL.String() != "https://router.local/api/v4.0/security/mac-mode" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/security/mac-mode")
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			if string(body) != `{"mode":"blacklist"}` {
				t.Fatalf("body = %q, want %q", string(body), `{"mode":"blacklist"}`)
			}

			return jsonResponse(`{"code":0,"message":"saved","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"mac", "set-mode", "--data", `{"mode":"blacklist"}`})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := "{\"message\":\"saved\"}\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestACLCreateMissingRequiredFlags(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "tok"}

	cmd := New(app)
	cmd.SetArgs([]string{"acl", "create"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing required flags")
	}
	if !strings.Contains(err.Error(), "missing required flags") {
		t.Fatalf("error = %q, want it to contain 'missing required flags'", err.Error())
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
