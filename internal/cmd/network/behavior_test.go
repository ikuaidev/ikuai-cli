package network

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
			if q.Get("limit") != "50" {
				t.Fatalf("limit = %q, want %q", q.Get("limit"), "50")
			}
			if q.Has("page_size") || q.Has("filter") || q.Has("order") || q.Has("order_by") {
				t.Fatalf("unsupported query params sent: %s", req.URL.RawQuery)
			}
			if got := req.Header.Get("Authorization"); got != "Bearer token-abc" {
				t.Fatalf("Authorization = %q, want %q", got, "Bearer token-abc")
			}
			return jsonResponse(`{"code":0,"data":{"items":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"nat", "list", "--page", "2", "--page-size", "50"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := `{"items":[]}` + "\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestWANListUsesConfigEndpointWithoutPagination(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-wan"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/interfaces/wan-config" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/interfaces/wan-config")
			}
			if req.URL.RawQuery != "" {
				t.Fatalf("query = %q, want empty", req.URL.RawQuery)
			}
			return jsonResponse(`{"code":0,"data":[{"id":1,"tagname":"wan1","internet":3,"mtu":1480}],"total":1}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"wan", "list"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := `[{"id":1,"internet":3,"mtu":1480,"tagname":"wan1"}]` + "\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestNetworkListsUseLimitQueryParam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		path string
	}{
		{name: "DNS proxy", args: []string{"dns", "proxy", "list"}, path: "/api/v4.0/network/dns/proxy/rules"},
		{name: "DHCP services", args: []string{"dhcp", "list"}, path: "/api/v4.0/network/dhcp/services"},
		{name: "DHCPv4", args: []string{"dhcp", "clients"}, path: "/api/v4.0/network/dhcp/clients"},
		{name: "DHCPv6", args: []string{"dhcp6", "clients"}, path: "/api/v4.0/network/dhcp6/clients"},
		{name: "DHCP static", args: []string{"dhcp", "static", "list"}, path: "/api/v4.0/network/dhcp/static"},
		{name: "DHCP access", args: []string{"dhcp", "access-rule", "list"}, path: "/api/v4.0/network/dhcp/access-control/rules"},
		{name: "DHCPv6 access", args: []string{"dhcp6", "access-rule", "list"}, path: "/api/v4.0/network/dhcp6/access-control/rules"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var out bytes.Buffer
			app := cliapp.New(&out, &out)
			app.Format = output.JSON
			app.Session = &session.Session{BaseURL: "https://router.local", Token: "test-token"}
			app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					if req.Method != http.MethodGet {
						t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
					}
					if req.URL.Path != tt.path {
						t.Fatalf("path = %q, want %q", req.URL.Path, tt.path)
					}
					query := req.URL.Query()
					if got := query.Get("page"); got != "2" {
						t.Fatalf("page = %q, want %q", got, "2")
					}
					if got := query.Get("limit"); got != "50" {
						t.Fatalf("limit = %q, want %q", got, "50")
					}
					if query.Has("page_size") {
						t.Fatalf("page_size should not be sent: %s", req.URL.RawQuery)
					}
					return jsonResponse(`{"code":0,"results":{"total":0,"data":[]}}`), nil
				}),
			})

			cmd := New(app)
			cmd.SetArgs(append(tt.args, "--page", "2", "--page-size", "50"))
			if err := cmd.Execute(); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
		})
	}
}

func TestVLANListUsesYAMLQueryParams(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "test-token"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			query := req.URL.Query()
			want := map[string]string{
				"page": "2", "limit": "50", "key": "vlan_name",
				"pattern": "office", "filter": "enabled==yes",
			}
			for key, value := range want {
				if got := query.Get(key); got != value {
					t.Fatalf("%s = %q, want %q", key, got, value)
				}
			}
			for _, unsupported := range []string{"page_size", "order", "order_by"} {
				if query.Has(unsupported) {
					t.Fatalf("unsupported %s query sent: %s", unsupported, req.URL.RawQuery)
				}
			}
			return jsonResponse(`{"code":0,"results":{"total":0,"data":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"vlan", "list", "--page", "2", "--page-size", "50", "--key", "vlan_name", "--pattern", "office", "--filter", "enabled==yes"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestDHCPCreateMissingRequiredFlags(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "tok"}

	cmd := New(app)
	cmd.SetArgs([]string{"dhcp", "create", "--name", "test"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing required flags")
	}
	if !strings.Contains(err.Error(), "missing required flags") {
		t.Fatalf("error = %q, want it to contain 'missing required flags'", err.Error())
	}
}

func TestDHCPToggleMissingEnabled(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "tok"}

	cmd := New(app)
	cmd.SetArgs([]string{"dhcp", "toggle", "1"})
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
