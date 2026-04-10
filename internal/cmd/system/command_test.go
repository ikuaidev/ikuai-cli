package system

import (
	"io"
	"testing"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
)

func TestNewRegistersExpectedSystemCommands(t *testing.T) {
	t.Parallel()

	cmd := New(cliapp.New(io.Discard, io.Discard))

	tests := []struct {
		name     string
		args     []string
		wantUse  string
		wantFlag string
	}{
		{name: "system set data flag", args: []string{"set"}, wantUse: "set", wantFlag: "data"},
		{name: "system set hostname flag", args: []string{"set"}, wantUse: "set", wantFlag: "hostname"},
		{name: "system set language flag", args: []string{"set"}, wantUse: "set", wantFlag: "language"},
		{name: "system set time-zone flag", args: []string{"set"}, wantUse: "set", wantFlag: "time-zone"},
		{name: "schedules list pagination", args: []string{"schedules", "list"}, wantUse: "list", wantFlag: "page"},
		{name: "schedules create data flag", args: []string{"schedules", "create"}, wantUse: "create", wantFlag: "data"},
		{name: "schedules create name flag", args: []string{"schedules", "create"}, wantUse: "create", wantFlag: "name"},
		{name: "schedules create enabled flag", args: []string{"schedules", "create"}, wantUse: "create", wantFlag: "enabled"},
		{name: "schedules update name flag", args: []string{"schedules", "update"}, wantUse: "update", wantFlag: "name"},
		{name: "remote access set data flag", args: []string{"remote-access", "set"}, wantUse: "set", wantFlag: "data"},
		{name: "remote access set telnet flag", args: []string{"remote-access", "set"}, wantUse: "set", wantFlag: "telnet"},
		{name: "remote access set ssh-port flag", args: []string{"remote-access", "set"}, wantUse: "set", wantFlag: "ssh-port"},
		{name: "vrrp set type flag", args: []string{"vrrp", "set"}, wantUse: "set", wantFlag: "type"},
		{name: "vrrp set enabled flag", args: []string{"vrrp", "set"}, wantUse: "set", wantFlag: "enabled"},
		{name: "alg set ftp flag", args: []string{"alg", "set"}, wantUse: "set", wantFlag: "ftp"},
		{name: "alg set ftp-ports flag", args: []string{"alg", "set"}, wantUse: "set", wantFlag: "ftp-ports"},
		{name: "kernel set bbr flag", args: []string{"kernel", "set"}, wantUse: "set", wantFlag: "bbr"},
		{name: "cpufreq set mode flag", args: []string{"cpufreq", "set"}, wantUse: "set", wantFlag: "mode"},
		{name: "cpufreq mode-set mode flag", args: []string{"cpufreq", "mode-set"}, wantUse: "mode-set", wantFlag: "mode"},
		{name: "web password reset ssh flag", args: []string{"web-passwd", "reset"}, wantUse: "reset", wantFlag: "ssh-user"},
		{name: "web password reset confirm flag", args: []string{"web-passwd", "reset"}, wantUse: "reset", wantFlag: "yes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, _, err := cmd.Find(tt.args)
			if err != nil {
				t.Fatalf("Find(%v) error = %v", tt.args, err)
			}
			if found == nil {
				t.Fatalf("Find(%v) returned nil command", tt.args)
			}
			if found.Name() != tt.wantUse {
				t.Fatalf("Find(%v) command = %q, want %q", tt.args, found.Name(), tt.wantUse)
			}
			if found.Flags().Lookup(tt.wantFlag) == nil {
				t.Fatalf("Find(%v) missing flag %q", tt.args, tt.wantFlag)
			}
		})
	}
}
