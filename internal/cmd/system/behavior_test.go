package system

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

func TestSchedulesListBuildsExpectedQueryParams(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-456"}
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.Path != "/api/v4.0/system/reboot-schedules" {
				t.Fatalf("path = %q, want %q", req.URL.Path, "/api/v4.0/system/reboot-schedules")
			}
			q := req.URL.Query()
			if q.Get("page") != "3" {
				t.Fatalf("page = %q, want %q", q.Get("page"), "3")
			}
			if q.Get("limit") != "10" {
				t.Fatalf("limit = %q, want %q", q.Get("limit"), "10")
			}
			if q.Get("filter") != "" {
				t.Fatalf("filter = %q, want empty", q.Get("filter"))
			}
			if q.Get("order") != "asc" {
				t.Fatalf("order = %q, want %q", q.Get("order"), "asc")
			}
			if q.Get("order_by") != "id" {
				t.Fatalf("order_by = %q, want %q", q.Get("order_by"), "id")
			}
			if got := req.Header.Get("Authorization"); got != "Bearer token-456" {
				t.Fatalf("Authorization = %q, want %q", got, "Bearer token-456")
			}
			return jsonResponse(`{"code":0,"data":{"items":[]}}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"schedules", "list", "--page", "3", "--page-size", "10", "--order", "asc", "--order-by", "id"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := `{"items":[]}` + "\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestSetSendsExpectedJSONBody(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-789"}
	calls := 0
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				if req.Method != http.MethodGet {
					t.Fatalf("method = %q, want %q", req.Method, http.MethodGet)
				}
				if req.URL.String() != "https://router.local/api/v4.0/system/basic/config" {
					t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/system/basic/config")
				}
				return jsonResponse(systemBasicBaselineResponse), nil
			}
			if req.Method != http.MethodPut {
				t.Fatalf("method = %q, want %q", req.Method, http.MethodPut)
			}
			if req.URL.String() != "https://router.local/api/v4.0/system/basic/config" {
				t.Fatalf("URL = %q, want %q", req.URL.String(), "https://router.local/api/v4.0/system/basic/config")
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			var got map[string]interface{}
			if err := json.Unmarshal(body, &got); err != nil {
				t.Fatalf("Unmarshal body error = %v", err)
			}
			if got["hostname"] != "ikuai-gw" {
				t.Fatalf("hostname = %v, want %q", got["hostname"], "ikuai-gw")
			}
			for _, key := range systemBasicFields {
				if _, ok := got[key]; !ok {
					t.Fatalf("body missing full-update field %q: %s", key, string(body))
				}
			}

			return jsonResponse(`{"code":0,"message":"saved","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"set", "--data", `{"hostname":"ikuai-gw"}`})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	want := "{\"message\":\"saved\"}\n"
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestSetSemanticFlagsSendCorrectBody(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-sf"}
	calls := 0
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				return jsonResponse(systemBasicBaselineResponse), nil
			}
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			var m map[string]interface{}
			if err := json.Unmarshal(body, &m); err != nil {
				t.Fatalf("Unmarshal body error = %v", err)
			}
			if m["hostname"] != "my-router" {
				t.Fatalf("hostname = %v, want %q", m["hostname"], "my-router")
			}
			if got := m["language"]; got != "2" && got != float64(2) {
				t.Fatalf("language = %v, want %v", got, float64(2))
			}
			return jsonResponse(`{"code":0,"message":"saved","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"set", "--hostname", "my-router", "--language", "2"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestSetSemanticFlagsOverrideData(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON
	app.Session = &session.Session{BaseURL: "https://router.local", Token: "token-ov"}
	calls := 0
	app.APIClient = api.NewWithHTTPClient(app.Session.BaseURL, app.Session.Token, &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				return jsonResponse(systemBasicBaselineResponse), nil
			}
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			var m map[string]interface{}
			if err := json.Unmarshal(body, &m); err != nil {
				t.Fatalf("Unmarshal body error = %v", err)
			}
			// --hostname flag should override the data value
			if m["hostname"] != "flag-wins" {
				t.Fatalf("hostname = %v, want %q", m["hostname"], "flag-wins")
			}
			// language from --data should be preserved
			if m["language"] != "cn" {
				t.Fatalf("language = %v, want %q", m["language"], "cn")
			}
			return jsonResponse(`{"code":0,"message":"saved","data":null}`), nil
		}),
	})

	cmd := New(app)
	cmd.SetArgs([]string{"set", "--data", `{"hostname":"data-value","language":"cn"}`, "--hostname", "flag-wins"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

const systemBasicBaselineResponse = `{
  "code": 0,
  "message": "Success",
  "results": {
    "total": 1,
    "data": [
      {
        "id": 1,
        "hostname": "Router",
        "language": 1,
        "time_zone": 8,
        "time_zone_full": "0800",
        "switch_nat": 1,
        "switch_ntp": 1,
        "switch_ntpd": 0,
        "switch_ntpserver": 0,
        "ntpserver_list": "",
        "ntp_sync_cycle": 60,
        "link_mode": 0,
        "lan_nat": 1,
        "backport": "wan1",
        "listenport": "lan1",
        "fast_nat": 0
      }
    ]
  }
}`

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
