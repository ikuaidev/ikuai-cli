package network

import (
	"io"
	"testing"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
)

func TestNewRegistersExpectedNetworkCommands(t *testing.T) {
	t.Parallel()

	cmd := New(cliapp.New(io.Discard, io.Discard))

	tests := []struct {
		name     string
		args     []string
		wantUse  string
		wantFlag string
		noFlag   string
	}{
		{name: "wan list no pagination", args: []string{"wan", "list"}, wantUse: "list", noFlag: "page"},
		{name: "wan-vlan list no pagination", args: []string{"wan-vlan", "list"}, wantUse: "list", noFlag: "page"},
		{name: "lan list no pagination", args: []string{"lan", "list"}, wantUse: "list", noFlag: "page-size"},
		{name: "physical list no pagination", args: []string{"physical", "list"}, wantUse: "list", noFlag: "page"},
		{name: "dns set data flag", args: []string{"dns", "set"}, wantUse: "set", wantFlag: "data"},
		{name: "dns proxy list pagination", args: []string{"dns", "proxy", "list"}, wantUse: "list", wantFlag: "page-size"},
		{name: "dhcp list pagination", args: []string{"dhcp", "list"}, wantUse: "list", wantFlag: "page"},
		{name: "dhcp access mode set data flag", args: []string{"dhcp", "access-mode", "set"}, wantUse: "set", wantFlag: "data"},
		{name: "dhcp6 rule update data flag", args: []string{"dhcp6", "access-rule", "update"}, wantUse: "update", wantFlag: "data"},
		{name: "nat create data flag", args: []string{"nat", "create"}, wantUse: "create", wantFlag: "data"},
		{name: "pppoe set data flag", args: []string{"pppoe", "set"}, wantUse: "set", wantFlag: "data"},
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
			if tt.wantFlag != "" && found.Flags().Lookup(tt.wantFlag) == nil {
				t.Fatalf("Find(%v) missing flag %q", tt.args, tt.wantFlag)
			}
			if tt.noFlag != "" && found.Flags().Lookup(tt.noFlag) != nil {
				t.Fatalf("Find(%v) unexpectedly has flag %q", tt.args, tt.noFlag)
			}
		})
	}
}
