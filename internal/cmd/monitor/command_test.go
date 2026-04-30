package monitor

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

func TestSystemCommandRequestsMonitoringSystemEndpoint(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-123"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.String() != "https://router.local/api/v4.0/monitoring/system" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/monitoring/system")
			}
			if got := req.Header.Get("Authorization"); got != "Bearer token-123" {
				t.Fatalf("Authorization = %q, want %q", got, "Bearer token-123")
			}
			return jsonResponse(`{"code":0,"message":"ok","data":{"uptime":123}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"system"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := `{"uptime":123}` + "\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestAppProtocolsHistoryPassesYamlQueryParameters(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-123"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/monitoring/app-protocols/history-load" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/monitoring/app-protocols/history-load")
			}
			query := req.URL.Query()
			for name, want := range map[string]string{
				"appids":    "2580003,2580004",
				"starttime": "1773215100",
				"stoptime":  "1773218700",
			} {
				if got := query.Get(name); got != want {
					t.Fatalf("query %s = %q, want %q", name, got, want)
				}
			}
			return jsonResponse(`{"code":0,"message":"ok","data":[{"appid":2580003}]}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"app-protocols-history", "--appids", "2580003,2580004", "--starttime", "1773215100", "--stoptime", "1773218700"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestAppTrafficSummaryMapsPageSizeToYamlLimit(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-123"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/monitoring/app-traffic-summary" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/monitoring/app-traffic-summary")
			}
			query := req.URL.Query()
			if got := query.Get("page"); got != "2" {
				t.Fatalf("page = %q, want %q", got, "2")
			}
			if got := query.Get("limit"); got != "100" {
				t.Fatalf("limit = %q, want %q", got, "100")
			}
			if got := query.Get("page_size"); got != "" {
				t.Fatalf("page_size = %q, want empty", got)
			}
			return jsonResponse(`{"code":0,"message":"ok","results":{"proto3_day":[],"proto3_day_total":0}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"app-traffic-summary", "--page", "2", "--page-size", "100"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
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
