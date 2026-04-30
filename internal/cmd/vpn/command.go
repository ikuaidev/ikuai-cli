package vpn

import (
	"encoding/json"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

// flagDescs provides human-readable descriptions for CLI flags.
var flagDescs = map[string]string{
	"name":         "Tunnel/client name",
	"server":       "VPN server address",
	"remote-addr":  "Remote server address",
	"username":     "Authentication username",
	"password":     "Authentication password",
	"interface":    "Bound interface (e.g. wan1, auto)",
	"comment":      "Comment",
	"left-id":      "Local identity string",
	"left-subnet":  "Local subnet (CIDR)",
	"right-subnet": "Remote subnet (CIDR)",
	"secret":       "Pre-shared key",
	"address":      "Tunnel address (CIDR)",
	"private-key":  "WireGuard private key (base64)",
	"public-key":   "WireGuard public key (base64)",
	"allow-ips":    "Allowed IPs (CIDR, comma-separated)",
	"port":         "Listen port",
	"mtu":          "MTU size",
	"keepalive":    "Persistent keepalive interval (seconds)",
	"ca":           "CA certificate (PEM)",
}

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
		"enabled":           "yes",
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
		"enabled":           "yes",
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
		"enabled":           "yes",
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

	pptpClientUpdateInputFields = []string{
		"id", "enabled", "name", "comment", "server", "server_port", "username", "passwd",
		"interface", "mtu", "mru", "check_link_mode", "check_link_host", "timing_rst_switch",
		"timing_rst_week", "timing_rst_time", "cycle_rst_time",
	}
	l2tpClientUpdateInputFields = []string{
		"id", "enabled", "name", "comment", "server", "server_port", "username", "passwd",
		"ipsec_secret", "interface", "leftid", "rightid", "mtu", "mru", "check_link_mode",
		"check_link_host", "timing_rst_switch", "timing_rst_week", "timing_rst_time", "cycle_rst_time",
	}
	openvpnClientUpdateInputFields = []string{
		"id", "enabled", "name", "comment", "remote_addr", "remote_port", "method", "username",
		"password", "interface", "proto", "dev_type", "cipher", "tls_auth", "ca", "cert", "key",
		"accept_push_route", "route", "comp_lzo", "tun_mtu", "check_link_mode", "check_link_host",
		"timing_rst_switch", "timing_rst_week", "timing_rst_time", "extra_config",
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
		"enabled":         "yes",
		"comment":         "",
		"authby":          "mschapv2",
		"check_link_mode": 3,
		"check_link_host": "www.baidu.com",
	}
	ikev2ClientUpdateInputFields = []string{
		"id", "enabled", "name", "comment", "remote_addr", "interface", "authby", "secret",
		"leftid", "rightid", "username", "passwd", "check_link_mode", "check_link_host",
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
		"enabled":     "yes",
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
	ipsecClientUpdateInputFields = []string{
		"id", "name", "comment", "remote_addr", "authby", "leftsubnet", "rightsubnet",
		"interface", "enabled", "keyexchange", "aggressive", "ikelifetime", "ike_enc",
		"ike_auth", "ike_dh", "secret", "leftid", "rightid", "privatekey", "leftcert",
		"rightcert", "lifetime", "esp_enc", "esp_auth", "dpdaction", "dpddelay",
		"dpdtimeout", "compress",
	}

	// WireGuard tunnel
	wireguardFieldMap = map[string]string{
		"name":      "name",
		"interface": "interface",
		"address":   "local_address",
		"port":      "local_listenport",
		"enabled":   "enabled",
	}
	wireguardDefaults = map[string]interface{}{
		"enabled":          "yes",
		"interface":        "auto",
		"local_listenport": 5000,
		"mtu":              1420,
	}
	wireguardInputFields = []string{
		"enabled", "name", "interface", "local_privatekey", "local_publickey", "local_address",
		"local_listenport", "mtu",
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
		"enabled":   "yes",
		"comment":   "",
		"keepalive": 0,
	}
	wireguardPeerInputFields = []string{
		"enabled", "comment", "interface", "peer_publickey", "presharedkey", "allowips",
		"endpoint", "endpoint_port", "keepalive",
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
	pptpServerInputFields = []string{
		"enabled", "dns1", "dns2", "addr_pool", "open_mppe", "server_ip", "server_port", "mtu", "mru",
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
	l2tpServerInputFields = []string{
		"enabled", "server_ip", "server_port", "addr_pool", "dns1", "dns2", "mtu", "mru",
		"ipsec_secret", "leftid", "rightid", "force_ipsec",
	}

	openvpnServerFieldMap = map[string]string{
		"enabled":            "enabled",
		"proto":              "proto",
		"port":               "port",
		"server-port":        "port",
		"subnet":             "subnet",
		"mask":               "mask",
		"tun-mtu":            "tun_mtu",
		"mtu":                "tun_mtu",
		"cipher":             "cipher",
		"comp-lzo":           "comp_lzo",
		"dev-type":           "dev_type",
		"topology":           "topology",
		"method":             "method",
		"tls-auth":           "tls_auth",
		"ca":                 "ca",
		"cert":               "cert",
		"key":                "key",
		"push-gateway":       "push_gateway",
		"push-route":         "push_route",
		"push-route-comment": "push_route_comment",
		"push-dns":           "push_dns",
		"extra-config":       "extra_config",
	}
	openvpnServerInputFields = []string{
		"enabled", "proto", "port", "subnet", "mask", "tun_mtu", "cipher", "comp_lzo",
		"dev_type", "topology", "method", "tls_auth", "ca", "cert", "key", "push_gateway",
		"push_route", "push_route_comment", "push_dns", "extra_config",
	}

	ikev2ServerFieldMap = map[string]string{
		"enabled":     "enabled",
		"authby":      "authby",
		"addrpool":    "addrpool",
		"addr-pool":   "addrpool",
		"secret":      "secret",
		"leftid":      "leftid",
		"rightid":     "rightid",
		"dns1":        "dns1",
		"dns2":        "dns2",
		"share-deny":  "share_deny",
		"mtu":         "mtu",
		"private-key": "privatekey",
		"left-cert":   "leftcert",
	}
	ikev2ServerInputFields = []string{
		"id", "enabled", "authby", "addrpool", "secret", "leftid", "rightid", "dns1",
		"dns2", "share_deny", "mtu", "privatekey", "leftcert", "aggressive", "keyexchange", "name",
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
		[]string{"id", "name", "server", "username", "interface", "enabled"}, pptpServerFieldMap,
		pptpServerInputFields, pptpClientUpdateInputFields, []string{"name", "server", "username", "password", "interface"}))
	vpnCmd.AddCommand(vpnServerGroup(app, "l2tp", "vpn/l2tp", l2tpClientFieldMap, l2tpClientDefaults,
		[]string{"id", "name", "server", "username", "interface", "enabled"}, l2tpServerFieldMap,
		l2tpServerInputFields, l2tpClientUpdateInputFields, []string{"name", "server", "username", "password", "interface"}))
	vpnCmd.AddCommand(vpnServerGroup(app, "openvpn", "vpn/openvpn", openvpnClientFieldMap, openvpnClientDefaults,
		[]string{"id", "name", "remote_addr", "remote_port", "proto", "interface", "enabled"}, openvpnServerFieldMap,
		openvpnServerInputFields, openvpnClientUpdateInputFields, []string{"name", "remote-addr", "interface", "username", "password", "ca"}))
	vpnCmd.AddCommand(vpnServerGroup(app, "ikev2", "vpn/ikev2", ikev2ClientFieldMap, ikev2ClientDefaults,
		[]string{"id", "name", "remote_addr", "interface", "authby", "enabled"}, ikev2ServerFieldMap,
		ikev2ServerInputFields, ikev2ClientUpdateInputFields, []string{"name", "remote-addr", "interface", "left-id"}))
	vpnCmd.AddCommand(ipsecGroup(app))
	vpnCmd.AddCommand(wireguardGroup(app))

	return vpnCmd
}

func vpnServerGroup(app *cliapp.Runtime, name, apiPath string, clientFieldMap map[string]string, clientDefaults map[string]interface{}, defaultColumns []string, serverFieldMap map[string]string, serverInputFields []string, clientUpdateInputFields []string, requiredCreateFlags []string) *cobra.Command {
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
	}, serverFieldMap, serverInputFields, cliapp.APIBase+"/"+apiPath+"/services")

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
			params := vpnListParams(cmd, page, pageSize, filter, order, orderBy)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/"+apiPath+"/clients",
				params)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	addVPNListFlags(clientsCmd)

	clientCreateCmd := writeCmd(app, "client-create", "Create a "+name+" client", false, clientFieldMap, nil, clientDefaults,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Post(cliapp.APIBase+"/"+apiPath+"/clients", body)
		})
	if len(requiredCreateFlags) > 0 {
		cliapp.MarkFlagsRequired(clientCreateCmd, requiredCreateFlags...)
		origRunE := clientCreateCmd.RunE
		clientCreateCmd.RunE = func(cmd *cobra.Command, args []string) error {
			if err := cliapp.RequireFlags(cmd, requiredCreateFlags...); err != nil {
				return err
			}
			return origRunE(cmd, args)
		}
	}
	clientGetCmd := getByIDCmd(app, "client-get ID", "Get a "+name+" client", "/"+apiPath+"/clients/")
	clientUpdateCmd := updateByIDCmd(app, "client-update ID", "Update a "+name+" client", clientFieldMap, clientUpdateInputFields, cliapp.APIBase+"/"+apiPath+"/clients/")
	clientToggleCmd := writeCmd(app, "client-toggle ID", "Enable/disable a "+name+" client", true, map[string]string{"enabled": "enabled"}, nil, nil,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Patch(cliapp.APIBase+"/"+apiPath+"/clients/"+id, body)
		})
	clientDeleteCmd := deleteByIDCmd(app, "client-delete ID", "Delete a "+name+" client", "/"+apiPath+"/clients/")
	kickCmd := deleteByIDCmd(app, "kick ID", "Kick / delete a "+name+" client", "/"+apiPath+"/clients/")

	grp.AddCommand(getCmd, setCmd, clientsCmd, clientGetCmd, clientCreateCmd, clientUpdateCmd, clientToggleCmd, clientDeleteCmd, kickCmd)
	return grp
}

func addVPNListFlags(cmd *cobra.Command) {
	cliapp.AddListFlags(cmd)
	cmd.Flags().String("key", "", "Fuzzy match fields, comma-separated")
	cmd.Flags().String("pattern", "", "Fuzzy match pattern")
}

func vpnListParams(cmd *cobra.Command, page, pageSize int, filter, order, orderBy string) map[string]string {
	params := cliapp.ListParamsWithPageSizeKey(page, pageSize, filter, order, orderBy, "limit")
	if key, _ := cmd.Flags().GetString("key"); key != "" {
		params["key"] = key
	}
	if pattern, _ := cmd.Flags().GetString("pattern"); pattern != "" {
		params["pattern"] = pattern
	}
	return params
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
			params := vpnListParams(cmd, page, pageSize, filter, order, orderBy)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/vpn/ipsec/clients",
				params)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	addVPNListFlags(clientsCmd)

	createCmd := writeCmd(app, "client-create", "Create an IPSec client", false, ipsecClientFieldMap, nil, ipsecClientDefaults,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Post(cliapp.APIBase+"/vpn/ipsec/clients", body)
		})
	cliapp.MarkFlagsRequired(createCmd, "name", "interface", "left-subnet", "right-subnet")
	{
		origRunE := createCmd.RunE
		createCmd.RunE = func(cmd *cobra.Command, args []string) error {
			if err := cliapp.RequireFlags(cmd, "name", "interface", "left-subnet", "right-subnet"); err != nil {
				return err
			}
			return origRunE(cmd, args)
		}
	}
	grp.AddCommand(
		clientsCmd,
		getByIDCmd(app, "client-get ID", "Get an IPSec client", "/vpn/ipsec/clients/"),
		createCmd,
		updateByIDCmd(app, "client-update ID", "Update an IPSec client", ipsecClientFieldMap, ipsecClientUpdateInputFields, cliapp.APIBase+"/vpn/ipsec/clients/"),
		writeCmd(app, "client-toggle ID", "Enable/disable an IPSec client", true, map[string]string{"enabled": "enabled"}, nil, nil,
			func(body interface{}, id string) (json.RawMessage, error) {
				return app.APIClient.Patch(cliapp.APIBase+"/vpn/ipsec/clients/"+id, body)
			}),
		deleteByIDCmd(app, "client-delete ID", "Delete an IPSec client", "/vpn/ipsec/clients/"),
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
			params := vpnListParams(cmd, page, pageSize, filter, order, orderBy)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/vpn/wireguard",
				params)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	addVPNListFlags(listCmd)

	peersListCmd := &cobra.Command{
		Use:   "peers ID",
		Short: "List peers",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			params := vpnListParams(cmd, page, pageSize, filter, order, orderBy)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/vpn/wireguard/"+args[0]+"/peers",
				params)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	addVPNListFlags(peersListCmd)

	// WireGuard tunnel needs private/public keys — these must be passed via flags.
	wgCreateFieldMap := make(map[string]string, len(wireguardFieldMap)+2)
	for k, v := range wireguardFieldMap {
		wgCreateFieldMap[k] = v
	}
	wgCreateFieldMap["private-key"] = "local_privatekey"
	wgCreateFieldMap["public-key"] = "local_publickey"

	wgCreateCmd := writeCmd(app, "create", "Create a WireGuard tunnel", false, wgCreateFieldMap, nil, wireguardDefaults,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Post(cliapp.APIBase+"/vpn/wireguard", body)
		})
	cliapp.MarkFlagsRequired(wgCreateCmd, "name", "address", "private-key", "public-key")
	{
		origRunE := wgCreateCmd.RunE
		wgCreateCmd.RunE = func(cmd *cobra.Command, args []string) error {
			if err := cliapp.RequireFlags(cmd, "name", "address", "private-key", "public-key"); err != nil {
				return err
			}
			return origRunE(cmd, args)
		}
	}
	wgUpdateCmd := updateByIDCmd(app, "update ID", "Update a WireGuard tunnel", wireguardFieldMap, wireguardInputFields, cliapp.APIBase+"/vpn/wireguard/")
	cliapp.MarkFlagsRequired(wgUpdateCmd, "interface")
	{
		origRunE := wgUpdateCmd.RunE
		wgUpdateCmd.RunE = func(cmd *cobra.Command, args []string) error {
			if err := cliapp.RequireFlags(cmd, "interface"); err != nil {
				return err
			}
			return origRunE(cmd, args)
		}
	}
	peerCreateCmd := writeCmd(app, "peer-create ID", "Create a peer on a WireGuard tunnel", true, wireguardPeerFieldMap, nil, wireguardPeerDefaults,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Post(cliapp.APIBase+"/vpn/wireguard/"+id+"/peers", body)
		})
	cliapp.MarkFlagsRequired(peerCreateCmd, "public-key", "allow-ips", "interface")
	{
		origRunE := peerCreateCmd.RunE
		peerCreateCmd.RunE = func(cmd *cobra.Command, args []string) error {
			if err := cliapp.RequireFlags(cmd, "public-key", "allow-ips", "interface"); err != nil {
				return err
			}
			return origRunE(cmd, args)
		}
	}
	grp.AddCommand(
		listCmd,
		getByIDCmd(app, "get ID", "Get a WireGuard tunnel", "/vpn/wireguard/"),
		wgCreateCmd,
		wgUpdateCmd,
		writeCmd(app, "toggle ID", "Enable/disable a WireGuard tunnel", true, map[string]string{"enabled": "enabled"}, nil, nil,
			func(body interface{}, id string) (json.RawMessage, error) {
				return app.APIClient.Patch(cliapp.APIBase+"/vpn/wireguard/"+id, body)
			}),
		deleteByIDCmd(app, "delete ID", "Delete a WireGuard tunnel", "/vpn/wireguard/"),
		peersListCmd,
		peerGetCmd(app),
		peerCreateCmd,
		peerUpdateCmd(app),
		peerToggleCmd(app),
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

func updateByIDCmd(app *cliapp.Runtime, use, short string, fieldMap map[string]string, inputFields []string, apiPathPrefix string) *cobra.Command {
	c := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
	}
	addBodyFlags(c, fieldMap)
	c.RunE = func(cmd *cobra.Command, args []string) error {
		if err := app.RequireAuth(); err != nil {
			return err
		}
		data, _ := cmd.Flags().GetString("data")
		body, err := buildFullBody(app, cmd, data, fieldMap, inputFields, apiPathPrefix+args[0], args[0])
		if err != nil {
			return err
		}
		raw, err := app.APIClient.Put(apiPathPrefix+args[0], body)
		if err != nil {
			return err
		}
		app.PrintRaw(raw)
		return nil
	}
	return c
}

func peerGetCmd(app *cliapp.Runtime) *cobra.Command {
	return &cobra.Command{
		Use:   "peer-get TUNNEL_ID PEER_ID",
		Short: "Get a WireGuard peer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/vpn/wireguard/"+args[0]+"/peers/"+args[1], nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
}

func peerUpdateCmd(app *cliapp.Runtime) *cobra.Command {
	c := &cobra.Command{
		Use:   "peer-update TUNNEL_ID PEER_ID",
		Short: "Update a WireGuard peer",
		Args:  cobra.ExactArgs(2),
	}
	addBodyFlags(c, wireguardPeerFieldMap)
	c.RunE = func(cmd *cobra.Command, args []string) error {
		if err := app.RequireAuth(); err != nil {
			return err
		}
		path := cliapp.APIBase + "/vpn/wireguard/" + args[0] + "/peers/" + args[1]
		data, _ := cmd.Flags().GetString("data")
		body, err := buildFullBody(app, cmd, data, wireguardPeerFieldMap, wireguardPeerInputFields, path, args[1])
		if err != nil {
			return err
		}
		raw, err := app.APIClient.Put(path, body)
		if err != nil {
			return err
		}
		app.PrintRaw(raw)
		return nil
	}
	return c
}

func peerToggleCmd(app *cliapp.Runtime) *cobra.Command {
	c := &cobra.Command{
		Use:   "peer-toggle TUNNEL_ID PEER_ID",
		Short: "Enable/disable a WireGuard peer",
		Args:  cobra.ExactArgs(2),
	}
	c.Flags().String("data", "{}", "JSON body (escape hatch)")
	cliapp.AddEnabledFlag(c)
	cliapp.MarkFlagsRequired(c, "enabled")
	c.RunE = func(cmd *cobra.Command, args []string) error {
		if err := app.RequireAuth(); err != nil {
			return err
		}
		if err := cliapp.RequireFlags(cmd, "enabled"); err != nil {
			return err
		}
		data, _ := cmd.Flags().GetString("data")
		body, err := cliapp.MergeDataWithFlags(data, cmd, map[string]string{"enabled": "enabled"})
		if err != nil {
			return err
		}
		raw, err := app.APIClient.Patch(cliapp.APIBase+"/vpn/wireguard/"+args[0]+"/peers/"+args[1], body)
		if err != nil {
			return err
		}
		app.PrintRaw(raw)
		return nil
	}
	return c
}

func addBodyFlags(c *cobra.Command, fieldMap map[string]string) {
	c.Flags().String("data", "{}", "JSON body (escape hatch)")
	for flagName := range fieldMap {
		if flagName == "enabled" {
			continue
		}
		desc := flagDescs[flagName]
		if desc == "" {
			desc = flagName + " value"
		}
		c.Flags().String(flagName, "", desc)
	}
	cliapp.AddEnabledFlag(c)
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
		desc := flagDescs[flagName]
		if desc == "" {
			desc = flagName + " value"
		}
		c.Flags().String(flagName, "", desc)
	}
	for flagName := range addrFields {
		desc := flagDescs[flagName]
		if desc == "" {
			desc = "Comma-separated " + flagName + " values"
		}
		c.Flags().String(flagName, "", desc)
	}
	cliapp.AddEnabledFlag(c)
	if isToggleUse(use) {
		cliapp.MarkFlagsRequired(c, "enabled")
	}

	c.RunE = func(cmd *cobra.Command, args []string) error {
		if err := app.RequireAuth(); err != nil {
			return err
		}
		if isToggleUse(use) {
			if err := cliapp.RequireFlags(cmd, "enabled"); err != nil {
				return err
			}
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

func isToggleUse(use string) bool {
	action := strings.Fields(use)
	return len(action) > 0 && strings.Contains(action[0], "toggle")
}

// configCmdWithFlags builds a full-body config PUT command with optional semantic flags.
func configCmdWithFlags(app *cliapp.Runtime, use, short string, addFlags func(*cobra.Command), fieldMap map[string]string, inputFields []string, apiPath string) *cobra.Command {
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
		body, err := buildFullBody(app, cmd, data, fieldMap, inputFields, apiPath, "")
		if err != nil {
			return err
		}
		raw, err := app.APIClient.Put(apiPath, body)
		if err != nil {
			return err
		}
		app.PrintRaw(raw)
		return nil
	}
	return c
}

func buildFullBody(app *cliapp.Runtime, cmd *cobra.Command, data string, fieldMap map[string]string, inputFields []string, getPath string, fallbackID string) (map[string]interface{}, error) {
	changes, err := cliapp.MergeDataWithFlags(data, cmd, fieldMap)
	if err != nil {
		return nil, err
	}
	readClient := app.APIClient
	if app.APIClient.DryRun {
		readClient = app.NewClient(app.Session.BaseURL, app.Session.Token)
	}
	raw, err := readClient.Get(getPath, nil)
	if err != nil {
		if hasAllInputFields(changes, inputFields) {
			return changes, nil
		}
		return nil, err
	}
	current, err := extractVPNInputObject(raw, inputFields)
	if err != nil {
		if hasAllInputFields(changes, inputFields) {
			return changes, nil
		}
		return nil, err
	}
	for k, v := range changes {
		current[k] = v
	}
	if fallbackID != "" && hasInputField(inputFields, "id") {
		if _, exists := current["id"]; !exists {
			current["id"] = fallbackID
		}
	}
	return current, nil
}

func extractVPNInputObject(raw json.RawMessage, inputFields []string) (map[string]interface{}, error) {
	var v interface{}
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, err
	}
	obj, err := findFirstVPNObject(v)
	if err != nil {
		return nil, err
	}
	if len(inputFields) == 0 {
		return obj, nil
	}
	out := map[string]interface{}{}
	for _, key := range inputFields {
		if val, ok := obj[key]; ok {
			out[key] = val
		}
	}
	return out, nil
}

func findFirstVPNObject(v interface{}) (map[string]interface{}, error) {
	switch data := v.(type) {
	case map[string]interface{}:
		if rows, ok := data["data"].([]interface{}); ok {
			return findFirstVPNObject(rows)
		}
		if rows, ok := data["results"].([]interface{}); ok {
			return findFirstVPNObject(rows)
		}
		if rows, ok := data["iface_data"].([]interface{}); ok {
			return findFirstVPNObject(rows)
		}
		return data, nil
	case []interface{}:
		if len(data) == 0 {
			return nil, &cliapp.ValidationError{Message: "empty VPN response"}
		}
		return findFirstVPNObject(data[0])
	default:
		return nil, &cliapp.ValidationError{Message: "unexpected VPN response"}
	}
}

func hasAllInputFields(body map[string]interface{}, inputFields []string) bool {
	for _, key := range inputFields {
		if _, ok := body[key]; !ok {
			return false
		}
	}
	return true
}

func hasInputField(fields []string, want string) bool {
	for _, field := range fields {
		if field == want {
			return true
		}
	}
	return false
}
