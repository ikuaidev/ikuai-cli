package routing

import (
	"bytes"
	"encoding/json"
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
			if q.Get("limit") != "30" {
				t.Fatalf("limit = %q, want %q", q.Get("limit"), "30")
			}
			if q.Has("page_size") {
				t.Fatalf("page_size query param should not be sent: %s", req.URL.RawQuery)
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

func TestStreamDomainListBuildsExpectedQueryParams(t *testing.T) {
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
			if req.URL.Path != "/api/v4.0/routing/domain-rules" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/routing/domain-rules")
			}
			q := req.URL.Query()
			if q.Get("page") != "2" {
				t.Fatalf("page = %q, want %q", q.Get("page"), "2")
			}
			if q.Get("limit") != "25" {
				t.Fatalf("limit = %q, want %q", q.Get("limit"), "25")
			}
			if q.Has("page_size") {
				t.Fatalf("page_size query param should not be sent: %s", req.URL.RawQuery)
			}
			if q.Get("filter") != "enabled==yes" {
				t.Fatalf("filter = %q, want %q", q.Get("filter"), "enabled==yes")
			}
			if q.Get("order") != "asc" {
				t.Fatalf("order = %q, want %q", q.Get("order"), "asc")
			}
			if q.Get("order_by") != "prio" {
				t.Fatalf("order_by = %q, want %q", q.Get("order_by"), "prio")
			}
			if got := req.Header.Get("Authorization"); got != "Bearer token-route" {
				t.Fatalf("Authorization = %q, want %q", got, "Bearer token-route")
			}
			return jsonResponse(`{"code":0,"data":{"items":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"stream", "domain", "list", "--page", "2", "--page-size", "25", "--filter", "enabled==yes", "--order", "asc", "--order-by", "prio"})
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

func TestStaticGetUsesYAMLEndpoint(t *testing.T) {
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
			if req.URL.Path != "/api/v4.0/routing/static-routes/7" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/routing/static-routes/7")
			}
			return jsonResponse(`{"code":0,"message":"Success","results":{"total":1,"data":[{"id":7,"tagname":"route-test"}]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"static", "get", "7"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got := out.String(); !strings.Contains(got, `"tagname":"route-test"`) {
		t.Fatalf("output = %q, want route-test JSON", got)
	}
}

func TestDomainUpdateReadsCurrentAndSendsFullJSONBody(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	requests := 0
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-route"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requests++
			if req.URL.Path != "/api/v4.0/routing/domain-rules/7" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/routing/domain-rules/7")
			}
			switch requests {
			case 1:
				if req.Method != http.MethodGet {
					t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
				}
				return jsonResponse(`{"code":0,"message":"Success","data":{"id":7,"tagname":"old-domain","enabled":"yes","comment":"","interface":"wan1","prio":31,"domain":{"custom":["old.example"],"object":{}},"src_addr":{"custom":{},"object":{}},"time":{"custom":[{"type":"weekly","weekdays":"1234567","start_time":"00:00","end_time":"23:59","comment":""}],"object":{}}}}`), nil
			case 2:
				if req.Method != http.MethodPut {
					t.Fatalf("method = %q, want %q", req.Method, http.MethodPut)
				}
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("ReadAll() error = %v", err)
				}
				var got map[string]interface{}
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatalf("json.Unmarshal() error = %v; body=%s", err, string(body))
				}
				if _, ok := got["id"]; ok {
					t.Fatalf("body should not include id: %s", string(body))
				}
				for key, want := range map[string]interface{}{"tagname": "new-domain", "comment": "changed", "interface": "wan1"} {
					if got[key] != want {
						t.Fatalf("%s = %#v, want %#v; body=%s", key, got[key], want, string(body))
					}
				}
				domain, ok := got["domain"].(map[string]interface{})
				if !ok {
					t.Fatalf("domain missing from body: %s", string(body))
				}
				custom, ok := domain["custom"].([]interface{})
				if !ok || len(custom) != 1 || custom[0] != "old.example" {
					t.Fatalf("domain.custom = %#v, want [old.example]", domain["custom"])
				}
				object, ok := domain["object"].([]interface{})
				if !ok || len(object) != 0 {
					t.Fatalf("domain.object = %#v, want empty array", domain["object"])
				}
				return jsonResponse(`{"code":0,"message":"success"}`), nil
			default:
				t.Fatalf("unexpected request %d", requests)
			}
			return nil, nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"stream", "domain", "update", "7", "--name", "new-domain", "--comment", "changed"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if requests != 2 {
		t.Fatalf("requests = %d, want 2", requests)
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
