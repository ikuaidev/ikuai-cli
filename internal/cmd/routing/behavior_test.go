package routing

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

func TestStaticListBuildsExpectedQueryParams(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-route"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/routing/static-routes" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/routing/static-routes")
			}
			q := req.URL.Query()
			if q.Get("page") != "5" {
				t.Fatalf("page = %q, want %q", q.Get("page"), "5")
			}
			if q.Get("page_size") != "30" {
				t.Fatalf("page_size = %q, want %q", q.Get("page_size"), "30")
			}
			if q.Get("filter") != "dst==10.0.0.0/24" {
				t.Fatalf("filter = %q, want %q", q.Get("filter"), "dst==10.0.0.0/24")
			}
			if q.Get("order") != "desc" {
				t.Fatalf("order = %q, want %q", q.Get("order"), "desc")
			}
			if q.Get("order_by") != "dst" {
				t.Fatalf("order_by = %q, want %q", q.Get("order_by"), "dst")
			}
			if got := req.Header.Get("Authorization"); got != "Bearer token-route" {
				t.Fatalf("Authorization = %q, want %q", got, "Bearer token-route")
			}
			return jsonResponse(`{"code":0,"data":{"items":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"static", "list", "--page", "5", "--page-size", "30", "--filter", "dst==10.0.0.0/24", "--order", "desc", "--order-by", "dst"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := `{"items":[]}` + "\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestStreamFiveTupleCreateSendsExpectedJSONBody(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-route"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodPost)
			}
			if req.URL.String() != "https://router.local/api/v4.0/routing/five-tuple-rules" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/routing/five-tuple-rules")
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			// writeCmd merges createDefaults into the body, so we check key fields.
			bs := string(body)
			for _, want := range []string{`"interface":"wan2"`, `"tagname":"office-stream"`, `"protocol":"any"`, `"prio":31`} {
				if !bytes.Contains(body, []byte(want)) {
					t.Fatalf("body missing %s: %s", want, bs)
				}
			}

			return jsonResponse(`{"code":0,"message":"created","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"stream", "five-tuple", "create", "--name", "office-stream", "--interface", "wan2"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := "{\"message\":\"created\"}\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestStaticCreateMissingRequiredFlags(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "tok"}

	cmd := New(app)
	cmd.SetArgs([]string{"static", "create", "--name", "test"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing required flags")
	}
	if !strings.Contains(err.Error(), "missing required flag") {
		t.Fatalf("error = %q, want it to contain 'missing required flag'", err.Error())
	}
}

func TestStreamDomainToggleMissingEnabled(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "tok"}

	cmd := New(app)
	cmd.SetArgs([]string{"stream", "domain", "toggle", "1"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing --enabled")
	}
	if !strings.Contains(err.Error(), "missing required flag: --enabled") {
		t.Fatalf("error = %q, want 'missing required flag: --enabled'", err.Error())
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
