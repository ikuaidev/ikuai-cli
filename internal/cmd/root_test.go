package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/ikuaidev/ikuai-cli/internal/buildinfo"
	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/ikuaidev/ikuai-cli/internal/output"
	"github.com/ikuaidev/ikuai-cli/internal/session"
	"gopkg.in/yaml.v3"
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

func TestVersionDefaultsToHumanText(t *testing.T) {
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
	if strings.Contains(got, "{") {
		t.Fatalf("default output should be human text, not JSON: %q", got)
	}
	if !strings.HasPrefix(got, "ikuai-cli 0.2.0\n") {
		t.Fatalf("human output should start with version: %q", got)
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

func TestCompletionRejectsUnsupportedShell(t *testing.T) {
	testRootCommand(t)

	var out bytes.Buffer
	setRootOutput(t, &out)

	rootCmd.SetArgs([]string{"completion", "elvish"})
	_, err := rootCmd.ExecuteC()
	if err == nil {
		t.Fatal("expected error for unsupported shell, got nil")
	}
	if !strings.Contains(err.Error(), `unknown command "elvish"`) {
		t.Fatalf("error = %q, want unknown command", err.Error())
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

func TestHandleExecuteErrorTTYPrintsSingleHumanError(t *testing.T) {
	var errOut bytes.Buffer
	err := &cliapp.ValidationError{Message: "bad input"}

	code := handleExecuteError(err, true, false, output.Table, &errOut)

	if code != ExitValidation {
		t.Fatalf("exit code = %d, want %d", code, ExitValidation)
	}
	got := errOut.String()
	if strings.Count(got, "Error: bad input") != 1 {
		t.Fatalf("stderr = %q, want exactly one human error", got)
	}
	if strings.Contains(got, `"ok":false`) {
		t.Fatalf("stderr = %q, did not want JSON envelope for TTY", got)
	}
}

func TestHandleExecuteErrorNonTTYPrintsOnlyJSONEnvelope(t *testing.T) {
	var errOut bytes.Buffer
	err := &cliapp.ValidationError{Message: "bad input"}

	code := handleExecuteError(err, false, false, output.Table, &errOut)

	if code != ExitValidation {
		t.Fatalf("exit code = %d, want %d", code, ExitValidation)
	}
	got := errOut.String()
	if strings.Contains(got, "Error:") {
		t.Fatalf("stderr = %q, did not want cobra-style Error prefix", got)
	}
	var envelope struct {
		OK    bool `json:"ok"`
		Error struct {
			Message  string `json:"message"`
			ExitCode int    `json:"exit_code"`
			Type     string `json:"type"`
		} `json:"error"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(errOut.Bytes()), &envelope); err != nil {
		t.Fatalf("stderr JSON error = %v; stderr = %q", err, got)
	}
	if envelope.OK {
		t.Fatalf("ok = true, want false")
	}
	if envelope.Error.Message != "bad input" {
		t.Fatalf("message = %q, want %q", envelope.Error.Message, "bad input")
	}
	if envelope.Error.ExitCode != ExitValidation {
		t.Fatalf("exit_code = %d, want %d", envelope.Error.ExitCode, ExitValidation)
	}
	if envelope.Error.Type != "validation_error" {
		t.Fatalf("type = %q, want validation_error", envelope.Error.Type)
	}
}

func TestHandleExecuteErrorNonTTYExplicitYAMLPrintsYAMLEnvelope(t *testing.T) {
	var errOut bytes.Buffer
	err := &cliapp.ValidationError{Message: "bad input"}

	code := handleExecuteError(err, false, true, output.YAML, &errOut)

	if code != ExitValidation {
		t.Fatalf("exit code = %d, want %d", code, ExitValidation)
	}
	got := errOut.String()
	if strings.Contains(got, "Error:") {
		t.Fatalf("stderr = %q, did not want human error for explicit YAML", got)
	}
	if strings.Contains(got, `{"error"`) {
		t.Fatalf("stderr = %q, did not want JSON for explicit YAML", got)
	}
	var envelope struct {
		OK    bool `yaml:"ok"`
		Error struct {
			Message  string `yaml:"message"`
			ExitCode int    `yaml:"exit_code"`
			Type     string `yaml:"type"`
		} `yaml:"error"`
	}
	if err := yaml.Unmarshal(errOut.Bytes(), &envelope); err != nil {
		t.Fatalf("stderr YAML error = %v; stderr = %q", err, got)
	}
	if envelope.OK {
		t.Fatalf("ok = true, want false")
	}
	if envelope.Error.Message != "bad input" {
		t.Fatalf("message = %q, want %q", envelope.Error.Message, "bad input")
	}
	if envelope.Error.ExitCode != ExitValidation {
		t.Fatalf("exit_code = %d, want %d", envelope.Error.ExitCode, ExitValidation)
	}
	if envelope.Error.Type != "validation_error" {
		t.Fatalf("type = %q, want validation_error", envelope.Error.Type)
	}
}

func TestHandleExecuteErrorExplicitJSONPrintsJSONEnvelopeInTTY(t *testing.T) {
	var errOut bytes.Buffer
	err := &cliapp.ValidationError{Message: "bad input"}

	code := handleExecuteError(err, true, true, output.JSON, &errOut)

	if code != ExitValidation {
		t.Fatalf("exit code = %d, want %d", code, ExitValidation)
	}
	got := errOut.String()
	if strings.Contains(got, "Error:") {
		t.Fatalf("stderr = %q, did not want human error for explicit JSON", got)
	}
	var envelope struct {
		OK    bool `json:"ok"`
		Error struct {
			Type string `json:"type"`
		} `json:"error"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(errOut.Bytes()), &envelope); err != nil {
		t.Fatalf("stderr JSON error = %v; stderr = %q", err, got)
	}
	if envelope.Error.Type != "validation_error" {
		t.Fatalf("type = %q, want validation_error", envelope.Error.Type)
	}
}

func TestRootCommandSilencesCobraErrorPrinting(t *testing.T) {
	testRootCommand(t)

	var out bytes.Buffer
	setRootOutput(t, &out)

	rootCmd.SetArgs([]string{"auth", "set-token"})
	_, err := rootCmd.ExecuteC()
	if err == nil {
		t.Fatal("expected auth set-token to fail")
	}

	got := out.String()
	if strings.Contains(got, "Error:") {
		t.Fatalf("stderr = %q, did not want cobra automatic error output", got)
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
