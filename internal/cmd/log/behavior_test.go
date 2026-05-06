package logcmd

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

func TestSystemLogBuildsExpectedQueryParams(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-log"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/log/system" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/log/system")
			}
			q := req.URL.Query()
			for name, want := range map[string]string{
				"page":     "3",
				"limit":    "9",
				"filter":   "level==error",
				"key":      "level,module",
				"pattern":  "error",
				"order":    "desc",
				"order_by": "timestamp",
			} {
				if got := q.Get(name); got != want {
					t.Fatalf("%s = %q, want %q; raw=%s", name, got, want, req.URL.RawQuery)
				}
			}
			if q.Get("page_size") != "" {
				t.Fatalf("page_size = %q, want empty; raw=%s", q.Get("page_size"), req.URL.RawQuery)
			}
			return jsonResponse(`{"code":0,"data":{"items":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"system", "list", "--page", "3", "--page-size", "9", "--filter", "level==error", "--key", "level,module", "--pattern", "error", "--order", "desc", "--order-by", "timestamp"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := `{"items":[]}` + "\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestSystemLogClearRequestsDelete(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-log"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodDelete)
			}
			if req.URL.String() != "https://router.local/api/v4.0/log/system" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/log/system")
			}
			return jsonResponse(`{"code":0,"message":"cleared","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"system", "delete", "--yes"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := "{\"message\":\"cleared\"}\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestDdnsLogListUsesResultColumn(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.Table
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-log"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/log/ddns" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/log/ddns")
			}
			return jsonResponse(`{"code":0,"message":"Success","results":{"data":[{"id":1,"timestamp":1778060866,"domain":"www.ali123.com","interface":"auto","ip_addr":"--","result":"失败","event":"错误: 认证失败"}],"total":1}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"ddns", "list", "--page", "1", "--page-size", "5"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "RESULT") || !strings.Contains(got, "失败") {
		t.Fatalf("ddns table should show result column and value: %q", got)
	}
	if strings.Contains(got, "STATUS") {
		t.Fatalf("ddns table should not show stale status column: %q", got)
	}
}

func TestNoticeLogListUsesEventColumn(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.Table
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-log"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/log/notice" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/log/notice")
			}
			return jsonResponse(`{"code":0,"message":"Success","results":{"data":[{"id":1,"timestamp":1778064504,"type":"实时通知","ip_addr":"192.168.99.100","event":"已收到通知"}],"total":1}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"notice", "list", "--page", "1", "--page-size", "5"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "EVENT") || !strings.Contains(got, "已收到通知") {
		t.Fatalf("notice table should show event column and value: %q", got)
	}
	if strings.Contains(got, "TITLE") || strings.Contains(got, "CONTENT") {
		t.Fatalf("notice table should not show stale title/content columns: %q", got)
	}
}

func TestSystemLogDeleteRequiresYes(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-log"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			t.Fatalf("unexpected HTTP request: %s %s", req.Method, req.URL.String())
			return nil, nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"system", "delete"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() error = nil, want confirmation error")
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
