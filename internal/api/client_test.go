package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func fakeResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func testCheck(client *Client, status int, body string) (json.RawMessage, error) {
	resp := fakeResponse(status, body)
	defer func() { _ = resp.Body.Close() }()
	return client.check(resp)
}

func TestCheckReturnsDataFieldOnly(t *testing.T) {
	client := New("https://192.168.1.1", "")
	raw, err := testCheck(client, 200, `{"code":0,"message":"ok","data":{"uptime":123}}`)
	if err != nil {
		t.Fatalf("check() error = %v", err)
	}
	want := `{"uptime":123}`
	if string(raw) != want {
		t.Fatalf("check() = %s, want %s", string(raw), want)
	}
}

func TestCheckRawModeReturnsFullEnvelope(t *testing.T) {
	client := New("https://192.168.1.1", "")
	client.RawMode = true
	raw, err := testCheck(client, 200, `{"code":0,"message":"ok","data":{"uptime":123}}`)
	if err != nil {
		t.Fatalf("check() error = %v", err)
	}
	if !strings.Contains(string(raw), `"code":0`) {
		t.Fatalf("RawMode check() should contain code field: %s", string(raw))
	}
	if !strings.Contains(string(raw), `"data":`) {
		t.Fatalf("RawMode check() should contain data field: %s", string(raw))
	}
}

func TestCheckReturnsNullDataSynthesizesMessage(t *testing.T) {
	client := New("https://192.168.1.1", "")
	raw, err := testCheck(client, 200, `{"code":0,"message":"ok","data":null}`)
	if err != nil {
		t.Fatalf("check() error = %v", err)
	}
	want := `{"message":"ok"}`
	if string(raw) != want {
		t.Fatalf("check() = %s, want %s", string(raw), want)
	}
}

func TestCheckReturnsAbsentDataSynthesizesMessage(t *testing.T) {
	client := New("https://192.168.1.1", "")
	raw, err := testCheck(client, 200, `{"code":0,"message":"saved"}`)
	if err != nil {
		t.Fatalf("check() error = %v", err)
	}
	want := `{"message":"saved"}`
	if string(raw) != want {
		t.Fatalf("check() = %s, want %s", string(raw), want)
	}
}

func TestCheckNormalizesNilToNull(t *testing.T) {
	client := New("https://192.168.1.1", "")
	raw, err := testCheck(client, 200, `{"code":0,"data":{"foo":nil}}`)
	if err != nil {
		t.Fatalf("check() error = %v", err)
	}
	if string(raw) != `{"foo":null}` {
		t.Fatalf("check() = %s, want nil normalized to null in data", string(raw))
	}
}

func TestCheckErrorWithHint3001(t *testing.T) {
	client := New("https://192.168.1.1", "")
	_, err := testCheck(client, 200, `{"code":3001,"message":"param error"}`)
	if err == nil {
		t.Fatal("expected error for code 3001")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Code != 3001 {
		t.Fatalf("Code = %d, want 3001", apiErr.Code)
	}
	if !strings.Contains(apiErr.Message, "Parameter error") {
		t.Fatalf("Message = %q, want hint about parameter error", apiErr.Message)
	}
}

func TestCheckErrorWithHint3007(t *testing.T) {
	client := New("https://192.168.1.1", "")
	_, err := testCheck(client, 200, `{"code":3007,"message":"Invalid token"}`)
	if err == nil {
		t.Fatal("expected error for code 3007")
	}
	apiErr := err.(*APIError)
	if !strings.Contains(apiErr.Message, "auth set-url") {
		t.Fatalf("Message = %q, want hint about auth set-url", apiErr.Message)
	}
}

func TestCheckErrorUnknownCode(t *testing.T) {
	client := New("https://192.168.1.1", "")
	_, err := testCheck(client, 200, `{"code":9999,"message":"unknown"}`)
	if err == nil {
		t.Fatal("expected error for code 9999")
	}
	apiErr := err.(*APIError)
	if apiErr.Message != "unknown" {
		t.Fatalf("Message = %q, want no hint appended", apiErr.Message)
	}
}

func TestCheckErrorWithHint1008(t *testing.T) {
	client := New("https://192.168.1.1", "")
	_, err := testCheck(client, 200, `{"code":1008,"message":"Session expired"}`)
	if err == nil {
		t.Fatal("expected error for code 1008")
	}
	apiErr := err.(*APIError)
	if !strings.Contains(apiErr.Message, "auth set-url") {
		t.Fatalf("Message = %q, want hint about auth set-url", apiErr.Message)
	}
}

func TestCheckCode20000IsSuccess(t *testing.T) {
	client := New("https://192.168.1.1", "")
	raw, err := testCheck(client, 200, `{"code":20000,"message":"ok","data":{"id":1}}`)
	if err != nil {
		t.Fatalf("check() error = %v", err)
	}
	want := `{"id":1}`
	if string(raw) != want {
		t.Fatalf("check() = %s, want %s", string(raw), want)
	}
}

func TestCheckNonJSONResponse(t *testing.T) {
	client := New("https://192.168.1.1", "")
	_, err := testCheck(client, 200, `<html>Not Found</html>`)
	if err == nil {
		t.Fatal("expected error for non-JSON response")
	}
	if !strings.Contains(err.Error(), "non-JSON") {
		t.Fatalf("error = %q, want 'non-JSON response'", err.Error())
	}
}

func TestCheckEmptyMessageFallsBackToBody(t *testing.T) {
	client := New("https://192.168.1.1", "")
	_, err := testCheck(client, 200, `{"code":5000,"message":""}`)
	if err == nil {
		t.Fatal("expected error for code 5000")
	}
	apiErr := err.(*APIError)
	// When message is empty, the full body is used as the error message.
	if !strings.Contains(apiErr.Message, "5000") {
		t.Fatalf("Message = %q, want body content", apiErr.Message)
	}
}

func TestSanitizeNilPreservesStringValues(t *testing.T) {
	// The word "nil" inside a JSON string value should NOT be replaced.
	input := `{"code":0,"message":"set to nil","data":{"comment":"nil value"}}`
	client := New("https://192.168.1.1", "")
	raw, err := testCheck(client, 200, input)
	if err != nil {
		t.Fatalf("check() error = %v", err)
	}
	got := string(raw)
	if !strings.Contains(got, "nil value") {
		t.Fatalf("sanitizeNil corrupted string value: %s", got)
	}
}

func TestCheckReturnsResultsWhenDataAbsent(t *testing.T) {
	client := New("https://192.168.1.1", "")
	raw, err := testCheck(client, 200, `{"code":0,"message":"Success","results":{"cpu":[{"cpu":0.5}]}}`)
	if err != nil {
		t.Fatalf("check() error = %v", err)
	}
	want := `{"cpu":[{"cpu":0.5}]}`
	if string(raw) != want {
		t.Fatalf("check() = %s, want %s", string(raw), want)
	}
}

func TestCheckPrefersDataOverResults(t *testing.T) {
	client := New("https://192.168.1.1", "")
	raw, err := testCheck(client, 200, `{"code":0,"message":"ok","data":{"id":1},"results":{"cpu":[]}}`)
	if err != nil {
		t.Fatalf("check() error = %v", err)
	}
	want := `{"id":1}`
	if string(raw) != want {
		t.Fatalf("check() = %s, want %s (should prefer data over results)", string(raw), want)
	}
}

func TestCheckRowIDReturned(t *testing.T) {
	client := New("https://192.168.1.1", "")
	raw, err := testCheck(client, 200, `{"code":0,"message":"success","rowid":42}`)
	if err != nil {
		t.Fatalf("check() error = %v", err)
	}
	want := `{"message":"success","rowid":42}`
	if string(raw) != want {
		t.Fatalf("check() = %s, want %s", string(raw), want)
	}
}

func TestCheckRowIDWithoutCode(t *testing.T) {
	client := New("https://192.168.1.1", "")
	raw, err := testCheck(client, 200, `{"message":"success","rowid":1}`)
	if err != nil {
		t.Fatalf("check() error = %v", err)
	}
	want := `{"message":"success","rowid":1}`
	if string(raw) != want {
		t.Fatalf("check() = %s, want %s", string(raw), want)
	}
}

func TestCheckDataTakesPriorityOverRowID(t *testing.T) {
	client := New("https://192.168.1.1", "")
	raw, err := testCheck(client, 200, `{"code":0,"message":"ok","data":{"id":5},"rowid":5}`)
	if err != nil {
		t.Fatalf("check() error = %v", err)
	}
	want := `{"id":5}`
	if string(raw) != want {
		t.Fatalf("check() = %s, want data to take priority over rowid", string(raw))
	}
}

func TestCheckResultsTakesPriorityOverRowID(t *testing.T) {
	client := New("https://192.168.1.1", "")
	raw, err := testCheck(client, 200, `{"code":0,"message":"ok","results":{"cpu":[]},"rowid":5}`)
	if err != nil {
		t.Fatalf("check() error = %v", err)
	}
	want := `{"cpu":[]}`
	if string(raw) != want {
		t.Fatalf("check() = %s, want results to take priority over rowid", string(raw))
	}
}

func TestCheckRowIDAsString(t *testing.T) {
	client := New("https://192.168.1.1", "")
	raw, err := testCheck(client, 200, `{"code":0,"message":"success","rowid":"42"}`)
	if err != nil {
		t.Fatalf("check() error = %v", err)
	}
	got := string(raw)
	if !strings.Contains(got, "rowid") {
		t.Fatalf("string rowid lost: %s", got)
	}
	if !strings.Contains(got, "42") {
		t.Fatalf("rowid value lost: %s", got)
	}
}

func TestCheckRowIDRawModeReturnsFullBody(t *testing.T) {
	client := New("https://192.168.1.1", "")
	client.RawMode = true
	raw, err := testCheck(client, 200, `{"code":0,"message":"success","rowid":42}`)
	if err != nil {
		t.Fatalf("check() error = %v", err)
	}
	if !strings.Contains(string(raw), `"rowid":42`) {
		t.Fatalf("RawMode should return full body with rowid: %s", string(raw))
	}
}

func TestCheckHTTP400(t *testing.T) {
	client := New("https://192.168.1.1", "")
	_, err := testCheck(client, 400, `{"code":3001,"message":"bad request"}`)
	if err == nil {
		t.Fatal("expected error for HTTP 400")
	}
}

func TestDryRunPostReturnsSyntheticJSON(t *testing.T) {
	client := New("https://192.168.1.1", "tok")
	client.DryRun = true
	raw, err := client.Post("/api/v4.0/system/basic/config", map[string]string{"hostname": "test"})
	if err != nil {
		t.Fatalf("Post() error = %v", err)
	}
	got := string(raw)
	if !strings.Contains(got, `"dry_run":true`) {
		t.Fatalf("missing dry_run field: %s", got)
	}
	if !strings.Contains(got, `"method":"POST"`) {
		t.Fatalf("missing method field: %s", got)
	}
	if !strings.Contains(got, `"hostname":"test"`) {
		t.Fatalf("missing body content: %s", got)
	}
}

func TestDryRunDeleteOmitsBody(t *testing.T) {
	client := New("https://192.168.1.1", "tok")
	client.DryRun = true
	raw, err := client.Delete("/api/v4.0/network/dhcp/services/1")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	got := string(raw)
	if !strings.Contains(got, `"dry_run":true`) {
		t.Fatalf("missing dry_run field: %s", got)
	}
	if strings.Contains(got, `"body"`) {
		t.Fatalf("DELETE dry-run should not include body field: %s", got)
	}
}

func TestDryRunGetReturnsSyntheticJSON(t *testing.T) {
	client := New("https://192.168.1.1", "tok")
	client.DryRun = true
	raw, err := client.Get("/api/v4.0/system/basic/config", map[string]string{"page": "1"})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	got := string(raw)
	if !strings.Contains(got, `"dry_run":true`) {
		t.Fatalf("missing dry_run field: %s", got)
	}
	if !strings.Contains(got, `"method":"GET"`) {
		t.Fatalf("missing method field: %s", got)
	}
	if !strings.Contains(got, "page=1") {
		t.Fatalf("missing query params in URL: %s", got)
	}
}
