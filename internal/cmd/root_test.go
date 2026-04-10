package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ikuaidev/ikuai-cli/internal/buildinfo"
	"github.com/ikuaidev/ikuai-cli/internal/session"
)

func TestAuthStatusWithFormatJSON(t *testing.T) {
	testRootCommand(t)
	if err := session.SaveBaseURL("https://192.168.1.1"); err != nil {
		t.Fatalf("SaveBaseURL() error = %v", err)
	}
	if err := session.SaveToken("token-123"); err != nil {
		t.Fatalf("SaveToken() error = %v", err)
	}

	var out bytes.Buffer
	oldStdout, oldStderr, oldFormat := stdout, stderr, formatStr
	stdout, stderr, formatStr = &out, &out, "json"
	t.Cleanup(func() {
		stdout, stderr, formatStr = oldStdout, oldStderr, oldFormat
		rootCmd.SetArgs(nil)
	})

	rootCmd.SetArgs([]string{"auth", "status", "--format", "json"})
	if _, err := rootCmd.ExecuteC(); err != nil {
		t.Fatalf("ExecuteC() error = %v", err)
	}

	got := strings.TrimSpace(out.String())
	want := `{"base_url":"https://192.168.1.1","source":"token"}`
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestVersionWithFormatJSON(t *testing.T) {
	testRootCommand(t)

	oldVersion, oldCommit, oldDate := buildinfo.Version, buildinfo.Commit, buildinfo.Date
	buildinfo.Version = "0.1.0"
	buildinfo.Commit = "abc1234"
	buildinfo.Date = "2026-04-07T07:00:00Z"
	t.Cleanup(func() {
		buildinfo.Version, buildinfo.Commit, buildinfo.Date = oldVersion, oldCommit, oldDate
	})

	var out bytes.Buffer
	setRootOutput(t, &out)
	rootCmd.SetArgs([]string{"version", "--format", "json"})
	if _, err := rootCmd.ExecuteC(); err != nil {
		t.Fatalf("ExecuteC() error = %v", err)
	}

	got := strings.TrimSpace(out.String())
	want := `{"commit":"abc1234","date":"2026-04-07T07:00:00Z","name":"ikuai-cli","version":"0.1.0"}`
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestVersionDefaultsToTableOutput(t *testing.T) {
	testRootCommand(t)

	oldVersion, oldCommit, oldDate := buildinfo.Version, buildinfo.Commit, buildinfo.Date
	buildinfo.Version = "0.2.0"
	buildinfo.Commit = "def5678"
	buildinfo.Date = "2026-04-07T09:00:00Z"
	t.Cleanup(func() {
		buildinfo.Version, buildinfo.Commit, buildinfo.Date = oldVersion, oldCommit, oldDate
	})

	var out bytes.Buffer
	setRootOutput(t, &out)

	rootCmd.SetArgs([]string{"version"})
	if _, err := rootCmd.ExecuteC(); err != nil {
		t.Fatalf("ExecuteC() error = %v", err)
	}

	got := out.String()
	// Default format is table — should show key-value pairs, not JSON braces.
	if strings.Contains(got, "{") {
		t.Fatalf("default output should be table, not JSON: %q", got)
	}
	if !strings.Contains(got, "0.2.0") {
		t.Fatalf("table output missing version: %q", got)
	}
}

func TestCompletionBashGeneratesScript(t *testing.T) {
	testRootCommand(t)

	var out bytes.Buffer
	setRootOutput(t, &out)

	rootCmd.SetArgs([]string{"completion", "bash"})
	if _, err := rootCmd.ExecuteC(); err != nil {
		t.Fatalf("ExecuteC() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "ikuai-cli") {
		t.Fatalf("completion output missing command name: %q", got)
	}
	if !strings.Contains(got, "bash completion") {
		t.Fatalf("completion output missing bash marker: %q", got)
	}
}

func TestCompletionZshGeneratesScript(t *testing.T) {
	testRootCommand(t)

	var out bytes.Buffer
	setRootOutput(t, &out)

	rootCmd.SetArgs([]string{"completion", "zsh"})
	if _, err := rootCmd.ExecuteC(); err != nil {
		t.Fatalf("ExecuteC() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "#compdef ikuai-cli") {
		t.Fatalf("zsh completion output missing compdef header: %q", got)
	}
}

func TestCompletionFishGeneratesScript(t *testing.T) {
	testRootCommand(t)

	var out bytes.Buffer
	setRootOutput(t, &out)

	rootCmd.SetArgs([]string{"completion", "fish"})
	if _, err := rootCmd.ExecuteC(); err != nil {
		t.Fatalf("ExecuteC() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "complete -c ikuai-cli") {
		t.Fatalf("fish completion output missing completion command: %q", got)
	}
}

func TestCompletionPowerShellGeneratesScript(t *testing.T) {
	testRootCommand(t)

	var out bytes.Buffer
	setRootOutput(t, &out)

	rootCmd.SetArgs([]string{"completion", "powershell"})
	if _, err := rootCmd.ExecuteC(); err != nil {
		t.Fatalf("ExecuteC() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "Register-ArgumentCompleter") {
		t.Fatalf("powershell completion output missing completer registration: %q", got)
	}
}

func TestRawAndFormatAreMutuallyExclusive(t *testing.T) {
	testRootCommand(t)

	var out bytes.Buffer
	setRootOutput(t, &out)

	rootCmd.SetArgs([]string{"version", "--raw", "--format", "json"})
	_, err := rootCmd.ExecuteC()
	if err == nil {
		t.Fatal("expected error for --raw + --format, got nil")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("error = %q, want mention of mutually exclusive", err.Error())
	}
}

func TestWideAndColumnsAreMutuallyExclusive(t *testing.T) {
	testRootCommand(t)

	var out bytes.Buffer
	setRootOutput(t, &out)

	rootCmd.SetArgs([]string{"version", "--wide", "--columns", "name,version"})
	_, err := rootCmd.ExecuteC()
	if err == nil {
		t.Fatal("expected error for --wide + --columns, got nil")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("error = %q, want mention of mutually exclusive", err.Error())
	}
}

func TestColumnsFlag(t *testing.T) {
	testRootCommand(t)

	var out bytes.Buffer
	setRootOutput(t, &out)

	rootCmd.SetArgs([]string{"version", "--columns", "name,version"})
	if _, err := rootCmd.ExecuteC(); err != nil {
		t.Fatalf("ExecuteC() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "NAME") || !strings.Contains(got, "VERSION") {
		t.Fatalf("--columns should show requested columns: %q", got)
	}
}

func testRootCommand(t *testing.T) {
	t.Helper()
	t.Setenv("IKUAI_CLI_CONFIG_FILE", t.TempDir()+"/config.json")
	t.Setenv("IKUAI_FORCE_TTY", "1") // Prevent TTY auto-detect from switching to JSON in CI
	rootCmd.SetArgs(nil)
}

func setRootOutput(t *testing.T, out *bytes.Buffer) {
	t.Helper()

	oldStdout, oldStderr, oldFormat := stdout, stderr, formatStr
	oldRaw, oldHumanTime, oldDryRun := rawOutput, humanTime, dryRun
	oldWide, oldColumns := wideOutput, columnsStr
	stdout, stderr, formatStr = out, out, "table"
	rawOutput, humanTime, dryRun = false, false, false
	wideOutput, columnsStr = false, ""
	t.Cleanup(func() {
		stdout, stderr, formatStr = oldStdout, oldStderr, oldFormat
		rawOutput, humanTime, dryRun = oldRaw, oldHumanTime, oldDryRun
		wideOutput, columnsStr = oldWide, oldColumns
		rootCmd.SetArgs(nil)
	})
}
