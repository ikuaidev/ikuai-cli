package monitor

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

func TestLoadCommandPassesYamlTimeParameters(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-123"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/api/v4.0/monitoring/cpu" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/monitoring/cpu")
			}
			query := req.URL.Query()
			for name, want := range map[string]string{
				"datetype":   "day",
				"start_time": "1773300000",
				"end_time":   "1773303600",
				"math":       "max",
			} {
				if got := query.Get(name); got != want {
					t.Fatalf("query %s = %q, want %q", name, got, want)
				}
			}
			return jsonResponse(`{"code":0,"message":"ok","results":{"cpu":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"cpu", "--time-range", "day", "--start-time", "1773300000", "--end-time", "1773303600", "--aggregate", "max"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestLoadCommandRequiresYamlRequiredTimeParameters(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-123"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			t.Fatalf("unexpected HTTP request: %s", req.URL.String())
			return nil, nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"cpu"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want missing required flags")
	}
	want := "missing required flags: --time-range, --start-time, --end-time, --aggregate"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

func TestProtocolsPassesYamlTimeParameters(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-123"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/api/v4.0/monitoring/protocols" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/monitoring/protocols")
			}
			query := req.URL.Query()
			if got := query.Get("starttime"); got != "1773215100" {
				t.Fatalf("starttime = %q, want %q", got, "1773215100")
			}
			if got := query.Get("stoptime"); got != "1773218700" {
				t.Fatalf("stoptime = %q, want %q", got, "1773218700")
			}
			return jsonResponse(`{"code":0,"message":"ok","results":{"data":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"protocols", "--starttime", "1773215100", "--stoptime", "1773218700"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestClientProtocolsHistoryPassesYamlTimeQueryParameters(t *testing.T) {
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
			if req.URL.Path != "/api/v4.0/monitoring/clients/protocols/history-load" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/monitoring/clients/protocols/history-load")
			}
			query := req.URL.Query()
			for name, want := range map[string]string{
				"ip":        "192.168.9.199",
				"mac":       "08:9b:4b:01:7e:7c",
				"starttime": "1773304236",
				"stoptime":  "1773304246",
			} {
				if got := query.Get(name); got != want {
					t.Fatalf("query %s = %q, want %q", name, got, want)
				}
			}
			return jsonResponse(`{"code":0,"message":"ok","results":{"data":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"client-protocols-history", "--ip", "192.168.9.199", "--mac", "08:9b:4b:01:7e:7c", "--starttime", "1773304236", "--stoptime", "1773304246"})
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

func TestAppTrafficSummaryOnlyExposesYamlFlags(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	cmd := New(app)
	appTrafficCmd, _, err := cmd.Find([]string{"app-traffic-summary"})
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	for _, name := range []string{"page", "page-size"} {
		if appTrafficCmd.Flags().Lookup(name) == nil {
			t.Fatalf("expected flag %q to exist", name)
		}
	}
	for _, name := range []string{"filter", "order", "order-by"} {
		if appTrafficCmd.Flags().Lookup(name) != nil {
			t.Fatalf("flag %q should not be exposed by app-traffic-summary", name)
		}
	}
}

func TestClientListPassesYamlSearchParameters(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-123"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/api/v4.0/monitoring/clients-online" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/monitoring/clients-online")
			}
			query := req.URL.Query()
			for name, want := range map[string]string{
				"key":     "mac,ip_addr",
				"pattern": "08:9b",
				"filter":  "interface==lan1",
			} {
				if got := query.Get(name); got != want {
					t.Fatalf("query %s = %q, want %q", name, got, want)
				}
			}
			return jsonResponse(`{"code":0,"message":"ok","results":{"data":[],"total":0}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"clients-online", "--key", "mac,ip_addr", "--pattern", "08:9b", "--filter", "interface==lan1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestClientListOnlyExposesTrimmedYamlFlags(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	cmd := New(app)
	clientsCmd, _, err := cmd.Find([]string{"clients-online"})
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	for _, name := range []string{"page", "page-size", "filter", "key", "pattern"} {
		if clientsCmd.Flags().Lookup(name) == nil {
			t.Fatalf("expected flag %q to exist", name)
		}
	}
	for _, name := range []string{"order", "order-by"} {
		if clientsCmd.Flags().Lookup(name) != nil {
			t.Fatalf("flag %q should not be exposed by clients-online", name)
		}
	}
}

func TestTrafficSummaryMapsPageSizeToYamlLimit(t *testing.T) {
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
			if req.URL.Path != "/api/v4.0/monitoring/clients-traffic-summary" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/monitoring/clients-traffic-summary")
			}
			query := req.URL.Query()
			if got := query.Get("page"); got != "2" {
				t.Fatalf("page = %q, want %q", got, "2")
			}
			if got := query.Get("limit"); got != "50" {
				t.Fatalf("limit = %q, want %q", got, "50")
			}
			for _, name := range []string{"page_size", "filter", "order", "order_by"} {
				if got := query.Get(name); got != "" {
					t.Fatalf("%s = %q, want empty", name, got)
				}
			}
			return jsonResponse(`{"code":0,"message":"ok","results":{"terminal":[],"terminal_total":0}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"traffic-summary", "--page", "2", "--page-size", "50"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestWirelessCommandsPassYamlFilterParameters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		args      []string
		path      string
		queryName string
		want      string
		body      string
	}{
		{
			name:      "wireless-traffic apmac",
			args:      []string{"wireless-traffic", "--apmac", "00:00:00:00:00:00"},
			path:      "/api/v4.0/monitoring/wireless-traffic",
			queryName: "apmac",
			want:      "00:00:00:00:00:00",
			body:      `{"code":0,"message":"ok","results":{"total_count_flow":[]}}`,
		},
		{
			name:      "ssid-clients ssid",
			args:      []string{"ssid-clients", "--ssid", "AP_test001"},
			path:      "/api/v4.0/monitoring/ssid-clients",
			queryName: "ssid",
			want:      "AP_test001",
			body:      `{"code":0,"message":"ok","results":{"ssid_sta_history":[]}}`,
		},
		{
			name:      "channel-clients channel",
			args:      []string{"channel-clients", "--channel", "3"},
			path:      "/api/v4.0/monitoring/channel-clients",
			queryName: "channel",
			want:      "3",
			body:      `{"code":0,"message":"ok","results":{"channel_sta_history":[]}}`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var out bytes.Buffer
			app := cliapp.New(&out, &out)
			app.Format = output.JSON
			app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-123"}
			app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					if req.URL.Path != tt.path {
						t.Fatalf("path = %q, want %q", req.URL.Path, tt.path)
					}
					if got := req.URL.Query().Get(tt.queryName); got != tt.want {
						t.Fatalf("query %s = %q, want %q", tt.queryName, got, tt.want)
					}
					return jsonResponse(tt.body), nil
				}),
			})
			cmd := New(app)
			cmd.SetArgs(tt.args)
			if err := cmd.Execute(); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
		})
	}
}

func TestAppProtocolsLoadMapsPageSizeToYamlLimit(t *testing.T) {
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
			if req.URL.Path != "/api/v4.0/monitoring/app-protocols/load" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/monitoring/app-protocols/load")
			}
			query := req.URL.Query()
			for name, want := range map[string]string{
				"page":     "3",
				"limit":    "25",
				"order":    "desc",
				"order_by": "total_down",
			} {
				if got := query.Get(name); got != want {
					t.Fatalf("query %s = %q, want %q", name, got, want)
				}
			}
			if got := query.Get("page_size"); got != "" {
				t.Fatalf("page_size = %q, want empty", got)
			}
			return jsonResponse(`{"code":0,"message":"ok","results":{"data":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"app-protocols-load", "--page", "3", "--page-size", "25", "--order", "desc", "--order-by", "total_down"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestClientAppProtocolsPassesYamlLimit(t *testing.T) {
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
			if req.URL.Path != "/api/v4.0/monitoring/clients/app-protocols/load" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/monitoring/clients/app-protocols/load")
			}
			query := req.URL.Query()
			for name, want := range map[string]string{
				"ip":    "192.168.9.199",
				"mac":   "08:9b:4b:01:7e:7c",
				"limit": "10",
			} {
				if got := query.Get(name); got != want {
					t.Fatalf("query %s = %q, want %q", name, got, want)
				}
			}
			return jsonResponse(`{"code":0,"message":"ok","results":{"data":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"client-app-protocols", "--ip", "192.168.9.199", "--mac", "08:9b:4b:01:7e:7c", "--page-size", "10"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestClientAppProtocolsDefaultTableUsesReturnedAppnameField(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.Table
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-123"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/api/v4.0/monitoring/clients/app-protocols/load" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/monitoring/clients/app-protocols/load")
			}
			return jsonResponse(`{"code":0,"message":"ok","data":[{"id":1,"appid":12345,"appname":"ChatGPT","conn_cnt":2,"upload":3,"download":4,"total":7}]}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"client-app-protocols", "--ip", "192.168.9.199", "--mac", "08:9b:4b:01:7e:7c", "--page-size", "10"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "APPID") || !strings.Contains(got, "12345") {
		t.Fatalf("default table should show appid value: %q", got)
	}
	if !strings.Contains(got, "APPNAME") || !strings.Contains(got, "ChatGPT") {
		t.Fatalf("default table should show appname value: %q", got)
	}
	if strings.Contains(got, "APP_NAME") {
		t.Fatalf("default table should not use stale app_name column: %q", got)
	}
}

func TestClientProtocolsPassesYamlTimeQueryParameters(t *testing.T) {
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
			if req.URL.Path != "/api/v4.0/monitoring/clients/protocols" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/monitoring/clients/protocols")
			}
			query := req.URL.Query()
			for name, want := range map[string]string{
				"ip":        "192.168.9.199",
				"mac":       "08:9b:4b:01:7e:7c",
				"starttime": "1773304236",
				"stoptime":  "1773304246",
			} {
				if got := query.Get(name); got != want {
					t.Fatalf("query %s = %q, want %q", name, got, want)
				}
			}
			return jsonResponse(`{"code":0,"message":"ok","results":{"data":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"client-protocols", "--ip", "192.168.9.199", "--mac", "08:9b:4b:01:7e:7c", "--starttime", "1773304236", "--stoptime", "1773304246"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestAppProtocolsTerminalsUsesIntegerAppIDAndCleanRequiredHelp(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-123"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/api/v4.0/monitoring/app-protocols/terminal-load" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/monitoring/app-protocols/terminal-load")
			}
			if got := req.URL.Query().Get("appid"); got != "2580003" {
				t.Fatalf("appid = %q, want %q", got, "2580003")
			}
			return jsonResponse(`{"code":0,"message":"ok","results":{"data":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"app-protocols-terminals", "--appid", "2580003"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	terminalsCmd, _, err := cmd.Find([]string{"app-protocols-terminals"})
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if got := terminalsCmd.Flags().Lookup("appid").Usage; got != "App protocol ID, e.g. 2580003 (required)" {
		t.Fatalf("appid usage = %q", got)
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
