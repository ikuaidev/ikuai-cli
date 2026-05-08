package advanced

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/ikuaidev/ikuai-cli/internal/api"
	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/ikuaidev/ikuai-cli/internal/output"
	"github.com/ikuaidev/ikuai-cli/internal/session"
)

func TestFTPListBuildsExpectedQueryParams(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-adv"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/advanced-service/ftp-users" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/advanced-service/ftp-users")
			}
			q := req.URL.Query()
			if q.Get("page") != "2" || q.Get("limit") != "8" || q.Get("key") != "username" || q.Get("pattern") != "alice" || q.Get("order") != "username" || q.Get("order_by") != "asc" {
				t.Fatalf("unexpected query params: %s", req.URL.RawQuery)
			}
			if q.Has("page_size") || q.Has("filter") {
				t.Fatalf("unsupported query params sent: %s", req.URL.RawQuery)
			}
			if got := req.Header.Get("Authorization"); got != "Bearer token-adv" {
				t.Fatalf("Authorization = %q, want %q", got, "Bearer token-adv")
			}
			return jsonResponse(`{"code":0,"data":{"items":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"ftp", "list", "--page", "2", "--page-size", "8", "--key", "username", "--pattern", "alice", "--order", "asc", "--order-by", "username"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := `{"items":[]}` + "\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestFTPCreateDoesNotExposeNameFlag(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	cmd := New(app)
	cmd.SetArgs([]string{"ftp", "create", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if bytes.Contains(out.Bytes(), []byte("--name")) {
		t.Fatalf("ftp create help exposes YAML-unsupported --name flag:\n%s", out.String())
	}
}

func TestFTPUpdateReadsBaselineFromListBeforePut(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-adv"}
	seenList := false
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			switch req.Method {
			case http.MethodGet:
				seenList = true
				if req.URL.Path != "/api/v4.0/advanced-service/ftp-users" {
					t.Fatalf("GET path = %q, want list path", req.URL.Path)
				}
				if req.URL.Query().Get("limit") != "500" {
					t.Fatalf("GET query = %s", req.URL.RawQuery)
				}
				return jsonResponse(`{"code":0,"message":"Success","results":{"total":1,"data":[{"id":7,"enabled":"yes","username":"olduser","tagname":"olduser","passwd":"oldpass","permission":"rw","home_dir":"/test","upload":0,"download":0}]}}`), nil
			case http.MethodPut:
				if req.URL.Path != "/api/v4.0/advanced-service/ftp-users/7" {
					t.Fatalf("PUT path = %q", req.URL.Path)
				}
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("ReadAll() error = %v", err)
				}
				var got map[string]interface{}
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatalf("Unmarshal body: %v", err)
				}
				for _, key := range []string{"enabled", "username", "passwd", "permission", "home_dir", "upload", "download", "tagname"} {
					if _, ok := got[key]; !ok {
						t.Fatalf("body missing %q: %s", key, string(body))
					}
				}
				if got["permission"] != "ro" || got["username"] != "olduser" {
					t.Fatalf("unexpected update body: %s", string(body))
				}
				return jsonResponse(`{"code":0,"message":"saved","data":null}`), nil
			default:
				t.Fatalf("unexpected method %s", req.Method)
			}
			return nil, nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"ftp", "update", "7", "--permission", "ro"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !seenList {
		t.Fatal("expected update to read list baseline before PUT")
	}
}

func TestSNMPDSetSendsExpectedJSONBody(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-adv"}
	seenGet := false
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method == http.MethodGet {
				seenGet = true
				if req.URL.String() != "https://router.local/api/v4.0/advanced-service/snmpd-config" {
					t.Fatalf("GET URL = %q", req.URL.String())
				}
				return jsonResponse(`{"code":0,"message":"Success","results":{"data":[{"id":1,"enabled":"no","listen_port":161,"syslocation":"","syscontact":"","sysname":"","version":2,"community":"public","source":"","rw":"ro","username":"","security":"authNoPriv","auth_proto":"","auth_pass":"","priv_proto":"","priv_pass":""}]}}`), nil
			}
			if req.Method != http.MethodPut {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodPut)
			}
			if req.URL.String() != "https://router.local/api/v4.0/advanced-service/snmpd-config" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/advanced-service/snmpd-config")
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			var got map[string]interface{}
			if err := json.Unmarshal(body, &got); err != nil {
				t.Fatalf("Unmarshal body: %v", err)
			}
			for _, key := range []string{"enabled", "listen_port", "syslocation", "syscontact", "sysname", "version", "community", "source", "rw", "username", "security", "auth_proto", "auth_pass", "priv_proto", "priv_pass"} {
				if _, ok := got[key]; !ok {
					t.Fatalf("body missing %q: %s", key, string(body))
				}
			}
			if got["enabled"] != "yes" || got["listen_port"] != float64(161) {
				t.Fatalf("unexpected body: %s", string(body))
			}

			return jsonResponse(`{"code":0,"message":"saved","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"snmpd", "set", "--enabled", "yes"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !seenGet {
		t.Fatal("expected set to read baseline before PUT")
	}

	got := out.String()
	want := "{\"message\":\"saved\"}\n"
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
