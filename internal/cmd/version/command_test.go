package version

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ikuaidev/ikuai-cli/internal/buildinfo"
	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/ikuaidev/ikuai-cli/internal/output"
)

func TestVersionCommandWithFormatJSON(t *testing.T) {
	oldVersion, oldCommit, oldDate := buildinfo.Version, buildinfo.Commit, buildinfo.Date
	buildinfo.Version = "1.2.3"
	buildinfo.Commit = "abc1234"
	buildinfo.Date = "2026-04-07T07:00:00Z"
	t.Cleanup(func() {
		buildinfo.Version, buildinfo.Commit, buildinfo.Date = oldVersion, oldCommit, oldDate
	})

	var out bytes.Buffer
	app := cliapp.New(&out, &out)
	app.Format = output.JSON

	cmd := New(app)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := strings.TrimSpace(out.String())
	want := `{"commit":"abc1234","date":"2026-04-07T07:00:00Z","name":"ikuai-cli","version":"1.2.3"}`
	if got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestVersionCommandDefaultsToTable(t *testing.T) {
	oldVersion, oldCommit, oldDate := buildinfo.Version, buildinfo.Commit, buildinfo.Date
	buildinfo.Version = "2.0.0"
	buildinfo.Commit = "fedcba9"
	buildinfo.Date = "2026-04-07T10:00:00Z"
	t.Cleanup(func() {
		buildinfo.Version, buildinfo.Commit, buildinfo.Date = oldVersion, oldCommit, oldDate
	})

	var out bytes.Buffer
	app := cliapp.New(&out, &out)

	cmd := New(app)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := out.String()
	// Default is table — should not contain JSON braces.
	if strings.Contains(got, "{") {
		t.Fatalf("default output should be table, not JSON: %q", got)
	}
	if !strings.Contains(got, "2.0.0") {
		t.Fatalf("table output missing version: %q", got)
	}
	if !strings.Contains(got, "fedcba9") {
		t.Fatalf("table output missing commit: %q", got)
	}
}
