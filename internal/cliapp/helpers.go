package cliapp

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const APIBase = "/api/v4.0"

func ParseJSON(s string) (interface{}, error) {
	if s == "" || s == "{}" {
		return map[string]interface{}{}, nil
	}
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return nil, &ValidationError{Message: fmt.Sprintf("invalid JSON: %v", err)}
	}
	return v, nil
}

func ListParams(page, pageSize int, filter, order, orderBy string) map[string]string {
	p := map[string]string{
		"page":      fmt.Sprint(page),
		"page_size": fmt.Sprint(pageSize),
	}
	if filter != "" {
		p["filter"] = filter
	}
	if order != "" {
		p["order"] = order
	}
	if orderBy != "" {
		p["order_by"] = orderBy
	}
	return p
}

func AddListFlags(cmd *cobra.Command) {
	cmd.Flags().IntP("page", "p", 1, "Page number")
	cmd.Flags().Int("page-size", 20, "Items per page")
	// --limit is planned but not yet implemented; omitted from registration to avoid misleading users.
	cmd.Flags().String("filter", "", "Filter: field==value, & for AND, comma for OR")
	cmd.Flags().String("order", "", "Sort direction: asc|desc")
	cmd.Flags().String("order-by", "", "Sort field")
	_ = cmd.RegisterFlagCompletionFunc("order", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"asc", "desc"}, cobra.ShellCompDirectiveNoFileComp
	})
}

func GetListParams(cmd *cobra.Command) (page, pageSize int, filter, order, orderBy string) {
	page, _ = cmd.Flags().GetInt("page")
	pageSize, _ = cmd.Flags().GetInt("page-size")
	filter, _ = cmd.Flags().GetString("filter")
	order, _ = cmd.Flags().GetString("order")
	orderBy, _ = cmd.Flags().GetString("order-by")
	return
}

// integerAPIFields lists API field names that expect integer values.
// Only these fields are auto-coerced from string to int64 by MergeDataWithFlags.
// All other fields are sent as strings to avoid corrupting values that happen
// to look numeric (e.g. SNMP community "161", VLAN name "0010").
var integerAPIFields = map[string]bool{
	"sshd_port": true, "http_port": true, "https_port": true,
	"server_port": true, "open_mppe": true, "force_ipsec": true,
	"local_listenport": true, "endpoint_port": true, "method": true,
	"ikelifetime": true, "lifetime": true, "dpddelay": true, "dpdtimeout": true,
	"keepalive": true,
	"prio":      true, "priority": true,
	"syn_recv_timeout": true, "established_timeout": true,
	"upload": true, "download": true,
	"ftp_ports": true, "sip_ports": true, "tftp_ports": true,
	"limits": true, "hit_rate": true,
	"authmode": true, "authport": true, "accountport": true,
	"force_verify_name": true, "rate_limit_lan": true, "drop_client": true,
	"force_pppoe": true, "enhance_check": true, "share_deny": true,
	"bind_vlan": true, "verify_vlan": true, "bind_iface": true,
	"mtu": true, "mru": true, "lcp_echo_interval": true,
	"lcp_echo_failure": true, "maxconnect": true, "restart_timer": true,
	"acl_mac": true, "mode": true,
	"type": true, "iface_band": true, "dst_type": true,
	"src_addr_inv": true, "dst_addr_inv": true,
	"noping_lan": true, "noping_wan": true, "notracert": true,
	"hijack_ping": true, "invalid": true, "dos_lan": true,
	"dos_lan_num": true, "tcp_mss": true, "tcp_mss_num": true,
	"nol2rt": true, "ttl_num": true,
}

// MergeDataWithFlags parses --data JSON, then overlays any explicitly set
// semantic flags (from fieldMap). Flags override data fields.
// fieldMap key = CLI flag name, value = API field name.
func MergeDataWithFlags(dataJSON string, cmd *cobra.Command, fieldMap map[string]string) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	if dataJSON != "" && dataJSON != "{}" {
		if err := json.Unmarshal([]byte(dataJSON), &result); err != nil {
			return nil, &ValidationError{Message: fmt.Sprintf("invalid --data JSON: %v", err)}
		}
	}
	for flagName, apiField := range fieldMap {
		f := cmd.Flags().Lookup(flagName)
		if f == nil || !f.Changed {
			continue
		}
		val := f.Value.String()
		// Only coerce to int64 for known integer API fields.
		if integerAPIFields[apiField] {
			if n, err := strconv.ParseInt(val, 10, 64); err == nil {
				result[apiField] = n
				continue
			}
		}
		result[apiField] = val
	}
	return result, nil
}

// RequireFlags returns a ValidationError if any of the named flags were not
// explicitly set by the user. Call this inside RunE before making API calls.
func RequireFlags(cmd *cobra.Command, flags ...string) error {
	var missing []string
	for _, name := range flags {
		f := cmd.Flags().Lookup(name)
		if f == nil || !f.Changed {
			missing = append(missing, "--"+name)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	if len(missing) == 1 {
		return &ValidationError{Message: "missing required flag: " + missing[0]}
	}
	return &ValidationError{Message: "missing required flags: " + strings.Join(missing, ", ")}
}

// MarkFlagsRequired appends " (required)" to the usage text of each named flag.
// Call after all flags are registered so --help shows which flags are mandatory.
func MarkFlagsRequired(cmd *cobra.Command, flags ...string) {
	for _, name := range flags {
		f := cmd.Flags().Lookup(name)
		if f != nil && !strings.HasSuffix(f.Usage, "(required)") {
			f.Usage += " (required)"
		}
	}
}

// AddEnabledFlag adds --enabled flag accepting "yes" or "no".
func AddEnabledFlag(cmd *cobra.Command) {
	cmd.Flags().String("enabled", "", "Enable/disable: yes|no")
	_ = cmd.RegisterFlagCompletionFunc("enabled", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"yes", "no"}, cobra.ShellCompDirectiveNoFileComp
	})
}
