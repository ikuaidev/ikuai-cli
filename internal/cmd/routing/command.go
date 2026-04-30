package routing

import (
	"encoding/json"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

// flagDescs provides human-readable descriptions for CLI flags.
// writeCmd falls back to "flagName value" when a key is absent.
var flagDescs = map[string]string{
	// common
	"name":      "Rule name",
	"interface": "Outbound interface (e.g. wan1)",
	"priority":  "Priority (lower = higher)",
	"comment":   "Comment",
	// static / five-tuple shared
	"dst-addr": "Destination address",
	"gateway":  "Next-hop gateway IP",
	"ip-type":  "IP type (4/6)",
	"netmask":  "Subnet mask",
	// stream domain
	"domain": "Domain list (comma-separated)",
	// stream five-tuple
	"area-code":    "Destination area code",
	"dst-addr-inv": "Invert destination address match (0/1)",
	"dst-port":     "Destination port (comma-separated)",
	"dst-type":     "Destination type (0=address, 1=area)",
	"iface-band":   "Bind interface (0/1)",
	"nexthop":      "Next-hop gateway for LAN forwarding",
	"protocol":     "Protocol (tcp/udp/any)",
	"src-addr":     "Source address (comma-separated)",
	"src-addr-inv": "Invert source address match (0/1)",
	"src-port":     "Source port (comma-separated)",
	"type":         "Forwarding type (0=WAN, 1=LAN)",
	// stream l7
	"app-proto": "App protocol (comma-separated)",
	// stream load-balance
	"mode":     "Balance mode",
	"weight":   "Weight",
	"isp-name": "ISP filter (default: all)",
	// stream updown
	"upiface":   "Upstream interface",
	"downiface": "Downstream interface",
}

// --- Field maps, addr fields, and defaults for stream subcommands ---

var (
	domainFieldMap = map[string]string{
		"name":      "tagname",
		"interface": "interface",
		"priority":  "prio",
		"comment":   "comment",
		"enabled":   "enabled",
	}
	domainAddrFields = map[string]string{
		"domain":   "domain",
		"src-addr": "src_addr",
	}
	domainDefaults = map[string]interface{}{
		"enabled":  "yes",
		"comment":  "",
		"prio":     31,
		"domain":   map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"src_addr": map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"time": map[string]interface{}{
			"custom": []interface{}{
				map[string]interface{}{"type": "weekly", "weekdays": "1234567", "start_time": "00:00", "end_time": "23:59", "comment": ""},
			},
			"object": []interface{}{},
		},
	}

	fiveTupleFieldMap = map[string]string{
		"name":         "tagname",
		"interface":    "interface",
		"type":         "type",
		"nexthop":      "nexthop",
		"protocol":     "protocol",
		"mode":         "mode",
		"priority":     "prio",
		"dst-type":     "dst_type",
		"area-code":    "area_code",
		"iface-band":   "iface_band",
		"src-addr-inv": "src_addr_inv",
		"dst-addr-inv": "dst_addr_inv",
		"comment":      "comment",
		"enabled":      "enabled",
	}
	fiveTupleAddrFields = map[string]string{
		"src-addr": "src_addr",
		"dst-addr": "dst_addr",
		"src-port": "src_port",
		"dst-port": "dst_port",
	}
	fiveTupleDefaults = map[string]interface{}{
		"enabled":      "yes",
		"comment":      "",
		"type":         0,
		"mode":         0,
		"prio":         31,
		"iface_band":   0,
		"src_addr_inv": 0,
		"dst_addr_inv": 0,
		"protocol":     "any",
		"src_addr":     map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"dst_addr":     map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"src_port":     map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"dst_port":     map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"time": map[string]interface{}{
			"custom": []interface{}{
				map[string]interface{}{"type": "weekly", "weekdays": "1234567", "start_time": "00:00", "end_time": "23:59", "comment": ""},
			},
			"object": []interface{}{},
		},
	}

	l7FieldMap = map[string]string{
		"name":       "tagname",
		"interface":  "interface",
		"mode":       "mode",
		"priority":   "prio",
		"iface-band": "iface_band",
		"comment":    "comment",
		"enabled":    "enabled",
	}
	l7AddrFields = map[string]string{
		"app-proto": "app_proto",
		"src-addr":  "src_addr",
	}
	l7Defaults = map[string]interface{}{
		"enabled":    "yes",
		"comment":    "",
		"app_proto":  map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"mode":       0,
		"prio":       31,
		"iface_band": 0,
		"src_addr":   map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"time": map[string]interface{}{
			"custom": []interface{}{
				map[string]interface{}{"type": "weekly", "weekdays": "1234567", "start_time": "00:00", "end_time": "23:59", "comment": ""},
			},
			"object": []interface{}{},
		},
	}

	loadBalanceFieldMap = map[string]string{
		"name":      "tagname",
		"interface": "interface",
		"mode":      "mode",
		"weight":    "weight",
		"isp-name":  "isp_name",
		"comment":   "comment",
		"enabled":   "enabled",
	}
	loadBalanceDefaults = map[string]interface{}{
		"enabled":  "yes",
		"comment":  "",
		"mode":     0,
		"weight":   "1",
		"isp_name": "all",
	}

	updownFieldMap = map[string]string{
		"name":      "tagname",
		"upiface":   "upiface",
		"downiface": "downiface",
		"protocol":  "protocol",
		"comment":   "comment",
		"enabled":   "enabled",
	}
	updownAddrFields = map[string]string{
		"src-addr": "src_addr",
		"dst-addr": "dst_addr",
		"src-port": "src_port",
		"dst-port": "dst_port",
	}
	updownDefaults = map[string]interface{}{
		"enabled":  "yes",
		"comment":  "",
		"protocol": "any",
		"src_addr": map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"dst_addr": map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"src_port": map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"dst_port": map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
	}
)

func New(app *cliapp.Runtime) *cobra.Command {
	routingCmd := &cobra.Command{
		Use:   "routing",
		Short: "Routing & traffic shunting",
	}

	routingCmd.AddCommand(staticGroup(app))

	streamCmd := &cobra.Command{Use: "stream", Short: "Traffic shunting rules"}
	streamCmd.AddCommand(
		ruleGroup(app, "domain", "Domain-based stream rules", "routing/domain-rules", ruleGroupOpts{
			fieldMap:            domainFieldMap,
			addrFields:          domainAddrFields,
			createDefaults:      domainDefaults,
			defaultColumns:      []string{"id", "tagname", "domain", "interface", "src_addr", "prio", "enabled"},
			requiredCreateFlags: []string{"name", "interface"},
		}),
		ruleGroup(app, "five-tuple", "5-tuple stream rules (src/dst IP, port, protocol)", "routing/five-tuple-rules", ruleGroupOpts{
			fieldMap:            fiveTupleFieldMap,
			addrFields:          fiveTupleAddrFields,
			createDefaults:      fiveTupleDefaults,
			defaultColumns:      []string{"id", "tagname", "src_addr", "dst_addr", "protocol", "dst_port", "interface", "prio", "enabled"},
			requiredCreateFlags: []string{"name", "interface"},
		}),
		ruleGroup(app, "l7", "L7 application protocol stream rules", "routing/app-protocols", ruleGroupOpts{
			fieldMap:            l7FieldMap,
			addrFields:          l7AddrFields,
			createDefaults:      l7Defaults,
			defaultColumns:      []string{"id", "tagname", "app_proto", "interface", "src_addr", "prio", "enabled"},
			requiredCreateFlags: []string{"name", "interface"},
		}),
		ruleGroup(app, "load-balance", "Load balance rules", "routing/load-balance-rules", ruleGroupOpts{
			fieldMap:            loadBalanceFieldMap,
			createDefaults:      loadBalanceDefaults,
			defaultColumns:      []string{"id", "tagname", "interface", "mode", "weight", "isp_name", "enabled"},
			requiredCreateFlags: []string{"name", "interface"},
		}),
		ruleGroup(app, "updown", "Upstream/downstream rules", "routing/updown", ruleGroupOpts{
			fieldMap:            updownFieldMap,
			addrFields:          updownAddrFields,
			createDefaults:      updownDefaults,
			defaultColumns:      []string{"id", "tagname", "upiface", "downiface", "protocol", "src_addr", "enabled"},
			requiredCreateFlags: []string{"name", "upiface", "downiface"},
		}),
	)
	routingCmd.AddCommand(streamCmd)

	return routingCmd
}

func staticGroup(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "static", Short: "Static routes"}

	staticFieldMap := map[string]string{
		"name":      "tagname",
		"interface": "interface",
		"dst-addr":  "dst_addr",
		"gateway":   "gateway",
		"ip-type":   "ip_type",
		"netmask":   "netmask",
		"priority":  "prio",
		"comment":   "comment",
		"enabled":   "enabled",
	}
	staticDefaults := map[string]interface{}{
		"enabled": "yes",
		"ip_type": "4",
		"comment": "",
		"prio":    1,
	}

	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List static routes",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "tagname", "dst_addr", "gateway", "netmask", "interface", "prio", "enabled"}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/routing/static-routes",
				cliapp.ListParamsWithPageSizeKey(page, pageSize, filter, order, orderBy, "limit"))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(listCmd)

	createCmd := writeCmd(app, "create", "Create a static route", false, staticFieldMap, nil, staticDefaults, "",
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Post(cliapp.APIBase+"/routing/static-routes", body)
		})
	cliapp.MarkFlagsRequired(createCmd, "name", "dst-addr", "gateway", "netmask", "interface")
	{
		origRunE := createCmd.RunE
		createCmd.RunE = func(cmd *cobra.Command, args []string) error {
			if err := cliapp.RequireFlags(cmd, "name", "dst-addr", "gateway", "netmask", "interface"); err != nil {
				return err
			}
			return origRunE(cmd, args)
		}
	}
	getCmd := getByIDCmd(app, "get ID", "Get a static route", "/routing/static-routes/")
	updateCmd := writeCmd(app, "update ID", "Update a static route", true, staticFieldMap, nil, nil, "routing/static-routes",
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Put(cliapp.APIBase+"/routing/static-routes/"+id, body)
		})

	toggleFieldMap := map[string]string{"enabled": "enabled"}
	group.AddCommand(
		listCmd,
		getCmd,
		createCmd,
		updateCmd,
		writeCmd(app, "toggle ID", "Enable/disable a static route", true, toggleFieldMap, nil, nil, "",
			func(body interface{}, id string) (json.RawMessage, error) {
				return app.APIClient.Patch(cliapp.APIBase+"/routing/static-routes/"+id, body)
			}),
		deleteByIDCmd(app, "delete ID", "Delete a static route", "/routing/static-routes/"),
	)
	return group
}

type ruleGroupOpts struct {
	fieldMap            map[string]string
	addrFields          map[string]string
	createDefaults      map[string]interface{}
	defaultColumns      []string
	requiredCreateFlags []string
}

func ruleGroup(app *cliapp.Runtime, use, short, apiPath string, opts ruleGroupOpts) *cobra.Command {
	grp := &cobra.Command{Use: use, Short: short}
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List " + use + " rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if len(opts.defaultColumns) > 0 {
				app.DefaultColumns = opts.defaultColumns
			}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/"+apiPath,
				cliapp.ListParamsWithPageSizeKey(page, pageSize, filter, order, orderBy, "limit"))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(listCmd)

	createCmd := writeCmd(app, "create", "Create a "+use+" rule", false, opts.fieldMap, opts.addrFields, opts.createDefaults, "",
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Post(cliapp.APIBase+"/"+apiPath, body)
		})
	if len(opts.requiredCreateFlags) > 0 {
		cliapp.MarkFlagsRequired(createCmd, opts.requiredCreateFlags...)
		origRunE := createCmd.RunE
		createCmd.RunE = func(cmd *cobra.Command, args []string) error {
			if err := cliapp.RequireFlags(cmd, opts.requiredCreateFlags...); err != nil {
				return err
			}
			return origRunE(cmd, args)
		}
	}
	getCmd := getByIDCmd(app, "get ID", "Get a "+use+" rule", "/"+apiPath+"/")
	updateCmd := writeCmd(app, "update ID", "Update a "+use+" rule", true, opts.fieldMap, opts.addrFields, nil, apiPath,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Put(cliapp.APIBase+"/"+apiPath+"/"+id, body)
		})
	toggleFieldMap := map[string]string{"enabled": "enabled"}
	toggleCmd := writeCmd(app, "toggle ID", "Enable/disable a "+use+" rule", true, toggleFieldMap, nil, nil, "",
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Patch(cliapp.APIBase+"/"+apiPath+"/"+id, body)
		})

	grp.AddCommand(
		listCmd,
		getCmd,
		createCmd,
		updateCmd,
		toggleCmd,
		deleteByIDCmd(app, "delete ID", "Delete a "+use+" rule", "/"+apiPath+"/"),
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

// writeCmd builds a create/update/toggle command with semantic flags, addrFields, and defaults.
// --data is kept as escape hatch for complex cases.
func writeCmd(app *cliapp.Runtime, use, short string, withID bool, fieldMap map[string]string, addrFields map[string]string, defaults map[string]interface{}, fullUpdatePath string, fn callWithBody) *cobra.Command {
	c := &cobra.Command{Use: use, Short: short}
	if use == "create" {
		c.Aliases = []string{"new"}
	}
	if withID {
		c.Args = cobra.ExactArgs(1)
	}
	c.Flags().String("data", "{}", "JSON body (escape hatch)")

	// Register semantic flags from fieldMap.
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
	// Register addrFields flags.
	for flagName := range addrFields {
		desc := flagDescs[flagName]
		if desc == "" {
			desc = "Comma-separated " + flagName + " values"
		}
		c.Flags().String(flagName, "", desc)
	}
	cliapp.AddEnabledFlag(c)
	if strings.HasPrefix(use, "toggle") {
		cliapp.MarkFlagsRequired(c, "enabled")
	}

	c.RunE = func(cmd *cobra.Command, args []string) error {
		if err := app.RequireAuth(); err != nil {
			return err
		}
		if strings.HasPrefix(use, "toggle") {
			if err := cliapp.RequireFlags(cmd, "enabled"); err != nil {
				return err
			}
		}
		data, _ := cmd.Flags().GetString("data")
		body, err := cliapp.MergeDataWithFlags(data, cmd, fieldMap)
		if err != nil {
			return err
		}
		// Apply defaults for missing keys.
		for k, v := range defaults {
			if _, exists := body[k]; !exists {
				body[k] = v
			}
		}
		// Parse addrFields: --flag "v1,v2" → {"custom":["v1","v2"],"object":[]}
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
		if withID && fullUpdatePath != "" {
			body, err = fullUpdateBody(app, fullUpdatePath, args[0], body)
			if err != nil {
				return err
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

func fullUpdateBody(app *cliapp.Runtime, apiPath, id string, updates map[string]interface{}) (map[string]interface{}, error) {
	readClient := app.APIClient
	if app.APIClient.DryRun {
		readClient = app.NewClient(app.Session.BaseURL, app.Session.Token)
	}
	raw, err := readClient.Get(cliapp.APIBase+"/"+apiPath+"/"+id, nil)
	if err != nil {
		return nil, err
	}
	current, err := inputBodyFromGet(raw)
	if err != nil {
		return nil, err
	}
	for k, v := range updates {
		current[k] = v
	}
	return current, nil
}

func inputBodyFromGet(raw json.RawMessage) (map[string]interface{}, error) {
	var value interface{}
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, err
	}
	body, ok := findRuleObject(value)
	if !ok {
		return nil, &cliapp.ValidationError{Message: "unexpected routing get response"}
	}
	result := map[string]interface{}{}
	for k, v := range body {
		if k == "id" || strings.HasSuffix(k, "_int") {
			continue
		}
		result[k] = v
	}
	normalizeInputBody(result)
	return result, nil
}

func findRuleObject(value interface{}) (map[string]interface{}, bool) {
	switch typed := value.(type) {
	case map[string]interface{}:
		if _, ok := typed["tagname"]; ok {
			return typed, true
		}
		for _, key := range []string{"data", "results"} {
			if nested, ok := typed[key]; ok {
				if found, ok := findRuleObject(nested); ok {
					return found, true
				}
			}
		}
	case []interface{}:
		if len(typed) == 0 {
			return nil, false
		}
		return findRuleObject(typed[0])
	}
	return nil, false
}

func normalizeInputBody(value interface{}) {
	switch typed := value.(type) {
	case map[string]interface{}:
		for k, child := range typed {
			if (k == "custom" || k == "object") && emptyMap(child) {
				typed[k] = []interface{}{}
				continue
			}
			normalizeInputBody(child)
		}
	case []interface{}:
		for _, child := range typed {
			normalizeInputBody(child)
		}
	}
}

func emptyMap(value interface{}) bool {
	if typed, ok := value.(map[string]interface{}); ok {
		return len(typed) == 0
	}
	return false
}
