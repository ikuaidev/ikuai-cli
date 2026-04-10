package advanced

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
			if q.Get("page") != "2" || q.Get("page_size") != "8" || q.Get("filter") != "enabled==yes" || q.Get("order") != "asc" || q.Get("order_by") != "username" {
				t.Fatalf("unexpected query params: %s", req.URL.RawQuery)
			}
			if got := req.Header.Get("Authorization"); got != "Bearer token-adv" {
				t.Fatalf("Authorization = %q, want %q", got, "Bearer token-adv")
			}
			return jsonResponse(`{"code":0,"data":{"items":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"ftp", "list", "--page", "2", "--page-size", "8", "--filter", "enabled==yes", "--order", "asc", "--order-by", "username"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := `{"items":[]}` + "\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestSNMPDSetSendsExpectedJSONBody(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-adv"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
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
			if string(body) != `{"enable":"yes"}` {
				t.Fatalf("body = %q, want %q", string(body), `{"enable":"yes"}`)
			}

			return jsonResponse(`{"code":0,"message":"saved","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"snmpd", "set", "--data", `{"enable":"yes"}`})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
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
