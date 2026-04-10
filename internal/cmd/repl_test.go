package cmd

import (
	"bytes"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/ikuaidev/ikuai-cli/internal/session"
)

func TestSplitArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "plain args",
			input: "auth status --format json",
			want:  []string{"auth", "status", "--format", "json"},
		},
		{
			name:  "double quoted arg",
			input: `objects ip create --data "{\"name\":\"office wan\"}"`,
			want:  []string{"objects", "ip", "create", "--data", `{"name":"office wan"}`},
		},
		{
			name:  "single quoted arg",
			input: "network dns set --data '{\"upstream\":\"1.1.1.1\"}'",
			want:  []string{"network", "dns", "set", "--data", `{"upstream":"1.1.1.1"}`},
		},
		{
			name:  "mixed spacing",
			input: `  system   web-passwd-reset   --ssh-user   admin  `,
			want:  []string{"system", "web-passwd-reset", "--ssh-user", "admin"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := splitArgs(tt.input)
			if err != nil {
				t.Fatalf("splitArgs(%q) error: %v", tt.input, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("splitArgs(%q) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSplitForCompletion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		pre         string
		wantTokens  []string
		wantPartial string
	}{
		{"", nil, ""},
		{"auth ", []string{"auth"}, ""},
		{"auth lo", []string{"auth"}, "lo"},
		{"network dhcp ", []string{"network", "dhcp"}, ""},
		{"network dhcp l", []string{"network", "dhcp"}, "l"},
		{"  network  dhcp  li", []string{"network", "dhcp"}, "li"},
		{"a", []string{}, "a"},
	}
	for _, tt := range tests {
		tokens, partial := splitForCompletion(tt.pre)
		if !reflect.DeepEqual(tokens, tt.wantTokens) {
			t.Errorf("splitForCompletion(%q) tokens = %v, want %v", tt.pre, tokens, tt.wantTokens)
		}
		if partial != tt.wantPartial {
			t.Errorf("splitForCompletion(%q) partial = %q, want %q", tt.pre, partial, tt.wantPartial)
		}
	}
}

func TestCompleteFromTree(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tokens  []string
		partial string
		minLen  int    // expect at least this many results
		must    string // at least one result must contain this
		empty   bool   // expect no results
	}{
		{"top-level empty", nil, "", 5, "auth", false},
		{"top-level prefix", nil, "au", 1, "auth", false},
		{"auth subcommands", []string{"auth"}, "", 4, "set-url", false},
		{"auth prefix se", []string{"auth"}, "se", 1, "set", false},
		{"network dhcp", []string{"network", "dhcp"}, "", 2, "list", false},
		{"network dhcp l", []string{"network", "dhcp"}, "l", 1, "list", false},
		{"bogus token", []string{"bogus"}, "", 0, "", true},
		{"flag partial", nil, "--", 0, "", true}, // filtered at caller level, but tree also returns nothing
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := completeFromTree(rootCmd, tt.tokens, tt.partial)
			if tt.empty {
				if len(got) != 0 {
					t.Errorf("completeFromTree(%v, %q) = %v, want empty", tt.tokens, tt.partial, got)
				}
				return
			}
			if len(got) < tt.minLen {
				t.Errorf("completeFromTree(%v, %q) returned %d results, want >= %d: %v", tt.tokens, tt.partial, len(got), tt.minLen, got)
			}
			if tt.must != "" {
				sort.Strings(got)
				found := false
				for _, c := range got {
					if c == tt.must || len(c) >= len(tt.must) && c[:len(tt.must)] == tt.must[:len(tt.must)] {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("completeFromTree(%v, %q) = %v, expected to contain %q", tt.tokens, tt.partial, got, tt.must)
				}
			}
		})
	}
}

func saveBannerState(t *testing.T, buf *bytes.Buffer) {
	t.Helper()
	origStdout, origColor := stdout, useColor
	stdout = buf
	t.Cleanup(func() { stdout, useColor = origStdout, origColor })
}

func TestBannerUnauthed(t *testing.T) {
	t.Setenv("IKUAI_FORCE_TTY", "1")
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer
	saveBannerState(t, &buf)

	app.Session = &session.Session{}
	app.CredSource = "none"

	printBanner()
	got := buf.String()

	if strings.Contains(got, "Router:") {
		t.Errorf("unauth banner should NOT contain Router: %q", got)
	}
	if !strings.Contains(got, "Not authenticated") {
		t.Errorf("unauth banner should contain 'Not authenticated': %q", got)
	}
	if !strings.Contains(got, "auth set-url") {
		t.Errorf("unauth banner should contain setup hint: %q", got)
	}
	// NO_COLOR=1 should suppress ANSI escapes.
	if strings.Contains(got, "\x1b[") {
		t.Errorf("banner should not contain ANSI escapes with NO_COLOR=1: %q", got)
	}
}

func TestBannerAuthedSession(t *testing.T) {
	t.Setenv("IKUAI_FORCE_TTY", "1")
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer
	saveBannerState(t, &buf)

	app.Session = &session.Session{BaseURL: "https://10.66.0.20", Token: "token-abc"}
	app.CredSource = "token"

	printBanner()
	got := buf.String()

	if !strings.Contains(got, "Router:") {
		t.Errorf("authed banner should contain 'Router:': %q", got)
	}
	if !strings.Contains(got, "10.66.0.20") {
		t.Errorf("authed banner should contain URL: %q", got)
	}
	if !strings.Contains(got, "authenticated") {
		t.Errorf("authed banner should contain 'authenticated': %q", got)
	}
	if !strings.Contains(got, "via token") {
		t.Errorf("authed banner should show source: %q", got)
	}
}

func TestBannerAuthedEnv(t *testing.T) {
	t.Setenv("IKUAI_FORCE_TTY", "1")
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer
	saveBannerState(t, &buf)

	app.Session = &session.Session{BaseURL: "https://env-router", Token: "env-tok"}
	app.CredSource = "env"

	printBanner()
	got := buf.String()

	if !strings.Contains(got, "via env") {
		t.Errorf("env-mode banner should show 'via env': %q", got)
	}
}

func TestBannerNonTTYSkipped(t *testing.T) {
	t.Setenv("IKUAI_FORCE_TTY", "")

	var buf bytes.Buffer
	saveBannerState(t, &buf)

	app.Session = &session.Session{BaseURL: "https://10.0.0.1", Token: "tok"}
	app.CredSource = "token"

	printBanner()

	if buf.Len() > 0 {
		t.Errorf("banner should be empty when stdout is not TTY, got %d bytes", buf.Len())
	}
}

func TestSplitArgs_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{"trailing backslash", `foo bar\`},
		{"unterminated double quote", `foo "bar baz`},
		{"unterminated single quote", `foo 'bar baz`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := splitArgs(tt.input)
			if err == nil {
				t.Fatalf("splitArgs(%q) expected error, got nil", tt.input)
			}
		})
	}
}
