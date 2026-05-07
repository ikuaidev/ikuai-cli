package network

import (
	"encoding/json"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

// natFlagDescs provides human-readable descriptions for NAT/DNAT CLI flags.
var natFlagDescs = map[string]string{
	"name":          "Rule name",
	"action":        "NAT action (filter/dnat/snat)",
	"protocol":      "Protocol (tcp/udp/any)",
	"in-interface":  "Inbound interface",
	"out-interface": "Outbound interface",
	"comment":       "Comment",
	"wan-port":      "WAN port",
	"lan-addr":      "LAN target address",
	"lan-port":      "LAN target port",
	"src-addr":      "Source address (comma-separated)",
	"dst-addr":      "Destination address (comma-separated)",
	"src-port":      "Source port (comma-separated)",
	"dst-port":      "Destination port (comma-separated)",
}

func readOnlyNetworkRunE(app *cliapp.Runtime, apiPath string, defaultColumns []string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := app.RequireAuth(); err != nil {
			return err
		}
		app.DefaultColumns = defaultColumns
		raw, err := app.APIClient.Get(cliapp.APIBase+apiPath, nil)
		if err != nil {
			return err
		}
		app.PrintRaw(raw)
		return nil
	}
}

func New(app *cliapp.Runtime) *cobra.Command {
	networkCmd := &cobra.Command{
		Use:   "network",
		Short: "Network config",
		Long:  `Manage network configuration: DNS, DHCP, WAN/LAN interfaces, VLAN, NAT/DNAT, DMZ, PPPoE, and DNS proxy.`,
		Example: `  ikuai-cli network dns get
  ikuai-cli network dhcp list
  ikuai-cli network vlan list
  ikuai-cli network nat list`,
	}

	wanRunE := readOnlyNetworkRunE(app, "/interfaces/wan-config", []string{"id", "tagname", "internet", "ip_mask", "gateway", "mac", "mtu", "speed"})
	networkWANCmd := &cobra.Command{Use: "wan", Short: "WAN interface config", Args: cobra.NoArgs, RunE: wanRunE}
	networkWANListCmd := &cobra.Command{Use: "list", Short: "List WAN interface configs", Args: cobra.NoArgs, RunE: wanRunE}

	wanVLANRunE := readOnlyNetworkRunE(app, "/interfaces/wan-vlan-config", []string{"id", "vlan_name", "interface", "vlan_id", "ip_mask", "gateway", "enabled"})
	networkWANVLANCmd := &cobra.Command{Use: "wan-vlan", Short: "WAN VLAN config", Args: cobra.NoArgs, RunE: wanVLANRunE}
	networkWANVLANListCmd := &cobra.Command{Use: "list", Short: "List WAN VLAN configs", Args: cobra.NoArgs, RunE: wanVLANRunE}

	lanRunE := readOnlyNetworkRunE(app, "/interfaces/lan-config", []string{"id", "tagname", "ip_mask", "bandeth", "dhcp_server", "vlan"})
	networkLANCmd := &cobra.Command{Use: "lan", Short: "LAN interface config", Args: cobra.NoArgs, RunE: lanRunE}
	networkLANListCmd := &cobra.Command{Use: "list", Short: "List LAN interface configs", Args: cobra.NoArgs, RunE: lanRunE}

	physicalRunE := readOnlyNetworkRunE(app, "/interfaces/physical", []string{"name", "interface", "link", "speed", "duplex", "mac", "type", "driver"})
	networkPhysicalCmd := &cobra.Command{Use: "physical", Short: "Physical NIC info", Args: cobra.NoArgs, RunE: physicalRunE}
	networkPhysicalListCmd := &cobra.Command{Use: "list", Short: "List physical NIC info", Args: cobra.NoArgs, RunE: physicalRunE}

	networkDNSCmd := &cobra.Command{Use: "dns", Short: "DNS config"}

	networkDNSGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get DNS config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "dns1", "dns2", "enabled", "cachemode", "cache_ttl", "proxy_force", "network"}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/dns/config", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	dnsSetFieldMap := map[string]string{
		"dns1": "dns1",
		"dns2": "dns2",
	}

	networkDNSSetCmd := &cobra.Command{
		Use:   "set",
		Short: "Set DNS config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := buildDNSSetBody(app, cmd, data, dnsSetFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/network/dns/config", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDNSStatsCmd := &cobra.Command{
		Use:   "stats",
		Short: "DNS query statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/dns/stats", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDNSProxyCmd := &cobra.Command{Use: "proxy", Short: "DNS proxy rules"}

	networkDNSProxyListCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List DNS proxy rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/dns/proxy/rules",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	dnsProxyFieldMap := map[string]string{
		"domain":     "domain",
		"dns-addr":   "dns_addr",
		"src-addr":   "src_addr",
		"parse-type": "parse_type",
		"comment":    "comment",
		"enabled":    "enabled",
	}
	dnsProxyCreateDefaults := map[string]interface{}{
		"enabled":  "yes",
		"comment":  "",
		"src_addr": "",
	}

	networkDNSProxyCreateCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create DNS proxy rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "domain", "dns-addr", "parse-type"); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, dnsProxyFieldMap)
			if err != nil {
				return err
			}
			for k, v := range dnsProxyCreateDefaults {
				if _, exists := body[k]; !exists {
					body[k] = v
				}
			}
			applyDNSProxyDerivedDefaults(body)
			raw, err := app.APIClient.Post(cliapp.APIBase+"/network/dns/proxy/rules", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDNSProxyUpdateCmd := &cobra.Command{
		Use:   "update ID",
		Short: "Update DNS proxy rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "domain", "dns-addr", "parse-type", "enabled"); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, dnsProxyFieldMap)
			if err != nil {
				return err
			}
			for k, v := range dnsProxyCreateDefaults {
				if _, exists := body[k]; !exists {
					body[k] = v
				}
			}
			applyDNSProxyDerivedDefaults(body)
			raw, err := app.APIClient.Put(cliapp.APIBase+"/network/dns/proxy/rules/"+args[0], body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDNSProxyDeleteCmd := newDeleteCmd(app, "Delete DNS proxy rule", "/network/dns/proxy/rules/")

	networkDHCPCmd := &cobra.Command{Use: "dhcp", Short: "DHCP service management"}

	networkDHCPListCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List DHCP services",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "interface", "addr_pool", "gateway", "netmask", "dns1", "lease", "enabled"}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/dhcp/services",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCPGetCmd := &cobra.Command{
		Use:   "get ID",
		Short: "Get DHCP service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "interface", "addr_pool", "gateway", "netmask", "dns1", "lease", "enabled"}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/dhcp/services/"+args[0], nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	dhcpFieldMap := map[string]string{
		"name":         "tagname",
		"interface":    "interface",
		"phy-ifnames":  "phy_ifnames",
		"addr-pool":    "addr_pool",
		"exclude-pool": "exclude_pool",
		"netmask":      "netmask",
		"gateway":      "gateway",
		"dns1":         "dns1",
		"dns2":         "dns2",
		"lease":        "lease",
		"enabled":      "enabled",
	}

	networkDHCPCreateCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create DHCP service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "name", "interface", "addr-pool", "netmask", "gateway", "lease", "phy-ifnames"); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, dhcpFieldMap)
			if err != nil {
				return err
			}
			for k, v := range map[string]interface{}{
				"enabled": "yes", "delay": 0, "check_addr_valid": 0, "check_relay_only": 0,
			} {
				if _, exists := body[k]; !exists {
					body[k] = v
				}
			}
			raw, err := app.APIClient.Post(cliapp.APIBase+"/network/dhcp/services", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCPUpdateCmd := &cobra.Command{
		Use:   "update ID",
		Short: "Update DHCP service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := buildDHCPUpdateBody(app, cmd, data, dhcpFieldMap, args[0])
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/network/dhcp/services/"+args[0], body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	toggleFieldMap := map[string]string{"enabled": "enabled"}

	networkDHCPToggleCmd := &cobra.Command{
		Use:   "toggle ID",
		Short: "Toggle DHCP service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "enabled"); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, toggleFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Patch(cliapp.APIBase+"/network/dhcp/services/"+args[0], body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCPDeleteCmd := newDeleteCmd(app, "Delete DHCP service", "/network/dhcp/services/")

	networkDHCPClientsCmd := &cobra.Command{
		Use:   "clients",
		Short: "DHCP lease clients",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "ip_addr", "mac", "hostname", "interface", "status", "timeout"}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/dhcp/clients",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCPStaticCmd := &cobra.Command{Use: "static", Short: "DHCP static bindings"}

	networkDHCPStaticListCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List DHCP static bindings",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "tagname", "ip_addr", "mac", "interface", "gateway", "hostname", "enabled"}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/dhcp/static",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	dhcpStaticFieldMap := map[string]string{
		"name":      "tagname",
		"ip":        "ip_addr",
		"mac":       "mac",
		"interface": "interface",
		"gateway":   "gateway",
		"dns1":      "dns1",
		"dns2":      "dns2",
		"hostname":  "hostname",
		"comment":   "comment",
		"enabled":   "enabled",
	}

	networkDHCPStaticCreateCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create static binding",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "name", "ip", "mac", "interface"); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, dhcpStaticFieldMap)
			if err != nil {
				return err
			}
			if _, exists := body["enabled"]; !exists {
				body["enabled"] = "yes"
			}
			raw, err := app.APIClient.Post(cliapp.APIBase+"/network/dhcp/static", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCPStaticUpdateCmd := &cobra.Command{
		Use:   "update ID",
		Short: "Update static binding",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := buildDHCPStaticUpdateBody(app, cmd, data, dhcpStaticFieldMap, args[0])
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/network/dhcp/static/"+args[0], body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCPStaticToggleCmd := &cobra.Command{
		Use:   "toggle ID",
		Short: "Toggle static binding",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "enabled"); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, toggleFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Patch(cliapp.APIBase+"/network/dhcp/static/"+args[0], body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCPStaticDeleteCmd := newDeleteCmd(app, "Delete static binding", "/network/dhcp/static/")

	networkDHCPAccessModeCmd := &cobra.Command{Use: "access-mode", Short: "DHCP access mode"}

	networkDHCPAccessModeGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get DHCP access mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/dhcp/access-control/mode", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	dhcpAccessModeFieldMap := map[string]string{"mode": "mode"}
	networkDHCPAccessModeSetCmd := &cobra.Command{
		Use:   "set",
		Short: "Set DHCP access mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "mode"); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, dhcpAccessModeFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/network/dhcp/access-control/mode", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCPAccessRuleCmd := &cobra.Command{Use: "access-rule", Short: "DHCP access rules"}

	networkDHCPAccessRulesCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List DHCP access rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "tagname", "mac", "ip_type", "enabled", "comment"}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/dhcp/access-control/rules",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	dhcpAccessRuleFieldMap := map[string]string{
		"name":    "tagname",
		"mac":     "mac",
		"comment": "comment",
		"enabled": "enabled",
	}
	networkDHCPAccessRuleCreateCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create DHCP access rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "name", "mac", "comment"); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, dhcpAccessRuleFieldMap)
			if err != nil {
				return err
			}
			if _, exists := body["enabled"]; !exists {
				body["enabled"] = "yes"
			}
			raw, err := app.APIClient.Post(cliapp.APIBase+"/network/dhcp/access-control/rules", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCPAccessRuleDeleteCmd := newDeleteCmd(app, "Delete DHCP access rule", "/network/dhcp/access-control/rules/")

	networkDHCPRestartCmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart DHCP",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Post(cliapp.APIBase+"/network/dhcp/services:restart", map[string]string{})
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCPStartCmd := &cobra.Command{
		Use:   "start",
		Short: "Start DHCP",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Post(cliapp.APIBase+"/network/dhcp/services:start", map[string]string{})
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCPStopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop DHCP",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Post(cliapp.APIBase+"/network/dhcp/services:stop", map[string]string{})
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCP6Cmd := &cobra.Command{Use: "dhcp6", Short: "DHCPv6 management"}

	networkDHCP6ClientsCmd := &cobra.Command{
		Use:   "clients",
		Short: "DHCPv6 lease clients",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/dhcp6/clients",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCP6AccessModeCmd := &cobra.Command{Use: "access-mode", Short: "DHCPv6 access mode"}

	networkDHCP6AccessModeGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get DHCPv6 access mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/dhcp6/access-control/mode", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCP6AccessModeSetCmd := &cobra.Command{
		Use:   "set",
		Short: "Set DHCPv6 access mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "mode"); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, dhcpAccessModeFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/network/dhcp6/access-control/mode", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCP6AccessRuleCmd := &cobra.Command{Use: "access-rule", Short: "DHCPv6 access rules"}
	dhcp6AccessRuleFieldMap := map[string]string{
		"name":    "tagname",
		"mac":     "mac",
		"comment": "comment",
		"enabled": "enabled",
	}
	dhcp6AccessRuleInputFields := []string{"enabled", "mac", "tagname", "comment"}

	networkDHCP6AccessRulesCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List DHCPv6 access rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/dhcp6/access-control/rules",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCP6AccessRuleCreateCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create DHCPv6 access rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "name", "mac"); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, dhcp6AccessRuleFieldMap)
			if err != nil {
				return err
			}
			if _, exists := body["enabled"]; !exists {
				body["enabled"] = "yes"
			}
			if _, exists := body["comment"]; !exists {
				body["comment"] = ""
			}
			raw, err := app.APIClient.Post(cliapp.APIBase+"/network/dhcp6/access-control/rules", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCP6AccessRuleUpdateCmd := &cobra.Command{
		Use:   "update ID",
		Short: "Update DHCPv6 access rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			updates, err := cliapp.MergeDataWithFlags(data, cmd, dhcp6AccessRuleFieldMap)
			if err != nil {
				return err
			}
			body, err := buildNATUpdateBody(app, "network/dhcp6/access-control/rules", args[0], updates, map[string]interface{}{"comment": ""}, dhcp6AccessRuleInputFields)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/network/dhcp6/access-control/rules/"+args[0], body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCP6AccessRuleToggleCmd := &cobra.Command{
		Use:   "toggle ID",
		Short: "Toggle DHCPv6 access rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "enabled"); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, toggleFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Patch(cliapp.APIBase+"/network/dhcp6/access-control/rules/"+args[0], body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDHCP6AccessRuleDeleteCmd := newDeleteCmd(app, "Delete DHCPv6 access rule", "/network/dhcp6/access-control/rules/")

	networkVLANCmd := &cobra.Command{Use: "vlan", Short: "VLAN management"}

	networkVLANListCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List VLANs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "vlan_name", "vlan_id", "interface", "ip_addr", "netmask", "enabled"}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/vlan",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	vlanFieldMap := map[string]string{
		"name":      "vlan_name",
		"vlan-id":   "vlan_id",
		"interface": "interface",
		"ip":        "ip_addr",
		"netmask":   "netmask",
		"ip-mask":   "ip_mask",
		"mac":       "mac",
		"comment":   "comment",
		"enabled":   "enabled",
	}

	networkVLANCreateCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create VLAN",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "name", "vlan-id", "interface", "netmask"); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, vlanFieldMap)
			if err != nil {
				return err
			}
			for k, v := range vlanCreateDefaults() {
				if _, exists := body[k]; !exists {
					body[k] = v
				}
			}
			raw, err := app.APIClient.Post(cliapp.APIBase+"/network/vlan", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkVLANUpdateCmd := &cobra.Command{
		Use:   "update ID",
		Short: "Update VLAN",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := buildVLANUpdateBody(app, cmd, data, vlanFieldMap, args[0])
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/network/vlan/"+args[0], body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkVLANToggleCmd := &cobra.Command{
		Use:   "toggle ID",
		Short: "Toggle VLAN",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "enabled"); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, toggleFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Patch(cliapp.APIBase+"/network/vlan/"+args[0], body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkVLANDeleteCmd := newDeleteCmd(app, "Delete VLAN", "/network/vlan/")

	networkPPPoECmd := &cobra.Command{Use: "pppoe", Short: "PPPoE server config"}

	networkPPPoEGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get PPPoE config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "enabled", "server_name", "server_ip", "addr_pool", "interface", "authmode", "mtu", "mru"}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/pppoe/services", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	pppoeFieldMap := map[string]string{
		"enabled":           "enabled",
		"server-name":       "server_name",
		"force-verify-name": "force_verify_name",
		"server-ip":         "server_ip",
		"dns1":              "dns1",
		"dns2":              "dns2",
		"authmode":          "authmode",
		"nas-identifier":    "nas_identifier",
		"nas-ip-address":    "nas_ip_address",
		"radius-ip":         "radius_ip",
		"secret":            "secret",
		"authport":          "authport",
		"accountport":       "accountport",
		"addr-pool":         "addr_pool",
		"interface":         "interface",
		"rate-limit-lan":    "rate_limit_lan",
		"drop-client":       "drop_client",
		"force-pppoe":       "force_pppoe",
		"enhance-check":     "enhance_check",
		"share-deny":        "share_deny",
		"bind-vlan":         "bind_vlan",
		"verify-vlan":       "verify_vlan",
		"bind-iface":        "bind_iface",
		"mtu":               "mtu",
		"mru":               "mru",
		"lcp-echo-interval": "lcp_echo_interval",
		"lcp-echo-failure":  "lcp_echo_failure",
		"maxconnect":        "maxconnect",
		"restart-timer":     "restart_timer",
		"restart-week":      "restart_week",
		"restart-time":      "restart_time",
		"comment":           "comment",
	}
	networkPPPoESetCmd := &cobra.Command{
		Use:   "set",
		Short: "Set PPPoE config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := buildPPPoESetBody(app, cmd, data, pppoeFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/network/pppoe/services", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkWANCmd.AddCommand(networkWANListCmd)
	networkWANVLANCmd.AddCommand(networkWANVLANListCmd)
	networkLANCmd.AddCommand(networkLANListCmd)
	networkPhysicalCmd.AddCommand(networkPhysicalListCmd)
	networkCmd.AddCommand(networkWANCmd, networkWANVLANCmd, networkLANCmd, networkPhysicalCmd)

	networkCmd.AddCommand(networkDNSCmd)
	networkDNSCmd.AddCommand(
		networkDNSGetCmd,
		networkDNSSetCmd,
		networkDNSStatsCmd,
		networkDNSProxyCmd,
	)
	networkDNSProxyCmd.AddCommand(
		networkDNSProxyListCmd,
		networkDNSProxyCreateCmd,
		networkDNSProxyUpdateCmd,
		networkDNSProxyDeleteCmd,
	)
	cliapp.AddListFlags(networkDNSProxyListCmd)
	for _, c := range []*cobra.Command{networkDNSSetCmd, networkDNSProxyCreateCmd, networkDNSProxyUpdateCmd} {
		c.Flags().String("data", "{}", "JSON body")
	}
	networkDNSSetCmd.Flags().String("dns1", "", "Primary DNS server")
	networkDNSSetCmd.Flags().String("dns2", "", "Secondary DNS server")
	// DNS proxy create/update semantic flags
	for _, c := range []*cobra.Command{networkDNSProxyCreateCmd, networkDNSProxyUpdateCmd} {
		c.Flags().String("domain", "", "Domain name")
		c.Flags().String("dns-addr", "", "DNS resolve IP address")
		c.Flags().String("src-addr", "", "Source IP range (comma-separated)")
		c.Flags().String("parse-type", "", "Parse type: ipv4, ipv6, proxy, proxy6")
		c.Flags().String("comment", "", "Comment")
		cliapp.AddEnabledFlag(c)
	}
	cliapp.MarkFlagsRequired(networkDNSProxyCreateCmd, "domain", "dns-addr", "parse-type")
	cliapp.MarkFlagsRequired(networkDNSProxyUpdateCmd, "domain", "dns-addr", "parse-type", "enabled")

	networkCmd.AddCommand(networkDHCPCmd)
	networkDHCPCmd.AddCommand(
		networkDHCPListCmd,
		networkDHCPGetCmd,
		networkDHCPCreateCmd,
		networkDHCPUpdateCmd,
		networkDHCPToggleCmd,
		networkDHCPDeleteCmd,
		networkDHCPClientsCmd,
		networkDHCPStaticCmd,
		networkDHCPAccessModeCmd,
		networkDHCPAccessRuleCmd,
		networkDHCPRestartCmd,
		networkDHCPStartCmd,
		networkDHCPStopCmd,
	)
	networkDHCPStaticCmd.AddCommand(
		networkDHCPStaticListCmd,
		networkDHCPStaticCreateCmd,
		networkDHCPStaticUpdateCmd,
		networkDHCPStaticToggleCmd,
		networkDHCPStaticDeleteCmd,
	)
	networkDHCPAccessModeCmd.AddCommand(
		networkDHCPAccessModeGetCmd,
		networkDHCPAccessModeSetCmd,
	)
	networkDHCPAccessRuleCmd.AddCommand(
		networkDHCPAccessRulesCmd,
		networkDHCPAccessRuleCreateCmd,
		networkDHCPAccessRuleDeleteCmd,
	)
	cliapp.AddListFlags(networkDHCPListCmd)
	cliapp.AddListFlags(networkDHCPClientsCmd)
	cliapp.AddListFlags(networkDHCPStaticListCmd)
	cliapp.AddListFlags(networkDHCPAccessRulesCmd)
	for _, c := range []*cobra.Command{
		networkDHCPCreateCmd, networkDHCPUpdateCmd, networkDHCPToggleCmd,
		networkDHCPStaticCreateCmd, networkDHCPStaticUpdateCmd, networkDHCPStaticToggleCmd,
		networkDHCPAccessModeSetCmd, networkDHCPAccessRuleCreateCmd,
	} {
		c.Flags().String("data", "{}", "JSON body")
	}
	cliapp.AddEnabledFlag(networkDHCPToggleCmd)
	cliapp.MarkFlagsRequired(networkDHCPToggleCmd, "enabled")
	cliapp.AddEnabledFlag(networkDHCPStaticToggleCmd)
	cliapp.MarkFlagsRequired(networkDHCPStaticToggleCmd, "enabled")
	// DHCP access-mode set semantic flags
	networkDHCPAccessModeSetCmd.Flags().String("mode", "", "Access mode (0=blacklist, 1=whitelist, 2=sync MAC ACL)")
	cliapp.MarkFlagsRequired(networkDHCPAccessModeSetCmd, "mode")
	// DHCP access-rule create semantic flags
	networkDHCPAccessRuleCreateCmd.Flags().String("name", "", "Rule name")
	networkDHCPAccessRuleCreateCmd.Flags().String("mac", "", "MAC address")
	networkDHCPAccessRuleCreateCmd.Flags().String("comment", "", "Comment")
	cliapp.AddEnabledFlag(networkDHCPAccessRuleCreateCmd)
	// DHCP create/update semantic flags
	for _, c := range []*cobra.Command{networkDHCPCreateCmd, networkDHCPUpdateCmd} {
		c.Flags().String("name", "", "Service name")
		c.Flags().String("interface", "", "Interface name")
		c.Flags().String("phy-ifnames", "", "Physical interface names")
		c.Flags().String("addr-pool", "", "Address pool range")
		c.Flags().String("exclude-pool", "", "Excluded address pool")
		c.Flags().String("netmask", "", "Subnet mask")
		c.Flags().String("gateway", "", "Gateway address")
		c.Flags().String("dns1", "", "Primary DNS server")
		c.Flags().String("dns2", "", "Secondary DNS server")
		c.Flags().String("lease", "", "Lease time")
		cliapp.AddEnabledFlag(c)
	}
	cliapp.MarkFlagsRequired(networkDHCPCreateCmd, "name", "interface", "addr-pool", "netmask", "gateway", "lease", "phy-ifnames")
	// DHCP static create/update semantic flags
	for _, c := range []*cobra.Command{networkDHCPStaticCreateCmd, networkDHCPStaticUpdateCmd} {
		c.Flags().String("name", "", "Binding name")
		c.Flags().String("ip", "", "IP address")
		c.Flags().String("mac", "", "MAC address")
		c.Flags().String("interface", "", "Interface name")
		c.Flags().String("gateway", "", "Gateway address")
		c.Flags().String("dns1", "", "Primary DNS server")
		c.Flags().String("dns2", "", "Secondary DNS server")
		c.Flags().String("hostname", "", "Hostname")
		c.Flags().String("comment", "", "Comment")
		cliapp.AddEnabledFlag(c)
	}
	cliapp.MarkFlagsRequired(networkDHCPStaticCreateCmd, "name", "ip", "mac", "interface")
	cliapp.MarkFlagsRequired(networkDHCPAccessRuleCreateCmd, "name", "mac", "comment")

	networkCmd.AddCommand(networkDHCP6Cmd)
	networkDHCP6Cmd.AddCommand(
		networkDHCP6ClientsCmd,
		networkDHCP6AccessModeCmd,
		networkDHCP6AccessRuleCmd,
	)
	networkDHCP6AccessModeCmd.AddCommand(
		networkDHCP6AccessModeGetCmd,
		networkDHCP6AccessModeSetCmd,
	)
	networkDHCP6AccessRuleCmd.AddCommand(
		networkDHCP6AccessRulesCmd,
		networkDHCP6AccessRuleCreateCmd,
		networkDHCP6AccessRuleUpdateCmd,
		networkDHCP6AccessRuleToggleCmd,
		networkDHCP6AccessRuleDeleteCmd,
	)
	cliapp.AddListFlags(networkDHCP6ClientsCmd)
	cliapp.AddListFlags(networkDHCP6AccessRulesCmd)
	for _, c := range []*cobra.Command{
		networkDHCP6AccessModeSetCmd, networkDHCP6AccessRuleCreateCmd,
		networkDHCP6AccessRuleUpdateCmd, networkDHCP6AccessRuleToggleCmd,
	} {
		c.Flags().String("data", "{}", "JSON body")
	}
	networkDHCP6AccessModeSetCmd.Flags().String("mode", "", "Access mode (0=blacklist, 1=whitelist)")
	cliapp.MarkFlagsRequired(networkDHCP6AccessModeSetCmd, "mode")
	for _, c := range []*cobra.Command{networkDHCP6AccessRuleCreateCmd, networkDHCP6AccessRuleUpdateCmd} {
		c.Flags().String("name", "", "Rule name")
		c.Flags().String("mac", "", "MAC address")
		c.Flags().String("comment", "", "Comment")
		cliapp.AddEnabledFlag(c)
	}
	cliapp.MarkFlagsRequired(networkDHCP6AccessRuleCreateCmd, "name", "mac")
	cliapp.AddEnabledFlag(networkDHCP6AccessRuleToggleCmd)
	cliapp.MarkFlagsRequired(networkDHCP6AccessRuleToggleCmd, "enabled")

	networkCmd.AddCommand(networkVLANCmd)
	networkVLANCmd.AddCommand(
		networkVLANListCmd,
		networkVLANCreateCmd,
		networkVLANUpdateCmd,
		networkVLANToggleCmd,
		networkVLANDeleteCmd,
	)
	cliapp.AddListFlags(networkVLANListCmd)
	for _, c := range []*cobra.Command{networkVLANCreateCmd, networkVLANUpdateCmd, networkVLANToggleCmd} {
		c.Flags().String("data", "{}", "JSON body")
	}
	cliapp.AddEnabledFlag(networkVLANToggleCmd)
	cliapp.MarkFlagsRequired(networkVLANToggleCmd, "enabled")
	// VLAN create/update semantic flags
	for _, c := range []*cobra.Command{networkVLANCreateCmd, networkVLANUpdateCmd} {
		c.Flags().String("name", "", "VLAN name")
		c.Flags().String("vlan-id", "", "VLAN ID")
		c.Flags().String("interface", "", "Interface name")
		c.Flags().String("ip", "", "IP address")
		c.Flags().String("netmask", "", "Subnet mask")
		c.Flags().String("ip-mask", "", "Additional IP/mask entries")
		c.Flags().String("mac", "", "MAC address")
		c.Flags().String("comment", "", "Comment")
		cliapp.AddEnabledFlag(c)
	}
	cliapp.MarkFlagsRequired(networkVLANCreateCmd, "name", "vlan-id", "interface", "netmask")

	natFieldMap := map[string]string{
		"name":          "tagname",
		"action":        "action",
		"in-interface":  "iinterface",
		"out-interface": "ointerface",
		"protocol":      "protocol",
		"nat-addr":      "nat_addr",
		"nat-port":      "nat_port",
		"comment":       "comment",
		"enabled":       "enabled",
	}
	natCreateDefaults := map[string]interface{}{
		"enabled":      "yes",
		"src_addr":     map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"dst_addr":     map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"src_port":     map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"dst_port":     map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"nat_addr":     "",
		"nat_port":     "",
		"protocol":     "",
		"comment":      "",
		"src_addr_inv": 0,
		"dst_addr_inv": 0,
	}
	natInputFields := []string{
		"tagname", "enabled", "action", "iinterface", "ointerface",
		"src_addr", "dst_addr", "nat_addr", "nat_port", "protocol", "comment",
		"src_port", "dst_port", "src_addr_inv", "dst_addr_inv",
	}
	natAddrFields := map[string]string{
		"src-addr": "src_addr",
		"dst-addr": "dst_addr",
		"src-port": "src_port",
		"dst-port": "dst_port",
	}
	dnatFieldMap := map[string]string{
		"name":      "tagname",
		"lan-addr":  "lan_addr",
		"lan-port":  "lan_port",
		"wan-port":  "wan_port",
		"protocol":  "protocol",
		"interface": "interface",
		"comment":   "comment",
		"enabled":   "enabled",
	}
	dnatCreateDefaults := map[string]interface{}{
		"enabled":  "yes",
		"src_addr": map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"comment":  "",
	}
	dnatInputFields := []string{
		"id", "tagname", "enabled", "lan_addr", "lan_port", "protocol", "interface", "wan_port", "src_addr", "comment",
	}
	dnatAddrFields := map[string]string{
		"src-addr": "src_addr",
	}
	dmzFieldMap := map[string]string{
		"name":      "tagname",
		"interface": "interface",
		"lan-addr":  "lan_addr",
		"protocol":  "protocol",
		"excl-port": "excl_port",
		"comment":   "comment",
		"enabled":   "enabled",
	}
	dmzCreateDefaults := map[string]interface{}{
		"enabled": "yes",
		"comment": "",
	}
	dmzInputFields := []string{
		"id", "tagname", "enabled", "interface", "lan_addr", "protocol", "excl_port", "comment",
	}

	networkCmd.AddCommand(
		natGroup(app, "nat", "NAT rules", "network/nat/rules", natFieldMap, natGroupOpts{
			createDefaults:      natCreateDefaults,
			addrFields:          natAddrFields,
			defaultColumns:      []string{"id", "tagname", "action", "src_addr", "dst_addr", "protocol", "iinterface", "ointerface", "enabled"},
			requiredCreateFlags: []string{"name", "action", "in-interface", "out-interface"},
			inputFields:         natInputFields,
		}),
		natGroup(app, "dnat", "DNAT (port-forward) rules", "network/dnat/rules", dnatFieldMap, natGroupOpts{
			createDefaults:      dnatCreateDefaults,
			addrFields:          dnatAddrFields,
			defaultColumns:      []string{"id", "tagname", "lan_addr", "lan_port", "wan_port", "protocol", "interface", "enabled"},
			requiredCreateFlags: []string{"name", "lan-addr", "lan-port", "wan-port", "protocol", "interface"},
			inputFields:         dnatInputFields,
		}),
		natGroup(app, "dmz", "DMZ rules", "network/dmz/rules", dmzFieldMap, natGroupOpts{
			createDefaults:      dmzCreateDefaults,
			defaultColumns:      []string{"id", "tagname", "interface", "lan_addr", "protocol", "excl_port", "enabled"},
			requiredCreateFlags: []string{"name", "interface", "lan-addr", "protocol", "excl-port"},
			inputFields:         dmzInputFields,
		}),
	)

	networkCmd.AddCommand(networkPPPoECmd)
	networkPPPoECmd.AddCommand(networkPPPoEGetCmd, networkPPPoESetCmd)
	networkPPPoESetCmd.Flags().String("data", "{}", "JSON body")
	networkPPPoESetCmd.Flags().String("enabled", "", "Service enabled (yes/no)")
	networkPPPoESetCmd.Flags().String("server-name", "", "Server name")
	networkPPPoESetCmd.Flags().String("force-verify-name", "", "Force server name verification (0/1)")
	networkPPPoESetCmd.Flags().String("server-ip", "", "Server IP address")
	networkPPPoESetCmd.Flags().String("dns1", "", "Primary DNS")
	networkPPPoESetCmd.Flags().String("dns2", "", "Secondary DNS")
	networkPPPoESetCmd.Flags().String("addr-pool", "", "Address pool range")
	networkPPPoESetCmd.Flags().String("interface", "", "Interface name")
	networkPPPoESetCmd.Flags().String("authmode", "", "Auth mode (0=local, 1=radius)")
	networkPPPoESetCmd.Flags().String("nas-identifier", "", "NAS identifier")
	networkPPPoESetCmd.Flags().String("nas-ip-address", "", "NAS IP address")
	networkPPPoESetCmd.Flags().String("authport", "", "RADIUS auth port")
	networkPPPoESetCmd.Flags().String("accountport", "", "RADIUS accounting port")
	networkPPPoESetCmd.Flags().String("rate-limit-lan", "", "Limit LAN access speed (0/1)")
	networkPPPoESetCmd.Flags().String("drop-client", "", "Block client-to-client access (0/1)")
	networkPPPoESetCmd.Flags().String("force-pppoe", "", "Force PPPoE access (0/1)")
	networkPPPoESetCmd.Flags().String("enhance-check", "", "Enhanced disconnect detection (0/1)")
	networkPPPoESetCmd.Flags().String("share-deny", "", "Shared connection overflow action (0=kick, 1=deny)")
	networkPPPoESetCmd.Flags().String("bind-vlan", "", "Enable VLAN passthrough binding (0/1)")
	networkPPPoESetCmd.Flags().String("verify-vlan", "", "Verify VLAN when bind-vlan is enabled (0/1)")
	networkPPPoESetCmd.Flags().String("bind-iface", "", "Enable interface binding (0/1)")
	networkPPPoESetCmd.Flags().String("mtu", "", "MTU value")
	networkPPPoESetCmd.Flags().String("mru", "", "MRU value")
	networkPPPoESetCmd.Flags().String("lcp-echo-interval", "", "LCP echo interval in seconds")
	networkPPPoESetCmd.Flags().String("lcp-echo-failure", "", "LCP echo failure count")
	networkPPPoESetCmd.Flags().String("maxconnect", "", "Maximum client connection duration in hours")
	networkPPPoESetCmd.Flags().String("restart-timer", "", "Enable scheduled PPPoE restart (0/1)")
	networkPPPoESetCmd.Flags().String("restart-week", "", "Scheduled restart weekdays, e.g. 1234567")
	networkPPPoESetCmd.Flags().String("restart-time", "", "Scheduled restart time, e.g. 06:00")
	networkPPPoESetCmd.Flags().String("radius-ip", "", "RADIUS server IP")
	networkPPPoESetCmd.Flags().String("secret", "", "RADIUS shared secret")
	networkPPPoESetCmd.Flags().String("comment", "", "Comment")

	return networkCmd
}

func buildDNSSetBody(app *cliapp.Runtime, cmd *cobra.Command, data string, fieldMap map[string]string) (map[string]interface{}, error) {
	body, err := cliapp.MergeDataWithFlags(data, cmd, fieldMap)
	if err != nil {
		return nil, err
	}

	readClient := app.APIClient
	if app.APIClient.DryRun {
		readClient = app.NewClient(app.Session.BaseURL, app.Session.Token)
	}
	raw, err := readClient.Get(cliapp.APIBase+"/network/dns/config", nil)
	if err != nil {
		if hasAllDNSConfigInputFields(body) {
			return body, nil
		}
		return nil, err
	}

	current, err := extractDNSConfig(raw)
	if err != nil {
		if hasAllDNSConfigInputFields(body) {
			return body, nil
		}
		return nil, err
	}
	for _, key := range dnsConfigInputFields {
		if _, exists := body[key]; !exists {
			if v, ok := current[key]; ok {
				body[key] = v
			}
		}
	}
	return body, nil
}

var dnsConfigInputFields = []string{
	"enabled",
	"forbid_dns_4a",
	"cache_ttl",
	"cachemode",
	"proxy_force",
	"proxy_force_dns",
	"dns1",
	"dns2",
	"query",
	"defense",
	"network",
	"query_args_ip",
	"query_head_ip",
}

func hasAllDNSConfigInputFields(body map[string]interface{}) bool {
	for _, key := range dnsConfigInputFields {
		if _, ok := body[key]; !ok {
			return false
		}
	}
	return true
}

func extractDNSConfig(raw json.RawMessage) (map[string]interface{}, error) {
	var payload interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	return findFirstObject(payload)
}

func findFirstObject(v interface{}) (map[string]interface{}, error) {
	switch data := v.(type) {
	case map[string]interface{}:
		if arr, ok := data["data"].([]interface{}); ok && len(arr) > 0 {
			return findFirstObject(arr[0])
		}
		return data, nil
	case []interface{}:
		if len(data) == 0 {
			return nil, &cliapp.ValidationError{Message: "empty DNS config response"}
		}
		return findFirstObject(data[0])
	default:
		return nil, &cliapp.ValidationError{Message: "unexpected DNS config response"}
	}
}

func applyDNSProxyDerivedDefaults(body map[string]interface{}) {
	if _, exists := body["is_ipv6"]; exists {
		return
	}
	parseType, _ := body["parse_type"].(string)
	switch parseType {
	case "ipv6", "proxy6":
		body["is_ipv6"] = 1
	default:
		body["is_ipv6"] = 0
	}
}

func buildDHCPStaticUpdateBody(app *cliapp.Runtime, cmd *cobra.Command, data string, fieldMap map[string]string, id string) (map[string]interface{}, error) {
	updates, err := cliapp.MergeDataWithFlags(data, cmd, fieldMap)
	if err != nil {
		return nil, err
	}

	readClient := app.APIClient
	if app.APIClient.DryRun {
		readClient = app.NewClient(app.Session.BaseURL, app.Session.Token)
	}
	raw, err := readClient.Get(cliapp.APIBase+"/network/dhcp/static/"+id, nil)
	if err != nil {
		if hasRequiredDHCPStaticInputFields(updates) {
			return updates, nil
		}
		return nil, err
	}
	current, err := extractDHCPStaticInput(raw)
	if err != nil {
		if hasRequiredDHCPStaticInputFields(updates) {
			return updates, nil
		}
		return nil, err
	}
	for k, v := range updates {
		current[k] = v
	}
	if _, exists := current["enabled"]; !exists {
		current["enabled"] = "yes"
	}
	return current, nil
}

func extractDHCPStaticInput(raw json.RawMessage) (map[string]interface{}, error) {
	var payload interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	current, err := findFirstDHCPStaticObject(payload)
	if err != nil {
		return nil, err
	}
	body := make(map[string]interface{})
	for _, key := range dhcpStaticInputFields {
		if v, ok := current[key]; ok {
			body[key] = v
		}
	}
	return body, nil
}

func findFirstDHCPStaticObject(v interface{}) (map[string]interface{}, error) {
	switch data := v.(type) {
	case map[string]interface{}:
		if results, ok := data["results"]; ok {
			return findFirstDHCPStaticObject(results)
		}
		if arr, ok := data["static_data"].([]interface{}); ok {
			if len(arr) == 0 {
				return nil, &cliapp.ValidationError{Message: "empty DHCP static response"}
			}
			return findFirstDHCPStaticObject(arr[0])
		}
		if arr, ok := data["data"].([]interface{}); ok {
			if len(arr) == 0 {
				return nil, &cliapp.ValidationError{Message: "empty DHCP static response"}
			}
			return findFirstDHCPStaticObject(arr[0])
		}
		return data, nil
	case []interface{}:
		if len(data) == 0 {
			return nil, &cliapp.ValidationError{Message: "empty DHCP static response"}
		}
		return findFirstDHCPStaticObject(data[0])
	default:
		return nil, &cliapp.ValidationError{Message: "unexpected DHCP static response"}
	}
}

func hasRequiredDHCPStaticInputFields(body map[string]interface{}) bool {
	for _, key := range dhcpStaticRequiredInputFields {
		if _, ok := body[key]; !ok {
			return false
		}
	}
	return true
}

var dhcpStaticRequiredInputFields = []string{
	"enabled",
	"mac",
	"ip_addr",
	"interface",
	"tagname",
}

var dhcpStaticInputFields = []string{
	"enabled",
	"mac",
	"ip_addr",
	"interface",
	"gateway",
	"dns1",
	"dns2",
	"comment",
	"tagname",
	"hostname",
}

func buildDHCPUpdateBody(app *cliapp.Runtime, cmd *cobra.Command, data string, fieldMap map[string]string, id string) (map[string]interface{}, error) {
	updates, err := cliapp.MergeDataWithFlags(data, cmd, fieldMap)
	if err != nil {
		return nil, err
	}

	readClient := app.APIClient
	if app.APIClient.DryRun {
		readClient = app.NewClient(app.Session.BaseURL, app.Session.Token)
	}
	raw, err := readClient.Get(cliapp.APIBase+"/network/dhcp/services/"+id, nil)
	if err != nil {
		if hasRequiredDHCPServiceInputFields(updates) {
			return updates, nil
		}
		return nil, err
	}
	current, err := extractDHCPServiceInput(raw)
	if err != nil {
		if hasRequiredDHCPServiceInputFields(updates) {
			return updates, nil
		}
		return nil, err
	}
	for k, v := range updates {
		current[k] = v
	}
	for k, v := range dhcpCreateDefaults() {
		if _, exists := current[k]; !exists {
			current[k] = v
		}
	}
	for k, v := range dhcpUpdateDefaults() {
		if _, exists := current[k]; !exists {
			current[k] = v
		}
	}
	return current, nil
}

func extractDHCPServiceInput(raw json.RawMessage) (map[string]interface{}, error) {
	current, err := extractDNSConfig(raw)
	if err != nil {
		return nil, err
	}
	body := make(map[string]interface{})
	for _, key := range dhcpServiceInputFields {
		v, ok := current[key]
		if !ok || shouldOmitEmptyDHCPInput(key, v) {
			continue
		}
		body[key] = v
	}
	return body, nil
}

func shouldOmitEmptyDHCPInput(key string, v interface{}) bool {
	s, ok := v.(string)
	if !ok || s != "" {
		return false
	}
	switch key {
	case "opt_type15", "opt_type28", "opt_type43", "opt_type60", "opt_type66", "opt_type67", "opt_type80", "opt_type119", "opt_type125", "opt_type128", "opt_type138", "opt_type121":
		return true
	default:
		return false
	}
}

func hasRequiredDHCPServiceInputFields(body map[string]interface{}) bool {
	for _, key := range dhcpServiceRequiredInputFields {
		if _, ok := body[key]; !ok {
			return false
		}
	}
	return true
}

func dhcpCreateDefaults() map[string]interface{} {
	return map[string]interface{}{
		"enabled":          "yes",
		"delay":            0,
		"check_addr_valid": 0,
		"check_relay_only": 0,
	}
}

func dhcpUpdateDefaults() map[string]interface{} {
	return map[string]interface{}{
		"opt_type15":  0,
		"opt_type28":  0,
		"opt_type43":  0,
		"opt_type60":  0,
		"opt_type66":  0,
		"opt_type67":  0,
		"opt_type80":  0,
		"opt_type119": 0,
		"opt_type125": 0,
		"opt_type128": 0,
		"opt_type138": 0,
		"opt_type121": 2,
	}
}

var dhcpServiceRequiredInputFields = []string{
	"enabled",
	"tagname",
	"interface",
	"phy_ifnames",
	"addr_pool",
	"netmask",
	"gateway",
	"lease",
	"delay",
	"check_addr_valid",
	"check_relay_only",
}

var dhcpServiceInputFields = []string{
	"enabled",
	"tagname",
	"interface",
	"phy_ifnames",
	"addr_pool",
	"exclude_pool",
	"netmask",
	"gateway",
	"dns1",
	"dns2",
	"wins1",
	"wins2",
	"domain",
	"next_server",
	"lease",
	"delay",
	"opt_type15",
	"opt15",
	"opt_type28",
	"opt28",
	"opt_type43",
	"opt43",
	"opt_type60",
	"opt60",
	"opt_type66",
	"opt66",
	"opt_type67",
	"opt67",
	"opt_type80",
	"opt80",
	"opt_type119",
	"opt119",
	"opt_type125",
	"opt125",
	"opt_type128",
	"opt128",
	"opt_type138",
	"opt138",
	"opt_type121",
	"opt121",
	"check_addr_valid",
	"check_relay_only",
}

func buildVLANUpdateBody(app *cliapp.Runtime, cmd *cobra.Command, data string, fieldMap map[string]string, id string) (map[string]interface{}, error) {
	updates, err := cliapp.MergeDataWithFlags(data, cmd, fieldMap)
	if err != nil {
		return nil, err
	}

	readClient := app.APIClient
	if app.APIClient.DryRun {
		readClient = app.NewClient(app.Session.BaseURL, app.Session.Token)
	}
	raw, err := readClient.Get(cliapp.APIBase+"/network/vlan/"+id, nil)
	if err != nil {
		if hasRequiredVLANInputFields(updates) {
			return updates, nil
		}
		return nil, err
	}
	current, err := extractVLANInput(raw)
	if err != nil {
		if hasRequiredVLANInputFields(updates) {
			return updates, nil
		}
		return nil, err
	}
	for k, v := range updates {
		current[k] = v
	}
	for k, v := range vlanCreateDefaults() {
		if _, exists := current[k]; !exists {
			current[k] = v
		}
	}
	if vlanName, ok := current["vlan_name"]; ok {
		current["tagname"] = vlanName
	}
	return current, nil
}

func buildPPPoESetBody(app *cliapp.Runtime, cmd *cobra.Command, data string, fieldMap map[string]string) (map[string]interface{}, error) {
	updates, err := cliapp.MergeDataWithFlags(data, cmd, fieldMap)
	if err != nil {
		return nil, err
	}

	readClient := app.APIClient
	if app.APIClient.DryRun {
		readClient = app.NewClient(app.Session.BaseURL, app.Session.Token)
	}
	raw, err := readClient.Get(cliapp.APIBase+"/network/pppoe/services", nil)
	if err != nil {
		if hasInputFields(updates, pppoeRequiredInputFields) {
			return updates, nil
		}
		return nil, err
	}
	current, err := extractInputObject(raw, pppoeInputFields)
	if err != nil {
		if hasInputFields(updates, pppoeRequiredInputFields) {
			return updates, nil
		}
		return nil, err
	}
	for k, v := range updates {
		current[k] = v
	}
	for k, v := range pppoeSetDefaults() {
		if _, exists := current[k]; !exists {
			current[k] = v
		}
	}
	return current, nil
}

func pppoeSetDefaults() map[string]interface{} {
	return map[string]interface{}{
		"server_name":       "iKuai",
		"force_verify_name": 0,
		"nas_identifier":    "iKuai",
		"nas_ip_address":    "",
		"authport":          1812,
		"accountport":       1813,
		"rate_limit_lan":    1,
		"drop_client":       1,
		"force_pppoe":       0,
		"enhance_check":     1,
		"share_deny":        0,
		"bind_vlan":         0,
		"verify_vlan":       1,
		"bind_iface":        0,
		"mtu":               1480,
		"mru":               1480,
		"lcp_echo_interval": 10,
		"lcp_echo_failure":  3,
		"maxconnect":        0,
		"restart_timer":     0,
		"comment":           "",
	}
}

var pppoeRequiredInputFields = []string{
	"enabled",
	"force_verify_name",
	"server_ip",
	"dns1",
	"dns2",
	"authmode",
	"nas_identifier",
	"nas_ip_address",
	"radius_ip",
	"secret",
	"authport",
	"accountport",
	"addr_pool",
	"interface",
	"rate_limit_lan",
	"drop_client",
	"force_pppoe",
	"enhance_check",
	"share_deny",
	"bind_vlan",
	"verify_vlan",
	"bind_iface",
	"mtu",
	"mru",
	"lcp_echo_interval",
	"lcp_echo_failure",
	"maxconnect",
	"restart_timer",
	"comment",
}

var pppoeInputFields = []string{
	"enabled",
	"server_name",
	"force_verify_name",
	"server_ip",
	"dns1",
	"dns2",
	"authmode",
	"nas_identifier",
	"nas_ip_address",
	"radius_ip",
	"secret",
	"authport",
	"accountport",
	"addr_pool",
	"interface",
	"rate_limit_lan",
	"drop_client",
	"force_pppoe",
	"enhance_check",
	"share_deny",
	"bind_vlan",
	"verify_vlan",
	"bind_iface",
	"mtu",
	"mru",
	"lcp_echo_interval",
	"lcp_echo_failure",
	"maxconnect",
	"restart_timer",
	"restart_week",
	"restart_time",
	"comment",
}

func extractVLANInput(raw json.RawMessage) (map[string]interface{}, error) {
	var payload interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	current, err := findFirstDHCPStaticObject(payload)
	if err != nil {
		return nil, err
	}
	body := make(map[string]interface{})
	for _, key := range vlanInputFields {
		if v, ok := current[key]; ok {
			body[key] = v
		}
	}
	return body, nil
}

func hasRequiredVLANInputFields(body map[string]interface{}) bool {
	for _, key := range vlanRequiredInputFields {
		if _, ok := body[key]; !ok {
			return false
		}
	}
	return true
}

func vlanCreateDefaults() map[string]interface{} {
	return map[string]interface{}{
		"enabled": "yes",
		"mac":     "",
		"ip_mask": "",
		"comment": "",
	}
}

var vlanRequiredInputFields = []string{
	"vlan_id",
	"vlan_name",
	"interface",
	"ip_addr",
	"netmask",
	"enabled",
	"mac",
	"ip_mask",
	"comment",
}

var vlanInputFields = []string{
	"vlan_id",
	"vlan_name",
	"tagname",
	"interface",
	"mac",
	"ip_addr",
	"netmask",
	"ip_mask",
	"enabled",
	"comment",
}

type natGroupOpts struct {
	createDefaults      map[string]interface{}
	addrFields          map[string]string
	defaultColumns      []string
	requiredCreateFlags []string
	inputFields         []string
}

func natGroup(app *cliapp.Runtime, use, short, apiPath string, fieldMap map[string]string, opts ...natGroupOpts) *cobra.Command {
	grp := &cobra.Command{Use: use, Short: short}
	var o natGroupOpts
	if len(opts) > 0 {
		o = opts[0]
	}

	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List " + use + " rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if len(o.defaultColumns) > 0 {
				app.DefaultColumns = o.defaultColumns
			}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/"+apiPath,
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(listCmd)

	createCmd := natWriteCmd(app, "create", "Create a "+use+" rule", false, apiPath, fieldMap, o.createDefaults, o.addrFields, nil,
		o.requiredCreateFlags,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Post(cliapp.APIBase+"/"+apiPath, body)
		})

	updateCmd := natWriteCmd(app, "update ID", "Update a "+use+" rule", true, apiPath, fieldMap, o.createDefaults, o.addrFields, o.inputFields,
		nil,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Put(cliapp.APIBase+"/"+apiPath+"/"+id, body)
		})

	toggleFieldMap := map[string]string{"enabled": "enabled"}
	toggleCmd := &cobra.Command{
		Use:   "toggle ID",
		Short: "Enable/disable a " + use + " rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "enabled"); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, toggleFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Patch(cliapp.APIBase+"/"+apiPath+"/"+args[0], body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	toggleCmd.Flags().String("data", "{}", "JSON body")
	cliapp.AddEnabledFlag(toggleCmd)
	cliapp.MarkFlagsRequired(toggleCmd, "enabled")

	deleteCmd := newDeleteCmd(app, "Delete a "+use+" rule", "/"+apiPath+"/")

	grp.AddCommand(listCmd, createCmd, updateCmd, toggleCmd, deleteCmd)
	return grp
}

// natWriteCmd builds a create or update command for natGroup.
// Supports fieldMap (semantic flags), createDefaults, and addrFields (nested objects).
func natWriteCmd(app *cliapp.Runtime, use, short string, withID bool, apiPath string, fieldMap map[string]string, defaults map[string]interface{}, addrFields map[string]string, inputFields []string, requiredCreateFlags []string, fn func(body interface{}, id string) (json.RawMessage, error)) *cobra.Command {
	c := &cobra.Command{Use: use, Short: short}
	if use == "create" {
		c.Aliases = []string{"new"}
	}
	if withID {
		c.Args = cobra.ExactArgs(1)
	}
	c.Flags().String("data", "{}", "JSON body")

	// Register address/port flags.
	for flagName := range addrFields {
		desc := natFlagDescs[flagName]
		if desc == "" {
			desc = "Comma-separated " + flagName + " values"
		}
		c.Flags().String(flagName, "", desc)
	}

	if fieldMap != nil {
		for flagName := range fieldMap {
			if flagName == "enabled" {
				continue
			}
			desc := natFlagDescs[flagName]
			if desc == "" {
				desc = flagName + " value"
			}
			c.Flags().String(flagName, "", desc)
		}
		cliapp.AddEnabledFlag(c)

		c.RunE = func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if !withID && len(requiredCreateFlags) > 0 {
				if err := cliapp.RequireFlags(cmd, requiredCreateFlags...); err != nil {
					return err
				}
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, fieldMap)
			if err != nil {
				return err
			}
			if !withID {
				for k, v := range defaults {
					if _, exists := body[k]; !exists {
						body[k] = v
					}
				}
			}
			for flagName, apiField := range addrFields {
				f := cmd.Flags().Lookup(flagName)
				if f == nil || !f.Changed {
					continue
				}
				parts := strings.Split(f.Value.String(), ",")
				custom := make([]interface{}, 0, len(parts))
				for _, v := range parts {
					v = strings.TrimSpace(v)
					if v != "" {
						custom = append(custom, v)
					}
				}
				body[apiField] = map[string]interface{}{
					"custom": custom,
					"object": []interface{}{},
				}
			}
			id := ""
			if withID {
				id = args[0]
				body, err = buildNATUpdateBody(app, apiPath, id, body, defaults, inputFields)
				if err != nil {
					return err
				}
			}
			raw, err := fn(body, id)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		}
	} else {
		c.RunE = func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.ParseJSON(data)
			if err != nil {
				return err
			}
			id := ""
			if withID {
				id = args[0]
			}
			raw, err := fn(body, id)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		}
	}
	if !withID && len(requiredCreateFlags) > 0 {
		cliapp.MarkFlagsRequired(c, requiredCreateFlags...)
	}
	return c
}

func buildNATUpdateBody(app *cliapp.Runtime, apiPath, id string, updates map[string]interface{}, defaults map[string]interface{}, inputFields []string) (map[string]interface{}, error) {
	if len(inputFields) == 0 {
		return updates, nil
	}

	readClient := app.APIClient
	if app.APIClient.DryRun {
		readClient = app.NewClient(app.Session.BaseURL, app.Session.Token)
	}
	raw, err := readClient.Get(cliapp.APIBase+"/"+apiPath+"/"+id, nil)
	if err != nil {
		if hasInputFields(updates, inputFields) {
			return updates, nil
		}
		return nil, err
	}

	current, err := extractInputObject(raw, inputFields)
	if err != nil {
		if hasInputFields(updates, inputFields) {
			return updates, nil
		}
		return nil, err
	}
	for k, v := range updates {
		current[k] = v
	}
	for k, v := range defaults {
		if _, exists := current[k]; !exists {
			current[k] = v
		}
	}
	return current, nil
}

func extractInputObject(raw json.RawMessage, inputFields []string) (map[string]interface{}, error) {
	var payload interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	current, err := findFirstAPIObject(payload)
	if err != nil {
		return nil, err
	}
	body := make(map[string]interface{})
	for _, key := range inputFields {
		if v, ok := current[key]; ok {
			body[key] = v
		}
	}
	return body, nil
}

func findFirstAPIObject(v interface{}) (map[string]interface{}, error) {
	switch data := v.(type) {
	case map[string]interface{}:
		if results, ok := data["results"]; ok {
			return findFirstAPIObject(results)
		}
		if arr, ok := data["data"].([]interface{}); ok {
			if len(arr) == 0 {
				return nil, &cliapp.ValidationError{Message: "empty API response"}
			}
			return findFirstAPIObject(arr[0])
		}
		return data, nil
	case []interface{}:
		if len(data) == 0 {
			return nil, &cliapp.ValidationError{Message: "empty API response"}
		}
		return findFirstAPIObject(data[0])
	default:
		return nil, &cliapp.ValidationError{Message: "unexpected API response"}
	}
}

func hasInputFields(body map[string]interface{}, inputFields []string) bool {
	for _, key := range inputFields {
		if _, ok := body[key]; !ok {
			return false
		}
	}
	return true
}

func newDeleteCmd(app *cliapp.Runtime, short, apiPath string) *cobra.Command {
	c := &cobra.Command{
		Use:     "delete ID",
		Aliases: []string{"rm"},
		Short:   short,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			yes, _ := cmd.Flags().GetBool("yes")
			resource := cmd.Parent().Use
			if err := cliapp.ConfirmDelete(app.Stdout, app.Stderr, resource, args[0], yes); err != nil {
				return err
			}
			raw, err := app.APIClient.Delete(cliapp.APIBase + apiPath + args[0])
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	c.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	return c
}
