package security

import (
	"encoding/json"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

// Field maps for semantic flags on write commands.
var (
	aclFieldMap = map[string]string{
		"name":          "tagname",
		"action":        "action",
		"direction":     "dir",
		"protocol":      "protocol",
		"in-interface":  "iinterface",
		"out-interface": "ointerface",
		"priority":      "prio",
		"comment":       "comment",
		"enabled":       "enabled",
	}
	// aclCreateDefaults provides required fields the iKuai API expects
	// but are not exposed as flags (sensible defaults for most use cases).
	aclCreateDefaults = map[string]interface{}{
		"ip_type":       "4",
		"ctdir":         0,
		"dir":           "forward",
		"iinterface":    "any",
		"ointerface":    "any",
		"prio":          50,
		"comment":       "",
		"src_addr":      map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"dst_addr":      map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"src_port":      map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"dst_port":      map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"src_addr_inv":  0,
		"dst_addr_inv":  0,
		"src_type":      0,
		"dst_type":      0,
		"src_area_code": "",
		"dst_area_code": "",
		"src6_addr":     map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"dst6_addr":     map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"src6_mode":     0,
		"dst6_mode":     0,
		"src6_suffix":   "",
		"dst6_suffix":   "",
		"time": map[string]interface{}{
			"custom": []interface{}{
				map[string]interface{}{
					"type":       "weekly",
					"weekdays":   "1234567",
					"start_time": "00:00",
					"end_time":   "23:59",
					"comment":    "",
				},
			},
			"object": []interface{}{},
		},
	}
	// aclAddrFields maps CLI flag names to API field names for address/port nested objects.
	// These flags accept comma-separated values and build {"custom": [...], "object": []} structures.
	aclAddrFields = map[string]string{
		"src-addr": "src_addr",
		"dst-addr": "dst_addr",
		"src-port": "src_port",
		"dst-port": "dst_port",
	}
	l7FieldMap = map[string]string{
		"name":     "tagname",
		"action":   "action",
		"priority": "prio",
		"comment":  "comment",
		"enabled":  "enabled",
	}
	l7CreateDefaults = map[string]interface{}{
		"comment":  "",
		"src_addr": "",
		"dst_addr": "",
		"time": map[string]interface{}{
			"custom": []interface{}{
				map[string]interface{}{
					"type":       "weekly",
					"weekdays":   "1234567",
					"start_time": "00:00",
					"end_time":   "23:59",
					"comment":    "",
				},
			},
			"object": []interface{}{},
		},
	}
	l7AddrFields = map[string]string{
		"app-proto": "app_proto",
	}
	macFieldMap = map[string]string{
		"name":    "tagname",
		"mac":     "mac",
		"comment": "comment",
		"enabled": "enabled",
	}
	macCreateDefaults = map[string]interface{}{
		"expires": 0,
		"comment": "",
		"time": map[string]interface{}{
			"custom": []interface{}{
				map[string]interface{}{
					"type":       "weekly",
					"weekdays":   "1234567",
					"start_time": "00:00",
					"end_time":   "23:59",
					"comment":    "",
				},
			},
			"object": []interface{}{},
		},
	}
	domainBlacklistFieldMap = map[string]string{
		"name":         "tagname",
		"domain-group": "domain_group",
		"comment":      "comment",
		"enabled":      "enabled",
	}
	peerconnFieldMap = map[string]string{
		"name":     "tagname",
		"limits":   "limits",
		"protocol": "protocol",
		"comment":  "comment",
		"enabled":  "enabled",
	}
	peerconnCreateDefaults = map[string]interface{}{
		"comment":  "",
		"src_addr": map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"dst_port": map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
	}
	peerconnAddrFields = map[string]string{
		"src-addr": "src_addr",
		"dst-port": "dst_port",
	}
	terminalsFieldMap = map[string]string{
		"name":    "tagname",
		"mac":     "mac",
		"comment": "comment",
	}
	urlKeywordsFieldMap = map[string]string{
		"name":        "tagname",
		"priority":    "prio",
		"mode":        "mode",
		"hit-rate":    "hit_rate",
		"src-url":     "src_url",
		"ori-keyword": "ori_keyword",
		"rep-keyword": "rep_keyword",
		"comment":     "comment",
		"enabled":     "enabled",
	}
	urlRedirectFieldMap = map[string]string{
		"name":     "tagname",
		"priority": "prio",
		"mode":     "mode",
		"hit-rate": "hit_rate",
		"src-url":  "src_url",
		"dst-url":  "dst_url",
		"comment":  "comment",
		"enabled":  "enabled",
	}
	urlReplaceFieldMap = map[string]string{
		"name":          "tagname",
		"priority":      "prio",
		"mode":          "mode",
		"hit-rate":      "hit_rate",
		"src-url":       "src_url",
		"param-keyword": "param_keyword",
		"rep-keyword":   "rep_keyword",
		"comment":       "comment",
		"enabled":       "enabled",
	}
	urlBlackFieldMap = map[string]string{
		"name":    "tagname",
		"mode":    "mode",
		"comment": "comment",
		"enabled": "enabled",
	}
	urlBlackCreateDefaults = map[string]interface{}{
		"comment":  "",
		"src_addr": map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"time":     map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
	}
	// urlBlackAddrFields maps CLI flag names to API field names for nested objects.
	// --domain and --src-addr accept comma-separated values → {"custom": [...], "object": []}.
	urlBlackAddrFields = map[string]string{
		"domain":   "domain",
		"src-addr": "src_addr",
	}

	// mac set-mode: acl_mac (0=blacklist, 1=whitelist)
	macModeFieldMap = map[string]string{
		"acl-mac": "acl_mac",
	}

	// advanced-set fields
	advancedFieldMap = map[string]string{
		"noping-lan":  "noping_lan",
		"noping-wan":  "noping_wan",
		"notracert":   "notracert",
		"hijack-ping": "hijack_ping",
		"invalid":     "invalid",
		"dos-lan":     "dos_lan",
		"dos-lan-num": "dos_lan_num",
		"tcp-mss":     "tcp_mss",
		"tcp-mss-num": "tcp_mss_num",
	}

	// secondary-route-set fields
	secondaryRouteFieldMap = map[string]string{
		"nol2rt":  "nol2rt",
		"ttl-num": "ttl_num",
	}
)

func New(app *cliapp.Runtime) *cobra.Command {
	securityCmd := &cobra.Command{
		Use:   "security",
		Short: "Security rules",
		Long:  `Manage security rules: ACL, MAC filtering, L7 application rules, URL filtering, domain blacklist, connection limits, and terminals.`,
		Example: `  ikuai-cli security acl list
  ikuai-cli security acl create --name "block_ssh" --action drop --protocol tcp --dst-port "22" --enabled yes
  ikuai-cli security mac list
  ikuai-cli security l7 list`,
	}

	securityCmd.AddCommand(secGroup(app, "acl", "IP ACL rules", "security/acl-rules", true, aclFieldMap, secGroupOpts{
		createDefaults: aclCreateDefaults,
		defaultColumns: []string{"id", "src_addr", "dst_addr", "src_port", "dst_port", "tagname", "action", "protocol", "enabled"},
		addrFields:     aclAddrFields,
	}))

	macCmd := &cobra.Command{Use: "mac", Short: "MAC access control"}
	macCmd.AddCommand(
		&cobra.Command{
			Use:   "get-mode",
			Short: "Get MAC filter mode",
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := app.RequireAuth(); err != nil {
					return err
				}
				raw, err := app.APIClient.Get(cliapp.APIBase+"/security/mac-mode", nil)
				if err != nil {
					return err
				}
				app.PrintRaw(raw)
				return nil
			},
		},
		dataCmdWithFlags(app, "set-mode", "Set MAC filter mode", func(c *cobra.Command) {
			c.Flags().String("acl-mac", "", "MAC mode (0=blacklist, 1=whitelist)")
		}, macModeFieldMap, func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Put(cliapp.APIBase+"/security/mac-mode", body)
		}),
	)
	addSecCRUD(app, macCmd, "security/mac-rules", false, macFieldMap, macCreateDefaults, nil, nil)
	securityCmd.AddCommand(macCmd)

	securityCmd.AddCommand(secGroup(app, "l7", "L7 application-layer rules", "security/app-protocols/professional/rules", true, l7FieldMap, secGroupOpts{
		createDefaults: l7CreateDefaults,
		addrFields:     l7AddrFields,
	}))

	urlCmd := &cobra.Command{Use: "url", Short: "URL filter rules"}
	urlCmd.AddCommand(
		secGroup(app, "black", "URL blacklist rules", "security/url-black/rules", false, urlBlackFieldMap, secGroupOpts{
			createDefaults: urlBlackCreateDefaults,
			addrFields:     urlBlackAddrFields,
		}),
		secGroup(app, "keywords", "URL keyword rules", "security/url-keywords/rules", false, urlKeywordsFieldMap, secGroupOpts{
			defaultColumns: []string{"id", "tagname", "mode", "src_url", "ori_keyword", "rep_keyword", "hit_rate", "prio", "enabled"},
		}),
		secGroup(app, "redirect", "URL redirect rules", "security/url-redirect/rules", false, urlRedirectFieldMap, secGroupOpts{
			defaultColumns: []string{"id", "tagname", "mode", "src_url", "dst_url", "hit_rate", "prio", "enabled"},
		}),
		secGroup(app, "replace", "URL replace rules", "security/url-replace/rules", false, urlReplaceFieldMap, secGroupOpts{
			defaultColumns: []string{"id", "tagname", "mode", "src_url", "param_keyword", "rep_keyword", "hit_rate", "prio", "enabled"},
		}),
	)
	securityCmd.AddCommand(urlCmd)

	securityCmd.AddCommand(secGroup(app, "domain-blacklist", "Domain blacklist rules", "security/domain-blacklist/rules", false, domainBlacklistFieldMap))
	securityCmd.AddCommand(secGroup(app, "peerconn", "Peer connection rules", "security/peerconn/rules", false, peerconnFieldMap, secGroupOpts{
		createDefaults: peerconnCreateDefaults,
		addrFields:     peerconnAddrFields,
	}))
	securityCmd.AddCommand(secGroup(app, "terminals", "Terminal device annotations", "security/terminals", false, terminalsFieldMap))

	securityCmd.AddCommand(
		dataCmd(app, "advanced-get", "Get advanced security config", func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Get(cliapp.APIBase+"/security/advanced/config", nil)
		}),
		dataCmdWithFlags(app, "advanced-set", "Set advanced security config", func(c *cobra.Command) {
			c.Flags().String("noping-lan", "", "Block LAN ping (0/1)")
			c.Flags().String("noping-wan", "", "Block WAN ping (0/1)")
			c.Flags().String("notracert", "", "Block tracert (0/1)")
			c.Flags().String("hijack-ping", "", "Hijack ping (0/1)")
			c.Flags().String("invalid", "", "Block invalid links (0/1)")
			c.Flags().String("dos-lan", "", "LAN DoS protection (0/1)")
			c.Flags().String("dos-lan-num", "", "LAN DoS connection limit")
			c.Flags().String("tcp-mss", "", "Enable TCP MSS (0/1)")
			c.Flags().String("tcp-mss-num", "", "TCP MSS max segment size")
		}, advancedFieldMap, func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Put(cliapp.APIBase+"/security/advanced/config", body)
		}),
		dataCmd(app, "secondary-route-get", "Get secondary (L2) route config", func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Get(cliapp.APIBase+"/security/secondary-route/config", nil)
		}),
		dataCmdWithFlags(app, "secondary-route-set", "Set secondary (L2) route config", func(c *cobra.Command) {
			c.Flags().String("nol2rt", "", "Block secondary routers (0=allow, 1=block)")
			c.Flags().String("nol2rt-ip", "", "IP addresses (comma-separated)")
			c.Flags().String("ttl-num", "", "Custom TTL value (1-255)")
		}, secondaryRouteFieldMap, func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Put(cliapp.APIBase+"/security/secondary-route/config", body)
		}),
	)

	return securityCmd
}

type secGroupOpts struct {
	createDefaults map[string]interface{}
	defaultColumns []string
	addrFields     map[string]string // CLI flag → API field for address/port nested objects.
}

func secGroup(app *cliapp.Runtime, use, short, apiPath string, withGet bool, fieldMap map[string]string, opts ...secGroupOpts) *cobra.Command {
	grp := &cobra.Command{Use: use, Short: short}
	var o secGroupOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	addSecCRUD(app, grp, apiPath, withGet, fieldMap, o.createDefaults, o.defaultColumns, o.addrFields)
	return grp
}

func addSecCRUD(app *cliapp.Runtime, grp *cobra.Command, apiPath string, withGet bool, fieldMap map[string]string, createDefaults map[string]interface{}, defaultColumns []string, addrFields map[string]string) {
	name := grp.Use
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List " + name + " rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if len(defaultColumns) > 0 {
				app.DefaultColumns = defaultColumns
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

	createCmd := secWriteCmd(app, "create", "Create a "+name+" rule", false, fieldMap, createDefaults, addrFields, func(body interface{}, id string) (json.RawMessage, error) {
		return app.APIClient.Post(cliapp.APIBase+"/"+apiPath, body)
	})
	updateCmd := secWriteCmd(app, "update ID", "Update a "+name+" rule", true, fieldMap, nil, addrFields, func(body interface{}, id string) (json.RawMessage, error) {
		return app.APIClient.Put(cliapp.APIBase+"/"+apiPath+"/"+id, body)
	})

	grp.AddCommand(
		listCmd,
		createCmd,
		updateCmd,
		dataCmdWithID(app, "toggle ID", "Enable/disable a "+name+" rule (--data JSON)", func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Patch(cliapp.APIBase+"/"+apiPath+"/"+id, body)
		}),
		deleteByIDCmd(app, "delete ID", "Delete a "+name+" rule", "/"+apiPath+"/"),
	)

	if withGet {
		grp.AddCommand(getByIDCmd(app, "get ID", "Get a single "+name+" rule", "/"+apiPath+"/"))
	}
}

// secWriteCmd builds a create or update command. When fieldMap is non-nil,
// semantic flags are registered and MergeDataWithFlags is used; otherwise
// it falls back to the plain --data / ParseJSON path.
// addrFields maps CLI flag names to API field names for address/port nested objects;
// these flags accept comma-separated values and build {"custom": [...], "object": []} structures.
func secWriteCmd(app *cliapp.Runtime, use, short string, withID bool, fieldMap map[string]string, defaults map[string]interface{}, addrFields map[string]string, fn callWithBody) *cobra.Command {
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
		addSemanticFlags(c, fieldMap)
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
			// Apply defaults: only fill keys not already set by --data or flags.
			for k, v := range defaults {
				if _, exists := body[k]; !exists {
					body[k] = v
				}
			}
			// Parse address/port flags into nested objects.
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

// addSemanticFlags registers a --<flag> string flag for every key in fieldMap,
// skipping "enabled" (handled by AddEnabledFlag).
func addSemanticFlags(cmd *cobra.Command, fieldMap map[string]string) {
	for flagName, apiField := range fieldMap {
		if flagName == "enabled" {
			continue
		}
		cmd.Flags().String(flagName, "", "Set "+apiField+" field")
	}
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

func dataCmd(app *cliapp.Runtime, use, short string, fn callWithBody) *cobra.Command {
	return dataCmdWithFlags(app, use, short, nil, nil, fn)
}

func dataCmdWithFlags(app *cliapp.Runtime, use, short string, addFlags func(*cobra.Command), fieldMap map[string]string, fn callWithBody) *cobra.Command {
	return dataCmdImpl(app, use, short, false, addFlags, fieldMap, fn)
}

func dataCmdWithID(app *cliapp.Runtime, use, short string, fn callWithBody) *cobra.Command {
	return dataCmdImpl(app, use, short, true, nil, nil, fn)
}

func dataCmdImpl(app *cliapp.Runtime, use, short string, withID bool, addFlags func(*cobra.Command), fieldMap map[string]string, fn callWithBody) *cobra.Command {
	c := &cobra.Command{Use: use, Short: short}
	if withID {
		c.Args = cobra.ExactArgs(1)
	}
	isToggle := strings.HasPrefix(use, "toggle")
	if isToggle {
		cliapp.AddEnabledFlag(c)
	}
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
		switch {
		case isToggle:
			body, err = cliapp.MergeDataWithFlags(data, cmd, map[string]string{"enabled": "enabled"})
		case fieldMap != nil:
			body, err = cliapp.MergeDataWithFlags(data, cmd, fieldMap)
		default:
			body, err = cliapp.ParseJSON(data)
		}
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
	if use != "advanced-get" && use != "secondary-route-get" {
		c.Flags().String("data", "{}", "JSON body")
	}
	return c
}
