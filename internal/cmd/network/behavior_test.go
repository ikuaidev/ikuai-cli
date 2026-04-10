package network

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

func TestNatListBuildsExpectedQueryParams(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-abc"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/network/nat/rules" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/network/nat/rules")
			}
			q := req.URL.Query()
			if q.Get("page") != "2" {
				t.Fatalf("page = %q, want %q", q.Get("page"), "2")
			}
			if q.Get("page_size") != "50" {
				t.Fatalf("page_size = %q, want %q", q.Get("page_size"), "50")
			}
			if q.Get("filter") != "enabled==true" {
				t.Fatalf("filter = %q, want %q", q.Get("filter"), "enabled==true")
			}
			if q.Get("order") != "asc" {
				t.Fatalf("order = %q, want %q", q.Get("order"), "asc")
			}
			if q.Get("order_by") != "id" {
				t.Fatalf("order_by = %q, want %q", q.Get("order_by"), "id")
			}
			if got := req.Header.Get("Authorization"); got != "Bearer token-abc" {
				t.Fatalf("Authorization = %q, want %q", got, "Bearer token-abc")
			}
			return jsonResponse(`{"code":0,"data":{"items":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"nat", "list", "--page", "2", "--page-size", "50", "--filter", "enabled==true", "--order", "asc", "--order-by", "id"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := `{"items":[]}` + "\n"
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
