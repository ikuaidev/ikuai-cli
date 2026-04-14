package network

import (
	"encoding/json"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

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

	networkWANCmd := &cobra.Command{
		Use:   "wan",
		Short: "WAN interface config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "tagname", "internet", "ip_mask", "gateway", "mac", "mtu", "speed", "enabled"}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/interfaces/wan-config", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkWANVLANCmd := &cobra.Command{
		Use:   "wan-vlan",
		Short: "WAN VLAN config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/interfaces/wan-vlan-config", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkLANCmd := &cobra.Command{
		Use:   "lan",
		Short: "LAN interface config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "tagname", "ip_mask", "bandeth", "dhcp_server", "vlan"}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/interfaces/lan-config", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkPhysicalCmd := &cobra.Command{
		Use:   "physical",
		Short: "Physical NIC info",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/interfaces/physical", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	networkDNSCmd := &cobra.Command{Use: "dns", Short: "DNS config"}

	networkDNSGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get DNS config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
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
			body, err := cliapp.MergeDataWithFlags(data, cmd, dnsSetFieldMap)
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
		"is_ipv6":  0,
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
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, dnsProxyFieldMap)
			if err != nil {
				return err
			}
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
			body, err := cliapp.MergeDataWithFlags(data, cmd, dhcpFieldMap)
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
			body, err := cliapp.MergeDataWithFlags(data, cmd, dhcpStaticFieldMap)
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
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.ParseJSON(data)
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
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.ParseJSON(data)
			if err != nil {
				return err
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
			body, err := cliapp.ParseJSON(data)
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
			if _, exists := body["enabled"]; !exists {
				body["enabled"] = "yes"
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
			body, err := cliapp.MergeDataWithFlags(data, cmd, vlanFieldMap)
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
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/pppoe/services", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	pppoeFieldMap := map[string]string{
		"enabled":     "enabled",
		"server-name": "server_name",
		"server-ip":   "server_ip",
		"dns1":        "dns1",
		"dns2":        "dns2",
		"addr-pool":   "addr_pool",
		"interface":   "interface",
		"authmode":    "authmode",
		"mtu":         "mtu",
		"mru":         "mru",
		"radius-ip":   "radius_ip",
		"secret":      "secret",
		"comment":     "comment",
	}
	networkPPPoESetCmd := &cobra.Command{
		Use:   "set",
		Short: "Set PPPoE config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, pppoeFieldMap)
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
	cliapp.AddEnabledFlag(networkDHCPStaticToggleCmd)
	// DHCP access-mode set semantic flags
	networkDHCPAccessModeSetCmd.Flags().String("mode", "", "Access mode (0=blacklist, 1=whitelist, 2=sync MAC ACL)")
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
	// DHCP static create/update semantic flags
	for _, c := range []*cobra.Command{networkDHCPStaticCreateCmd, networkDHCPStaticUpdateCmd} {
		c.Flags().String("name", "", "Binding name")
		c.Flags().String("ip", "", "IP address")
		c.Flags().String("mac", "", "MAC address")
		c.Flags().String("interface", "", "Interface name")
		c.Flags().String("gateway", "", "Gateway address")
		c.Flags().String("hostname", "", "Hostname")
		c.Flags().String("comment", "", "Comment")
		cliapp.AddEnabledFlag(c)
	}

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
	cliapp.AddEnabledFlag(networkDHCP6AccessRuleToggleCmd)

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
	// VLAN create/update semantic flags
	for _, c := range []*cobra.Command{networkVLANCreateCmd, networkVLANUpdateCmd} {
		c.Flags().String("name", "", "VLAN name")
		c.Flags().String("vlan-id", "", "VLAN ID")
		c.Flags().String("interface", "", "Interface name")
		c.Flags().String("ip", "", "IP address")
		c.Flags().String("netmask", "", "Subnet mask")
		c.Flags().String("mac", "", "MAC address")
		c.Flags().String("comment", "", "Comment")
		cliapp.AddEnabledFlag(c)
	}

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

	networkCmd.AddCommand(
		natGroup(app, "nat", "NAT rules", "network/nat/rules", natFieldMap, natGroupOpts{
			createDefaults: natCreateDefaults,
			addrFields:     natAddrFields,
			defaultColumns: []string{"id", "tagname", "action", "src_addr", "dst_addr", "protocol", "iinterface", "ointerface", "enabled"},
		}),
		natGroup(app, "dnat", "DNAT (port-forward) rules", "network/dnat/rules", dnatFieldMap),
		natGroup(app, "dmz", "DMZ rules", "network/dmz/rules", nil),
	)

	networkCmd.AddCommand(networkPPPoECmd)
	networkPPPoECmd.AddCommand(networkPPPoEGetCmd, networkPPPoESetCmd)
	networkPPPoESetCmd.Flags().String("data", "{}", "JSON body")
	networkPPPoESetCmd.Flags().String("enabled", "", "Service enabled (yes/no)")
	networkPPPoESetCmd.Flags().String("server-name", "", "Server name")
	networkPPPoESetCmd.Flags().String("server-ip", "", "Server IP address")
	networkPPPoESetCmd.Flags().String("dns1", "", "Primary DNS")
	networkPPPoESetCmd.Flags().String("dns2", "", "Secondary DNS")
	networkPPPoESetCmd.Flags().String("addr-pool", "", "Address pool range")
	networkPPPoESetCmd.Flags().String("interface", "", "Interface name")
	networkPPPoESetCmd.Flags().String("authmode", "", "Auth mode (0=local, 1=radius)")
	networkPPPoESetCmd.Flags().String("mtu", "", "MTU value")
	networkPPPoESetCmd.Flags().String("mru", "", "MRU value")
	networkPPPoESetCmd.Flags().String("radius-ip", "", "RADIUS server IP")
	networkPPPoESetCmd.Flags().String("secret", "", "RADIUS shared secret")
	networkPPPoESetCmd.Flags().String("comment", "", "Comment")

	return networkCmd
}

type natGroupOpts struct {
	createDefaults map[string]interface{}
	addrFields     map[string]string
	defaultColumns []string
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

	createCmd := natWriteCmd(app, "create", "Create a "+use+" rule", false, apiPath, fieldMap, o.createDefaults, o.addrFields,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Post(cliapp.APIBase+"/"+apiPath, body)
		})

	updateCmd := natWriteCmd(app, "update ID", "Update a "+use+" rule", true, apiPath, fieldMap, nil, o.addrFields,
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

	deleteCmd := newDeleteCmd(app, "Delete a "+use+" rule", "/"+apiPath+"/")

	grp.AddCommand(listCmd, createCmd, updateCmd, toggleCmd, deleteCmd)
	return grp
}

// natWriteCmd builds a create or update command for natGroup.
// Supports fieldMap (semantic flags), createDefaults, and addrFields (nested objects).
func natWriteCmd(app *cliapp.Runtime, use, short string, withID bool, apiPath string, fieldMap map[string]string, defaults map[string]interface{}, addrFields map[string]string, fn func(body interface{}, id string) (json.RawMessage, error)) *cobra.Command {
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
		c.Flags().String(flagName, "", "Comma-separated "+flagName+" values")
	}

	if fieldMap != nil {
		for flagName := range fieldMap {
			if flagName == "enabled" {
				continue
			}
			c.Flags().String(flagName, "", flagName+" value")
		}
		cliapp.AddEnabledFlag(c)

		c.RunE = func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, fieldMap)
			if err != nil {
				return err
			}
			for k, v := range defaults {
				if _, exists := body[k]; !exists {
					body[k] = v
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
	return c
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
