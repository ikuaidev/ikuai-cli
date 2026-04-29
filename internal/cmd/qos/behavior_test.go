package qos

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
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-qos"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/network/qos/ip" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/network/qos/ip")
			}
			q := req.URL.Query()
			if q.Get("page") != "6" {
				t.Fatalf("page = %q, want %q", q.Get("page"), "6")
			}
			if q.Get("page_size") != "40" {
				t.Fatalf("page_size = %q, want %q", q.Get("page_size"), "40")
			}
			if q.Get("filter") != "enabled==yes" {
				t.Fatalf("filter = %q, want %q", q.Get("filter"), "enabled==yes")
			}
			if q.Get("order") != "asc" {
				t.Fatalf("order = %q, want %q", q.Get("order"), "asc")
			}
			if q.Get("order_by") != "priority" {
				t.Fatalf("order_by = %q, want %q", q.Get("order_by"), "priority")
			}
			if got := req.Header.Get("Authorization"); got != "Bearer token-qos" {
				t.Fatalf("Authorization = %q, want %q", got, "Bearer token-qos")
			}
			return jsonResponse(`{"code":0,"data":{"items":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"ip", "list", "--page", "6", "--page-size", "40", "--filter", "enabled==yes", "--order", "asc", "--order-by", "priority"})
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
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-qos"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodPost)
			}
			if req.URL.String() != "https://router.local/api/v4.0/network/qos/ip" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/network/qos/ip")
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			// writeCmd merges createDefaults into the body, so we check key fields.
			bs := string(body)
			for _, want := range []string{`"tagname":"office-cap"`, `"upload":2000`, `"download":2000`, `"interface":"wan1"`} {
				if !bytes.Contains(body, []byte(want)) {
					t.Fatalf("body missing %s: %s", want, bs)
				}
			}

			return jsonResponse(`{"code":0,"message":"created","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"ip", "create", "--name", "office-cap", "--interface", "wan1", "--upload", "2000", "--download", "2000"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := "{\"message\":\"created\"}\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestIPUpdateSendsFullBodyFromGet(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	var seenGet bool
	var seenPut bool
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-qos"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			switch req.Method {
			case http.MethodGet:
				seenGet = true
				if req.URL.String() != "https://router.local/api/v4.0/network/qos/ip/7" {
					t.Fatalf("GET URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/network/qos/ip/7")
				}
				return jsonResponse(`{"code":0,"data":{"id":7,"tagname":"old-qos","enabled":"yes","upload":100,"download":200,"type":0,"protocol":"any","comment":"","interface":"wan1","ip_addr":{"custom":["192.0.2.10"],"object":[]},"src_port":{"custom":[],"object":[]},"dst_port":{"custom":[],"object":[]},"attr":0}}`), nil
			case http.MethodPut:
				seenPut = true
				if req.URL.String() != "https://router.local/api/v4.0/network/qos/ip/7" {
					t.Fatalf("PUT URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/network/qos/ip/7")
				}
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("ReadAll() error = %v", err)
				}
				bs := string(body)
				for _, want := range []string{
					`"tagname":"new-qos"`,
					`"enabled":"yes"`,
					`"upload":100`,
					`"download":300`,
					`"interface":"wan1"`,
					`"protocol":"any"`,
					`"attr":0`,
					`"src_port":{"custom":[],"object":[]}`,
					`"dst_port":{"custom":[],"object":[]}`,
					`"ip_addr":{"custom":["192.0.2.10"],"object":[]}`,
				} {
					if !bytes.Contains(body, []byte(want)) {
						t.Fatalf("body missing %s: %s", want, bs)
					}
				}
				for _, forbidden := range []string{`"id":`} {
					if bytes.Contains(body, []byte(forbidden)) {
						t.Fatalf("body contains %s: %s", forbidden, bs)
					}
				}
				return jsonResponse(`{"code":0,"message":"updated","data":null}`), nil
			default:
				t.Fatalf("method = %q, want GET or PUT", req.Method)
			}
			return nil, nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"ip", "update", "7", "--name", "new-qos", "--download", "300"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !seenGet || !seenPut {
		t.Fatalf("seenGet=%v seenPut=%v, want both true", seenGet, seenPut)
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
