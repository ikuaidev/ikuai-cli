package cliapp

import (
	"encoding/json"
	"fmt"
	"strconv"

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
	"prio": true, "priority": true,
	"syn_recv_timeout": true, "established_timeout": true,
	"upload": true, "download": true,
	"ftp_ports": true, "sip_ports": true, "tftp_ports": true,
	"limits": true, "hit_rate": true,
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

// AddEnabledFlag adds --enabled flag accepting "yes" or "no".
func AddEnabledFlag(cmd *cobra.Command) {
	cmd.Flags().String("enabled", "", "Enable/disable: yes|no")
	_ = cmd.RegisterFlagCompletionFunc("enabled", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"yes", "no"}, cobra.ShellCompDirectiveNoFileComp
	})
}
