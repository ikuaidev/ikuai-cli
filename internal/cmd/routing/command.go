package routing

import (
	"encoding/json"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

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
		"comment":  "",
		"prio":     31,
		"src_addr": map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"time": map[string]interface{}{
			"custom": []interface{}{
				map[string]interface{}{"type": "weekly", "weekdays": "1234567", "start_time": "00:00", "end_time": "23:59", "comment": ""},
			},
			"object": []interface{}{},
		},
	}

	fiveTupleFieldMap = map[string]string{
		"name":      "tagname",
		"interface": "interface",
		"protocol":  "protocol",
		"priority":  "prio",
		"comment":   "comment",
		"enabled":   "enabled",
	}
	fiveTupleAddrFields = map[string]string{
		"src-addr": "src_addr",
		"dst-addr": "dst_addr",
		"src-port": "src_port",
		"dst-port": "dst_port",
	}
	fiveTupleDefaults = map[string]interface{}{
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
		"name":      "tagname",
		"interface": "interface",
		"priority":  "prio",
		"comment":   "comment",
		"enabled":   "enabled",
	}
	l7AddrFields = map[string]string{
		"app-proto": "app_proto",
		"src-addr":  "src_addr",
	}
	l7Defaults = map[string]interface{}{
		"comment":    "",
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
		"weight":    "weight",
		"isp-name":  "isp_name",
		"comment":   "comment",
		"enabled":   "enabled",
	}
	loadBalanceDefaults = map[string]interface{}{
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
			fieldMap:       domainFieldMap,
			addrFields:     domainAddrFields,
			createDefaults: domainDefaults,
			defaultColumns: []string{"id", "tagname", "domain", "interface", "src_addr", "prio", "enabled"},
		}),
		ruleGroup(app, "five-tuple", "5-tuple stream rules (src/dst IP, port, protocol)", "routing/five-tuple-rules", ruleGroupOpts{
			fieldMap:       fiveTupleFieldMap,
			addrFields:     fiveTupleAddrFields,
			createDefaults: fiveTupleDefaults,
			defaultColumns: []string{"id", "tagname", "src_addr", "dst_addr", "protocol", "dst_port", "interface", "prio", "enabled"},
		}),
		ruleGroup(app, "l7", "L7 application protocol stream rules", "routing/app-protocols", ruleGroupOpts{
			fieldMap:       l7FieldMap,
			addrFields:     l7AddrFields,
			createDefaults: l7Defaults,
			defaultColumns: []string{"id", "tagname", "app_proto", "interface", "src_addr", "prio", "enabled"},
		}),
		ruleGroup(app, "load-balance", "Load balance rules", "routing/load-balance-rules", ruleGroupOpts{
			fieldMap:       loadBalanceFieldMap,
			createDefaults: loadBalanceDefaults,
			defaultColumns: []string{"id", "tagname", "interface", "mode", "weight", "isp_name", "enabled"},
		}),
		ruleGroup(app, "updown", "Upstream/downstream rules", "routing/updown", ruleGroupOpts{
			fieldMap:       updownFieldMap,
			addrFields:     updownAddrFields,
			createDefaults: updownDefaults,
			defaultColumns: []string{"id", "tagname", "upiface", "downiface", "protocol", "src_addr", "enabled"},
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
		"netmask":   "netmask",
		"priority":  "prio",
		"comment":   "comment",
		"enabled":   "enabled",
	}
	staticDefaults := map[string]interface{}{
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
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(listCmd)

	createCmd := writeCmd(app, "create", "Create a static route", false, staticFieldMap, nil, staticDefaults,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Post(cliapp.APIBase+"/routing/static-routes", body)
		})
	updateCmd := writeCmd(app, "update ID", "Update a static route", true, staticFieldMap, nil, nil,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Put(cliapp.APIBase+"/routing/static-routes/"+id, body)
		})

	toggleFieldMap := map[string]string{"enabled": "enabled"}
	group.AddCommand(
		listCmd,
		createCmd,
		updateCmd,
		writeCmd(app, "toggle ID", "Enable/disable a static route", true, toggleFieldMap, nil, nil,
			func(body interface{}, id string) (json.RawMessage, error) {
				return app.APIClient.Patch(cliapp.APIBase+"/routing/static-routes/"+id, body)
			}),
		deleteByIDCmd(app, "delete ID", "Delete a static route", "/routing/static-routes/"),
	)
	return group
}

type ruleGroupOpts struct {
	fieldMap       map[string]string
	addrFields     map[string]string
	createDefaults map[string]interface{}
	defaultColumns []string
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
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(listCmd)

	createCmd := writeCmd(app, "create", "Create a "+use+" rule", false, opts.fieldMap, opts.addrFields, opts.createDefaults,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Post(cliapp.APIBase+"/"+apiPath, body)
		})
	updateCmd := writeCmd(app, "update ID", "Update a "+use+" rule", true, opts.fieldMap, opts.addrFields, nil,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Put(cliapp.APIBase+"/"+apiPath+"/"+id, body)
		})
	toggleFieldMap := map[string]string{"enabled": "enabled"}
	toggleCmd := writeCmd(app, "toggle ID", "Enable/disable a "+use+" rule", true, toggleFieldMap, nil, nil,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Patch(cliapp.APIBase+"/"+apiPath+"/"+id, body)
		})

	grp.AddCommand(
		listCmd,
		createCmd,
		updateCmd,
		toggleCmd,
		deleteByIDCmd(app, "delete ID", "Delete a "+use+" rule", "/"+apiPath+"/"),
	)
	return grp
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
func writeCmd(app *cliapp.Runtime, use, short string, withID bool, fieldMap map[string]string, addrFields map[string]string, defaults map[string]interface{}, fn callWithBody) *cobra.Command {
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
		c.Flags().String(flagName, "", flagName+" value")
	}
	// Register addrFields flags.
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
