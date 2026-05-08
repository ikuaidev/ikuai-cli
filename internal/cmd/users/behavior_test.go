package users

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

func TestAccountsListBuildsExpectedQueryParams(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-users"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/auth/users" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/auth/users")
			}
			q := req.URL.Query()
			if q.Get("page") != "2" {
				t.Fatalf("page = %q, want %q", q.Get("page"), "2")
			}
			if q.Get("limit") != "25" {
				t.Fatalf("limit = %q, want %q", q.Get("limit"), "25")
			}
			if q.Get("page_size") != "" {
				t.Fatalf("page_size = %q, want empty", q.Get("page_size"))
			}
			if q.Get("filter") != "role==admin" {
				t.Fatalf("filter = %q, want %q", q.Get("filter"), "role==admin")
			}
			if q.Get("order") != "desc" {
				t.Fatalf("order = %q, want %q", q.Get("order"), "desc")
			}
			if q.Get("order_by") != "username" {
				t.Fatalf("order_by = %q, want %q", q.Get("order_by"), "username")
			}
			if got := req.Header.Get("Authorization"); got != "Bearer token-users" {
				t.Fatalf("Authorization = %q, want %q", got, "Bearer token-users")
			}
			return jsonResponse(`{"code":0,"data":{"items":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"accounts", "list", "--page", "2", "--page-size", "25", "--filter", "role==admin", "--order", "desc", "--order-by", "username"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := `{"items":[]}` + "\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestOnlineListBuildsExpectedQueryParams(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-users"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/auth/online-users" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/auth/online-users")
			}
			q := req.URL.Query()
			for name, want := range map[string]string{
				"page":     "2",
				"limit":    "25",
				"KEYWORDS": "alice",
				"FINDS":    "username,name",
				"ORDER":    "desc",
				"ORDER_BY": "auth_time",
			} {
				if got := q.Get(name); got != want {
					t.Fatalf("%s = %q, want %q", name, got, want)
				}
			}
			if q.Get("page_size") != "" || q.Get("filter") != "" || q.Get("order_by") != "" {
				t.Fatalf("unexpected lowercase pagination/search params: %s", req.URL.RawQuery)
			}
			return jsonResponse(`{"code":0,"data":{"items":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"online", "--page", "2", "--page-size", "25", "--keywords", "alice", "--finds", "username,name", "--order", "desc", "--order-by", "auth_time"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestOnlineGetFetchesOnlineUser(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-users"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.String() != "https://router.local/api/v4.0/auth/online-users/42" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/auth/online-users/42")
			}
			return jsonResponse(`{"code":0,"data":{"id":42,"username":"alice"}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"online", "get", "42"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := "{\"id\":42,\"username\":\"alice\"}\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestKickDeletesOnlineUser(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-users"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodDelete)
			}
			if req.URL.String() != "https://router.local/api/v4.0/auth/online-users/42" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/auth/online-users/42")
			}
			return jsonResponse(`{"code":0,"message":"success","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"kick", "42", "--yes"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestAccountsCreateSendsExpectedJSONBody(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-users"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodPost)
			}
			if req.URL.String() != "https://router.local/api/v4.0/auth/users" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/auth/users")
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			// writeCmd merges createDefaults into the body, so we check key fields.
			bs := string(body)
			for _, want := range []string{`"username":"alice"`, `"passwd":"secret"`, `"ppptype":"any"`, `"share":2`, `"upload":128`, `"download":256`} {
				if !bytes.Contains(body, []byte(want)) {
					t.Fatalf("body missing %s: %s", want, bs)
				}
			}

			return jsonResponse(`{"code":0,"message":"created","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"accounts", "create", "--username", "alice", "--password", "secret", "--share", "2", "--upload", "128", "--download", "256"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := "{\"message\":\"created\"}\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestPackagesCreateSendsNumericFields(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-users"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodPost)
			}
			if req.URL.String() != "https://router.local/api/v4.0/auth/packages" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/auth/packages")
			}
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			bs := string(body)
			for _, want := range []string{`"packname":"pkg1"`, `"packtime":"24h"`, `"price":100`, `"up_speed":500`, `"down_speed":1000`} {
				if !bytes.Contains(body, []byte(want)) {
					t.Fatalf("body missing %s: %s", want, bs)
				}
			}
			return jsonResponse(`{"code":0,"message":"created","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"packages", "create", "--name", "pkg1", "--time", "24h", "--price", "100", "--up-speed", "500", "--down-speed", "1000"})
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
