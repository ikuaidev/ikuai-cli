package objects

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
			if q.Get("page_size") != "15" {
				t.Fatalf("page_size = %q, want %q", q.Get("page_size"), "15")
			}
			if q.Get("filter") != "group==office" {
				t.Fatalf("filter = %q, want %q", q.Get("filter"), "group==office")
			}
			if q.Get("order") != "asc" {
				t.Fatalf("order = %q, want %q", q.Get("order"), "asc")
			}
			if q.Get("order_by") != "name" {
				t.Fatalf("order_by = %q, want %q", q.Get("order_by"), "name")
			}
			if got := req.Header.Get("Authorization"); got != "Bearer token-obj" {
				t.Fatalf("Authorization = %q, want %q", got, "Bearer token-obj")
			}
			return jsonResponse(`{"code":0,"data":{"items":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"ip", "list", "--page", "4", "--page-size", "15", "--filter", "group==office", "--order", "asc", "--order-by", "name"})
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
			if string(body) != `{"address":"10.0.0.1","name":"office-gateway"}` &&
				string(body) != `{"name":"office-gateway","address":"10.0.0.1"}` {
				t.Fatalf("body = %q", string(body))
			}

			return jsonResponse(`{"code":0,"message":"created","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"ip", "create", "--data", `{"name":"office-gateway","address":"10.0.0.1"}`})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := "{\"message\":\"created\"}\n"
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
