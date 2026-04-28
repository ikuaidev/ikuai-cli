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
			if string(body) != `{"acl_mac":1}` {
				t.Fatalf("body = %q, want %q", string(body), `{"acl_mac":1}`)
			}

			return jsonResponse(`{"code":0,"message":"saved","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"mac", "set-mode", "--acl-mac", "1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := "{\"message\":\"saved\"}\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestSecurityUpdateUsesFullBodyFromGet(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-sec"}
	step := 0
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			step++
			switch step {
			case 1:
				if req.Method != http.MethodGet {
					t.Fatalf("step 1 method = %q, want GET", req.Method)
				}
				if req.URL.Path != "/api/v4.0/security/acl-rules/42" {
					t.Fatalf("step 1 path = %q", req.URL.Path)
				}
				return jsonResponse(`{"code":0,"data":{"id":42,"enabled":"yes","tagname":"old","action":"drop","protocol":"tcp","dir":"forward","ctdir":0,"iinterface":"any","ointerface":"any","prio":50,"src_addr":{"custom":[],"object":[]},"dst_addr":{"custom":[],"object":[]},"src_port":{"custom":[],"object":[]},"dst_port":{"custom":[],"object":[]},"src_addr_inv":0,"dst_addr_inv":0,"src6_mode":0,"dst6_mode":0}}`), nil
			case 2:
				if req.Method != http.MethodPut {
					t.Fatalf("step 2 method = %q, want PUT", req.Method)
				}
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("ReadAll() error = %v", err)
				}
				got := string(body)
				if strings.Contains(got, `"id"`) {
					t.Fatalf("body should not include id: %s", got)
				}
				for _, want := range []string{`"tagname":"new-name"`, `"action":"drop"`, `"protocol":"tcp"`, `"dir":"forward"`} {
					if !strings.Contains(got, want) {
						t.Fatalf("body = %s, missing %s", got, want)
					}
				}
				return jsonResponse(`{"code":0,"message":"saved","data":null}`), nil
			default:
				t.Fatalf("unexpected request %d: %s %s", step, req.Method, req.URL.String())
			}
			return nil, nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"acl", "update", "42", "--name", "new-name"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if step != 2 {
		t.Fatalf("requests = %d, want 2", step)
	}
}

func TestTerminalsDoesNotExposeToggle(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "tok"}

	cmd := New(app)
	cmd.SetArgs([]string{"terminals", "toggle", "1", "--enabled", "no"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected unknown command error")
	}
	if !strings.Contains(err.Error(), "unknown") {
		t.Fatalf("error = %q, want unknown command/flag", err.Error())
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
