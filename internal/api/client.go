// Package api provides the HTTP client for the iKuai router local REST API.
package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const timeout = 15 * time.Second

// sanitizeNil replaces bare `nil` tokens in JSON value positions with `null`.
// Some iKuai firmware emits `nil` instead of `null`. This implementation
// tracks whether we are inside a quoted string to avoid corrupting string
// values that legitimately contain "nil" as a substring.
func sanitizeNil(body []byte) []byte {
	// Normalize CRLF to LF so patterns are consistent.
	body = bytes.ReplaceAll(body, []byte("\r\n"), []byte("\n"))

	n := len(body)
	if n < 3 { // "nil" is 3 bytes minimum
		return body
	}

	out := make([]byte, 0, n+8) // small extra room for "null" being 1 byte longer
	inString := false

	for i := 0; i < n; i++ {
		c := body[i]

		// Track quoted strings — skip escaped quotes inside strings.
		if c == '\\' && inString && i+1 < n {
			out = append(out, c, body[i+1])
			i++
			continue
		}
		if c == '"' {
			inString = !inString
			out = append(out, c)
			continue
		}
		if inString {
			out = append(out, c)
			continue
		}

		// Outside a string: check for bare "nil" token.
		if c == 'n' && i+2 < n && body[i+1] == 'i' && body[i+2] == 'l' {
			// Verify it's not part of a longer word (e.g. "null" itself).
			nextOK := i+3 >= n || !isAlpha(body[i+3])
			if nextOK {
				out = append(out, 'n', 'u', 'l', 'l')
				i += 2 // skip "il"
				continue
			}
		}
		out = append(out, c)
	}
	return out
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// Client is the iKuai HTTP API client.
type Client struct {
	BaseURL string
	Token   string
	RawMode bool // When true, check() returns the full envelope; otherwise only data.
	DryRun  bool // When true, write methods (Post/Put/Patch/Delete) print the request and return without executing.
	http    *http.Client
}

// New creates a Client. TLS verification is intentionally skipped (self-signed certs on routers).
func New(baseURL, token string) *Client {
	return NewWithHTTPClient(baseURL, token, &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		},
	})
}

// NewWithHTTPClient creates a Client with a custom http.Client.
func NewWithHTTPClient(baseURL, token string, httpClient *http.Client) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		http:    httpClient,
	}
}

// APIError is returned when the router responds with a non-zero code or 4xx/5xx status.
type APIError struct {
	Code    int
	Message string
	Details []APIErrorDetail
}

// APIErrorDetail holds per-field validation info from the iKuai API.
type APIErrorDetail struct {
	Field string `json:"field"`
	Type  string `json:"type"`
	Msg   string `json:"msg"`
}

func (e *APIError) Error() string {
	s := fmt.Sprintf("[%d] %s", e.Code, e.Message)
	for _, d := range e.Details {
		s += fmt.Sprintf("\n  - %s: %s", d.Field, d.Msg)
	}
	return s
}

// NetworkError wraps connection-level failures (timeout, refused, DNS).
type NetworkError struct {
	Message string
	Cause   error
}

func (e *NetworkError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *NetworkError) Unwrap() error { return e.Cause }

// errorHints maps known iKuai error codes to user-friendly hints.
var errorHints = map[int]string{
	3001: "Parameter error, check your --data or required flags",
	3007: "Token expired or invalid, run: ikuai-cli auth set-url <URL> && ikuai-cli auth set-token <TOKEN>",
	1008: "Session expired, run: ikuai-cli auth set-url <URL> && ikuai-cli auth set-token <TOKEN>",
}

func (c *Client) headers() http.Header {
	h := http.Header{}
	h.Set("Content-Type", "application/json")
	h.Set("Accept", "application/json")
	if c.Token != "" {
		h.Set("Authorization", "Bearer "+c.Token)
	}
	return h
}

func (c *Client) check(resp *http.Response) (json.RawMessage, error) {
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP %d: read error: %w", resp.StatusCode, err)
	}

	// Some iKuai firmware versions emit bare `nil` (not valid JSON).
	// Replace with `null` before parsing.
	body := sanitizeNil(raw)

	// iKuai API uses "data" for most endpoints, "results" for monitor/load,
	// and "rowid" at the top level for create responses.
	var envelope struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
		Results json.RawMessage `json:"results"`
		RowID   json.RawMessage `json:"rowid"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("HTTP %d: non-JSON response", resp.StatusCode)
	}

	if resp.StatusCode >= 400 || (envelope.Code != 0 && envelope.Code != 20000) {
		msg := envelope.Message
		if msg == "" {
			msg = string(body)
		}
		if hint, ok := errorHints[envelope.Code]; ok {
			msg = msg + " (" + hint + ")"
		}
		// Parse details array if present.
		var details []APIErrorDetail
		var detailEnv struct {
			Details []APIErrorDetail `json:"details"`
		}
		if err := json.Unmarshal(body, &detailEnv); err == nil && len(detailEnv.Details) > 0 {
			details = detailEnv.Details
		}
		return nil, &APIError{Code: envelope.Code, Message: msg, Details: details}
	}
	if c.RawMode {
		return body, nil
	}
	// Prefer "data", fall back to "results" (used by monitor/load endpoints).
	payload := envelope.Data
	if (len(payload) == 0 || string(payload) == "null") && len(envelope.Results) > 0 && string(envelope.Results) != "null" {
		payload = envelope.Results
	}

	// When both data and results are absent/null, check for rowid (create responses).
	if len(payload) == 0 || string(payload) == "null" {
		msg := envelope.Message
		if msg == "" {
			msg = "ok"
		}
		if len(envelope.RowID) > 0 && string(envelope.RowID) != "null" {
			// RowID may be integer (42) or string ("42") depending on firmware.
			// Parse as json.Number to handle both, then include raw value.
			var rowid interface{}
			if err := json.Unmarshal(envelope.RowID, &rowid); err == nil {
				synthetic, _ := json.Marshal(map[string]interface{}{
					"message": msg,
					"rowid":   rowid,
				})
				return synthetic, nil
			}
		}
		synthetic, _ := json.Marshal(map[string]interface{}{
			"message": msg,
		})
		return synthetic, nil
	}
	return payload, nil
}

func (c *Client) doJSON(method, path string, body interface{}) (json.RawMessage, error) {
	fullURL := c.BaseURL + path

	var bodyBytes []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyBytes = b
	} else {
		bodyBytes = []byte("{}")
	}

	// Dry-run: return the request preview without executing
	if c.DryRun {
		preview := map[string]interface{}{
			"dry_run": true,
			"method":  method,
			"url":     fullURL,
		}
		if body != nil {
			preview["body"] = json.RawMessage(bodyBytes)
		}
		return json.Marshal(preview)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, fullURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header = c.headers()

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, &NetworkError{Message: "connection failed", Cause: err}
	}
	defer func() { _ = resp.Body.Close() }()
	return c.check(resp)
}

// Get sends a GET request with optional query parameters.
func (c *Client) Get(path string, params map[string]string) (json.RawMessage, error) {
	fullURL := c.BaseURL + path
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			q.Set(k, v)
		}
		fullURL += "?" + q.Encode()
	}

	// Dry-run: return the request preview without executing
	if c.DryRun {
		preview := map[string]interface{}{
			"dry_run": true,
			"method":  "GET",
			"url":     fullURL,
		}
		return json.Marshal(preview)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header = c.headers()

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, &NetworkError{Message: "connection failed", Cause: err}
	}
	defer func() { _ = resp.Body.Close() }()
	return c.check(resp)
}

// Post sends a POST request with a JSON body.
func (c *Client) Post(path string, body interface{}) (json.RawMessage, error) {
	return c.doJSON(http.MethodPost, path, body)
}

// Put sends a PUT request with a JSON body.
func (c *Client) Put(path string, body interface{}) (json.RawMessage, error) {
	return c.doJSON(http.MethodPut, path, body)
}

// Patch sends a PATCH request with a JSON body.
func (c *Client) Patch(path string, body interface{}) (json.RawMessage, error) {
	return c.doJSON(http.MethodPatch, path, body)
}

// Delete sends a DELETE request.
func (c *Client) Delete(path string) (json.RawMessage, error) {
	return c.doJSON(http.MethodDelete, path, nil)
}
