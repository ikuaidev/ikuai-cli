package vpn

import (
	"encoding/json"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

// --- VPN client field maps and defaults ---

var (
	pptpClientFieldMap = map[string]string{
		"name":      "name",
		"server":    "server",
		"username":  "username",
		"password":  "passwd",
		"interface": "interface",
		"comment":   "comment",
		"enabled":   "enabled",
	}
	pptpClientDefaults = map[string]interface{}{
		"comment":           "",
		"server_port":       1723,
		"mtu":               1400,
		"mru":               1400,
		"check_link_mode":   3,
		"check_link_host":   "www.baidu.com",
		"timing_rst_switch": 0,
		"timing_rst_week":   "1234567",
		"timing_rst_time":   "12:00",
		"cycle_rst_time":    0,
	}

	l2tpClientFieldMap = map[string]string{
		"name":      "name",
		"server":    "server",
		"username":  "username",
		"password":  "passwd",
		"interface": "interface",
		"comment":   "comment",
		"enabled":   "enabled",
	}
	l2tpClientDefaults = map[string]interface{}{
		"comment":           "",
		"server_port":       1701,
		"mtu":               1400,
		"mru":               1400,
		"check_link_mode":   3,
		"check_link_host":   "www.baidu.com",
		"timing_rst_switch": 0,
		"cycle_rst_time":    0,
	}

	openvpnClientFieldMap = map[string]string{
		"name":        "name",
		"remote-addr": "remote_addr",
		"remote-port": "remote_port",
		"username":    "username",
		"password":    "password",
		"interface":   "interface",
		"proto":       "proto",
		"cipher":      "cipher",
		"ca":          "ca",
		"comment":     "comment",
		"enabled":     "enabled",
	}
	openvpnClientDefaults = map[string]interface{}{
		"comment":           "",
		"remote_port":       "1194",
		"method":            0,
		"proto":             "udp",
		"dev_type":          "tun",
		"cipher":            "BF-CBC",
		"tun_mtu":           "1400",
		"comp_lzo":          "1",
		"accept_push_route": 1,
		"check_link_mode":   3,
		"check_link_host":   "www.baidu.com",
		"timing_rst_switch": 0,
		"ca":                "",
	}

	// IKEv2 client
	ikev2ClientFieldMap = map[string]string{
		"name":        "name",
		"remote-addr": "remote_addr",
		"interface":   "interface",
		"authby":      "authby",
		"left-id":     "leftid",
		"right-id":    "rightid",
		"secret":      "secret",
		"username":    "username",
		"password":    "passwd",
		"comment":     "comment",
		"enabled":     "enabled",
	}
	ikev2ClientDefaults = map[string]interface{}{
		"comment":         "",
		"authby":          "mschapv2",
		"check_link_mode": 3,
		"check_link_host": "www.baidu.com",
	}

	// IPSec client
	ipsecClientFieldMap = map[string]string{
		"name":         "name",
		"remote-addr":  "remote_addr",
		"interface":    "interface",
		"left-subnet":  "leftsubnet",
		"right-subnet": "rightsubnet",
		"secret":       "secret",
		"comment":      "comment",
		"enabled":      "enabled",
	}
	ipsecClientDefaults = map[string]interface{}{
		"comment":     "",
		"keyexchange": "ikev2",
		"authby":      "secret",
		"dpdaction":   "none",
		"compress":    "0",
		"ikelifetime": 3,
		"lifetime":    1,
		"ike_enc":     "aes256",
		"ike_auth":    "sha256",
		"ike_dh":      "modp2048",
		"esp_enc":     "aes256",
		"esp_auth":    "sha256",
		"aggressive":  "0",
	}

	// WireGuard tunnel
	wireguardFieldMap = map[string]string{
		"name":      "name",
		"interface": "interface",
		"address":   "local_address",
		"port":      "local_listenport",
		"comment":   "comment",
		"enabled":   "enabled",
	}
	wireguardDefaults = map[string]interface{}{
		"interface":        "auto",
		"local_listenport": 5000,
		"mtu":              1420,
		"keepalive":        0,
	}

	// WireGuard peer
	wireguardPeerFieldMap = map[string]string{
		"public-key": "peer_publickey",
		"allow-ips":  "allowips",
		"endpoint":   "endpoint",
		"port":       "endpoint_port",
		"interface":  "interface",
		"comment":    "comment",
		"enabled":    "enabled",
	}
	wireguardPeerDefaults = map[string]interface{}{
		"comment":   "",
		"keepalive": 0,
	}

	// --- VPN server field maps ---

	pptpServerFieldMap = map[string]string{
		"enabled":     "enabled",
		"server-ip":   "server_ip",
		"server-port": "server_port",
		"addr-pool":   "addr_pool",
		"dns1":        "dns1",
		"dns2":        "dns2",
		"open-mppe":   "open_mppe",
		"mtu":         "mtu",
		"mru":         "mru",
	}

	l2tpServerFieldMap = map[string]string{
		"enabled":      "enabled",
		"server-ip":    "server_ip",
		"server-port":  "server_port",
		"addr-pool":    "addr_pool",
		"dns1":         "dns1",
		"dns2":         "dns2",
		"mtu":          "mtu",
		"mru":          "mru",
		"ipsec-secret": "ipsec_secret",
		"leftid":       "leftid",
		"rightid":      "rightid",
		"force-ipsec":  "force_ipsec",
	}

	openvpnServerFieldMap = map[string]string{
		"enabled":     "enabled",
		"server-ip":   "server_ip",
		"server-port": "server_port",
		"addr-pool":   "addr_pool",
		"dns1":        "dns1",
		"dns2":        "dns2",
		"proto":       "proto",
		"mtu":         "mtu",
	}

	ikev2ServerFieldMap = map[string]string{
		"enabled":   "enabled",
		"server-ip": "server_ip",
		"addr-pool": "addr_pool",
		"dns1":      "dns1",
		"dns2":      "dns2",
	}
)

func New(app *cliapp.Runtime) *cobra.Command {
	vpnCmd := &cobra.Command{
		Use:   "vpn",
		Short: "VPN management",
		Long:  `Manage VPN services: PPTP, L2TP, OpenVPN, IKEv2, IPSec, and WireGuard — server config, client CRUD, tunnels, and peers.`,
		Example: `  ikuai-cli vpn pptp get
  ikuai-cli vpn pptp clients
  ikuai-cli vpn wireguard list
  ikuai-cli vpn ipsec clients`,
	}

	vpnCmd.AddCommand(vpnServerGroup(app, "pptp", "vpn/pptp", pptpClientFieldMap, pptpClientDefaults,
		[]string{"id", "name", "server", "username", "interface", "enabled"}, pptpServerFieldMap))
	vpnCmd.AddCommand(vpnServerGroup(app, "l2tp", "vpn/l2tp", l2tpClientFieldMap, l2tpClientDefaults,
		[]string{"id", "name", "server", "username", "interface", "enabled"}, l2tpServerFieldMap))
	vpnCmd.AddCommand(vpnServerGroup(app, "openvpn", "vpn/openvpn", openvpnClientFieldMap, openvpnClientDefaults,
		[]string{"id", "name", "remote_addr", "remote_port", "proto", "interface", "enabled"}, openvpnServerFieldMap))
	vpnCmd.AddCommand(vpnServerGroup(app, "ikev2", "vpn/ikev2", ikev2ClientFieldMap, ikev2ClientDefaults,
		[]string{"id", "name", "remote_addr", "interface", "authby", "enabled"}, ikev2ServerFieldMap))
	vpnCmd.AddCommand(ipsecGroup(app))
	vpnCmd.AddCommand(wireguardGroup(app))

	return vpnCmd
}

func vpnServerGroup(app *cliapp.Runtime, name, apiPath string, clientFieldMap map[string]string, clientDefaults map[string]interface{}, defaultColumns []string, serverFieldMap map[string]string) *cobra.Command {
	grp := &cobra.Command{Use: name, Short: name + " VPN"}

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get " + name + " server configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/"+apiPath+"/services", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	setCmd := configCmdWithFlags(app, "set", "Update "+name+" server configuration", func(c *cobra.Command) {
		for flagName := range serverFieldMap {
			c.Flags().String(flagName, "", flagName)
		}
	}, serverFieldMap, func(body interface{}, id string) (json.RawMessage, error) {
		return app.APIClient.Put(cliapp.APIBase+"/"+apiPath+"/services", body)
	})

	clientsCmd := &cobra.Command{
		Use:   "clients",
		Short: "List " + name + " client accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if len(defaultColumns) > 0 {
				app.DefaultColumns = defaultColumns
			}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/"+apiPath+"/clients",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(clientsCmd)

	clientCreateCmd := writeCmd(app, "client-create", "Create a "+name+" client", false, clientFieldMap, nil, clientDefaults,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Post(cliapp.APIBase+"/"+apiPath+"/clients", body)
		})
	clientUpdateCmd := writeCmd(app, "client-update ID", "Update a "+name+" client", true, clientFieldMap, nil, nil,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Put(cliapp.APIBase+"/"+apiPath+"/clients/"+id, body)
		})
	kickCmd := deleteByIDCmd(app, "kick ID", "Kick / delete a "+name+" client", "/"+apiPath+"/clients/")

	grp.AddCommand(getCmd, setCmd, clientsCmd, clientCreateCmd, clientUpdateCmd, kickCmd)
	return grp
}

func ipsecGroup(app *cliapp.Runtime) *cobra.Command {
	grp := &cobra.Command{Use: "ipsec", Short: "IPSec clients"}

	clientsCmd := &cobra.Command{
		Use:   "clients",
		Short: "List active IPSec clients",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "name", "remote_addr", "leftsubnet", "rightsubnet", "interface", "keyexchange", "authby", "enabled"}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/vpn/ipsec/clients",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(clientsCmd)

	grp.AddCommand(
		clientsCmd,
		writeCmd(app, "client-create", "Create an IPSec client", false, ipsecClientFieldMap, nil, ipsecClientDefaults,
			func(body interface{}, id string) (json.RawMessage, error) {
				return app.APIClient.Post(cliapp.APIBase+"/vpn/ipsec/clients", body)
			}),
		writeCmd(app, "client-update ID", "Update an IPSec client", true, ipsecClientFieldMap, nil, nil,
			func(body interface{}, id string) (json.RawMessage, error) {
				return app.APIClient.Put(cliapp.APIBase+"/vpn/ipsec/clients/"+id, body)
			}),
		deleteByIDCmd(app, "kick ID", "Kick an active IPSec client", "/vpn/ipsec/clients/"),
	)
	return grp
}

func wireguardGroup(app *cliapp.Runtime) *cobra.Command {
	grp := &cobra.Command{Use: "wireguard", Short: "WireGuard tunnels & peers"}

	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List tunnels",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "name", "local_address", "local_listenport", "interface", "mtu", "enabled"}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/vpn/wireguard",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(listCmd)

	peersListCmd := &cobra.Command{
		Use:   "peers ID",
		Short: "List peers",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/vpn/wireguard/"+args[0]+"/peers",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(peersListCmd)

	// WireGuard tunnel needs private/public keys — these must be passed via flags.
	wgCreateFieldMap := make(map[string]string, len(wireguardFieldMap)+2)
	for k, v := range wireguardFieldMap {
		wgCreateFieldMap[k] = v
	}
	wgCreateFieldMap["private-key"] = "local_privatekey"
	wgCreateFieldMap["public-key"] = "local_publickey"

	grp.AddCommand(
		listCmd,
		getByIDCmd(app, "get ID", "Get a WireGuard tunnel", "/vpn/wireguard/"),
		writeCmd(app, "create", "Create a WireGuard tunnel", false, wgCreateFieldMap, nil, wireguardDefaults,
			func(body interface{}, id string) (json.RawMessage, error) {
				return app.APIClient.Post(cliapp.APIBase+"/vpn/wireguard", body)
			}),
		writeCmd(app, "update ID", "Update a WireGuard tunnel", true, wireguardFieldMap, nil, nil,
			func(body interface{}, id string) (json.RawMessage, error) {
				return app.APIClient.Put(cliapp.APIBase+"/vpn/wireguard/"+id, body)
			}),
		writeCmd(app, "toggle ID", "Enable/disable a WireGuard tunnel", true, map[string]string{"enabled": "enabled"}, nil, nil,
			func(body interface{}, id string) (json.RawMessage, error) {
				return app.APIClient.Patch(cliapp.APIBase+"/vpn/wireguard/"+id, body)
			}),
		deleteByIDCmd(app, "delete ID", "Delete a WireGuard tunnel", "/vpn/wireguard/"),
		peersListCmd,
		writeCmd(app, "peer-create ID", "Create a peer on a WireGuard tunnel", true, wireguardPeerFieldMap, nil, wireguardPeerDefaults,
			func(body interface{}, id string) (json.RawMessage, error) {
				return app.APIClient.Post(cliapp.APIBase+"/vpn/wireguard/"+id+"/peers", body)
			}),
		peerDeleteCmd(app),
	)
	return grp
}

func getByIDCmd(app *cliapp.Runtime, use, short, apiPath string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+apiPath+args[0], nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
}

func peerDeleteCmd(app *cliapp.Runtime) *cobra.Command {
	c := &cobra.Command{
		Use:   "peer-delete TUNNEL_ID PEER_ID",
		Short: "Delete peer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			yes, _ := cmd.Flags().GetBool("yes")
			if err := cliapp.ConfirmDelete(app.Stdout, app.Stderr, "peer", args[1], yes); err != nil {
				return err
			}
			raw, err := app.APIClient.Delete(cliapp.APIBase + "/vpn/wireguard/" + args[0] + "/peers/" + args[1])
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

func deleteByIDCmd(app *cliapp.Runtime, use, short, apiPath string) *cobra.Command {
	c := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
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
	if use == "delete ID" {
		c.Aliases = []string{"rm"}
	}
	c.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	return c
}

type callWithBody func(body interface{}, id string) (json.RawMessage, error)

// writeCmd builds a create/update command with semantic flags + defaults.
func writeCmd(app *cliapp.Runtime, use, short string, withID bool, fieldMap map[string]string, addrFields map[string]string, defaults map[string]interface{}, fn callWithBody) *cobra.Command {
	c := &cobra.Command{Use: use, Short: short}
	if use == "create" {
		c.Aliases = []string{"new"}
	}
	if withID {
		c.Args = cobra.ExactArgs(1)
	}
	c.Flags().String("data", "{}", "JSON body (escape hatch)")
	for flagName := range fieldMap {
		if flagName == "enabled" {
			continue
		}
		c.Flags().String(flagName, "", flagName+" value")
	}
	for flagName := range addrFields {
		c.Flags().String(flagName, "", "Comma-separated "+flagName+" values")
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
	return c
}

// configCmdWithFlags builds a --data config command with optional semantic flags.
func configCmdWithFlags(app *cliapp.Runtime, use, short string, addFlags func(*cobra.Command), fieldMap map[string]string, fn callWithBody) *cobra.Command {
	c := &cobra.Command{Use: use, Short: short}
	c.Flags().String("data", "{}", "JSON body")
	if addFlags != nil {
		addFlags(c)
	}
	c.RunE = func(cmd *cobra.Command, args []string) error {
		if err := app.RequireAuth(); err != nil {
			return err
		}
		data, _ := cmd.Flags().GetString("data")
		var body interface{}
		var err error
		if fieldMap != nil {
			body, err = cliapp.MergeDataWithFlags(data, cmd, fieldMap)
		} else {
			body, err = cliapp.ParseJSON(data)
		}
		if err != nil {
			return err
		}
		raw, err := fn(body, "")
		if err != nil {
			return err
		}
		app.PrintRaw(raw)
		return nil
	}
	return c
}
