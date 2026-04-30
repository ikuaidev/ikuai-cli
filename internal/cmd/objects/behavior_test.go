package objects

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

func TestIPListBuildsExpectedQueryParams(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-obj"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/ip-objects" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/ip-objects")
			}
			q := req.URL.Query()
			if q.Get("page") != "4" {
				t.Fatalf("page = %q, want %q", q.Get("page"), "4")
			}
			if q.Get("limit") != "15" {
				t.Fatalf("limit = %q, want %q", q.Get("limit"), "15")
			}
			for _, unsupported := range []string{"page_size", "filter", "order", "order_by"} {
				if q.Has(unsupported) {
					t.Fatalf("query unexpectedly includes %q: %v", unsupported, q)
				}
			}
			if got := req.Header.Get("Authorization"); got != "Bearer token-obj" {
				t.Fatalf("Authorization = %q, want %q", got, "Bearer token-obj")
			}
			return jsonResponse(`{"code":0,"data":{"items":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"ip", "list", "--page", "4", "--page-size", "15"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := `{"items":[]}` + "\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestIPCreateSendsExpectedJSONBody(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-obj"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodPost)
			}
			if req.URL.String() != "https://router.local/api/v4.0/ip-objects" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/ip-objects")
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			bs := string(body)
			for _, want := range []string{`"group_name":"office-gateway"`, `"group_value":`} {
				if !bytes.Contains(body, []byte(want)) {
					t.Fatalf("body missing %s: %s", want, bs)
				}
			}

			return jsonResponse(`{"code":0,"message":"created","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"ip", "create", "--name", "office-gateway", "--value", "10.0.0.1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := "{\"message\":\"created\"}\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestTimeCreateSendsSemanticFlagsAsGroupValue(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-obj"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodPost)
			}
			if req.URL.String() != "https://router.local/api/v4.0/time-objects" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/time-objects")
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			bs := string(body)
			for _, want := range []string{
				`"group_name":"office"`,
				`"type":"weekly"`,
				`"weekdays":"12345"`,
				`"start_time":"09:00"`,
				`"end_time":"18:00"`,
				`"comment":"work"`,
			} {
				if !bytes.Contains(body, []byte(want)) {
					t.Fatalf("body missing %s: %s", want, bs)
				}
			}

			return jsonResponse(`{"code":0,"message":"created","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"time", "create", "--name", "office", "--type", "weekly", "--weekdays", "12345", "--start-time", "09:00", "--end-time", "18:00", "--comment", "work"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestObjectCommandsDoNotExposeUnsupportedToggle(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	cmd := New(app)
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"ip", "--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if strings.Contains(out.String(), "toggle") {
		t.Fatalf("help unexpectedly exposes toggle: %s", out.String())
	}
}

func TestObjectListDoesNotExposeUnsupportedQueryFlags(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	cmd := New(app)
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"ip", "list", "--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	help := out.String()
	for _, unsupported := range []string{"--filter", "--order ", "--order-by"} {
		if strings.Contains(help, unsupported) {
			t.Fatalf("help unexpectedly exposes %s: %s", unsupported, help)
		}
	}
	for _, expected := range []string{"--page", "--page-size"} {
		if !strings.Contains(help, expected) {
			t.Fatalf("help missing %s: %s", expected, help)
		}
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
