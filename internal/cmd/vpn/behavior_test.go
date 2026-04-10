package vpn

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

func TestOpenVPNClientsBuildExpectedQueryParams(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-vpn"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/vpn/openvpn/clients" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/vpn/openvpn/clients")
			}
			q := req.URL.Query()
			if q.Get("page") != "2" {
				t.Fatalf("page = %q, want %q", q.Get("page"), "2")
			}
			if q.Get("page_size") != "16" {
				t.Fatalf("page_size = %q, want %q", q.Get("page_size"), "16")
			}
			if q.Get("filter") != "enabled==yes" {
				t.Fatalf("filter = %q, want %q", q.Get("filter"), "enabled==yes")
			}
			if q.Get("order") != "asc" {
				t.Fatalf("order = %q, want %q", q.Get("order"), "asc")
			}
			if q.Get("order_by") != "username" {
				t.Fatalf("order_by = %q, want %q", q.Get("order_by"), "username")
			}
			if got := req.Header.Get("Authorization"); got != "Bearer token-vpn" {
				t.Fatalf("Authorization = %q, want %q", got, "Bearer token-vpn")
			}
			return jsonResponse(`{"code":0,"data":{"items":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"openvpn", "clients", "--page", "2", "--page-size", "16", "--filter", "enabled==yes", "--order", "asc", "--order-by", "username"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := `{"items":[]}` + "\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestWireGuardPeerCreateSendsExpectedJSONBody(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-vpn"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodPost)
			}
			if req.URL.String() != "https://router.local/api/v4.0/vpn/wireguard/tunnel-1/peers" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/vpn/wireguard/tunnel-1/peers")
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			// writeCmd merges createDefaults into the body, so we check key fields.
			bs := string(body)
			for _, want := range []string{`"name":"peer-a"`, `"public_key":"pub-key"`} {
				if !bytes.Contains(body, []byte(want)) {
					t.Fatalf("body missing %s: %s", want, bs)
				}
			}

			return jsonResponse(`{"code":0,"message":"created","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"wireguard", "peer-create", "tunnel-1", "--data", `{"name":"peer-a","public_key":"pub-key"}`})
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
