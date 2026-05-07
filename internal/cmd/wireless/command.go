package wireless

import (
	"encoding/json"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

// flagDescs provides human-readable descriptions for CLI flags.
var flagDescs = map[string]string{
	"name":    "Rule name",
	"mode":    "Mode (0=blacklist, 1=whitelist)",
	"ssid":    "SSID filter (default: ALL)",
	"ap":      "AP filter (default: ALL)",
	"mac":     "MAC address (comma-separated)",
	"week":    "Active days (e.g. 1234567)",
	"time":    "Active time range (e.g. 00:00-23:59)",
	"vlan-id": "VLAN ID",
	"comment": "Comment",
}

// --- Field maps, addr fields, and defaults for wireless subcommands ---

var (
	acAPDefaultColumns = []string{"id", "tagname", "ip_addr", "mac", "status", "connected", "online", "version", "ap_model"}

	blacklistFieldMap = map[string]string{
		"name":    "tagname",
		"mode":    "mode",
		"ssid":    "lssid",
		"ap":      "lap",
		"week":    "week",
		"time":    "time",
		"comment": "comment",
		"enabled": "enabled",
	}
	blacklistAddrFields = map[string]string{
		"mac": "lmac",
	}
	blacklistDefaults = map[string]interface{}{
		"enabled": "yes",
		"mode":    0,
		"lmac":    map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"lssid":   "ALL",
		"lap":     "ALL",
		"week":    "1234567",
		"time":    "00:00-23:59",
		"comment": "",
	}

	vlanFieldMap = map[string]string{
		"name":    "tagname",
		"vlan-id": "vlanid",
		"ssid":    "lssid",
		"comment": "comment",
		"enabled": "enabled",
	}
	vlanAddrFields = map[string]string{
		"mac": "lmac",
	}
	vlanDefaults = map[string]interface{}{
		"enabled": "yes",
		"lssid":   "ALL",
		"lmac":    map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
		"comment": "",
	}
)

func New(app *cliapp.Runtime) *cobra.Command {
	wirelessCmd := &cobra.Command{
		Use:   "wireless",
		Short: "Wireless control",
	}

	wirelessCmd.AddCommand(ruleGroup(app, "blacklist", "Wireless access control (blacklist)", "wireless/access-control/rules", ruleGroupOpts{
		fieldMap:            blacklistFieldMap,
		addrFields:          blacklistAddrFields,
		createDefaults:      blacklistDefaults,
		defaultColumns:      []string{"id", "tagname", "mode", "lmac", "lssid", "lap", "enabled"},
		inputFields:         []string{"enabled", "tagname", "mode", "lmac", "lssid", "lap", "week", "time", "comment"},
		requiredCreateFlags: []string{"name"},
	}))
	wirelessCmd.AddCommand(ruleGroup(app, "vlan", "Wireless VLAN rules", "wireless/vlan/rules", ruleGroupOpts{
		fieldMap:            vlanFieldMap,
		addrFields:          vlanAddrFields,
		createDefaults:      vlanDefaults,
		defaultColumns:      []string{"id", "tagname", "vlanid", "lmac", "lssid", "enabled"},
		inputFields:         []string{"enabled", "tagname", "vlanid", "lmac", "lssid", "comment"},
		requiredCreateFlags: []string{"name", "vlan-id"},
	}))
	wirelessCmd.AddCommand(acGroup(app))
	return wirelessCmd
}

type ruleGroupOpts struct {
	fieldMap            map[string]string
	addrFields          map[string]string
	createDefaults      map[string]interface{}
	defaultColumns      []string
	inputFields         []string
	requiredCreateFlags []string
}

func ruleGroup(app *cliapp.Runtime, use, short, apiPath string, opts ruleGroupOpts) *cobra.Command {
	group := &cobra.Command{Use: use, Short: short}

	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List " + short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if len(opts.defaultColumns) > 0 {
				app.DefaultColumns = opts.defaultColumns
			}
			page, pageSize, _, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/"+apiPath,
				cliapp.ListParamsWithPageSizeKey(page, pageSize, "", order, orderBy, "limit"))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	addWirelessRuleListFlags(listCmd)

	createCmd := writeCmd(app, "create", "Create a "+short, false, opts.fieldMap, opts.addrFields, opts.createDefaults,
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
	getCmd := &cobra.Command{
		Use:   "get ID",
		Short: "Get a " + short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if len(opts.defaultColumns) > 0 {
				app.DefaultColumns = opts.defaultColumns
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/"+apiPath+"/"+args[0], nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	updateCmd := writeCmd(app, "update ID", "Update a "+short, true, opts.fieldMap, opts.addrFields, nil,
		func(body interface{}, id string) (json.RawMessage, error) {
			updates, _ := body.(map[string]interface{})
			fullBody, err := fullWirelessUpdateBody(app, apiPath, id, updates, opts.inputFields)
			if err != nil {
				return nil, err
			}
			return app.APIClient.Put(cliapp.APIBase+"/"+apiPath+"/"+id, fullBody)
		})
	toggleFieldMap := map[string]string{"enabled": "enabled"}
	toggleCmd := writeCmd(app, "toggle ID", "Enable/disable a "+short, true, toggleFieldMap, nil, nil,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Patch(cliapp.APIBase+"/"+apiPath+"/"+id, body)
		})

	group.AddCommand(
		listCmd,
		getCmd,
		createCmd,
		updateCmd,
		toggleCmd,
		deleteByIDCmd(app, "delete ID", "Delete a "+short, "/"+apiPath+"/"),
	)
	return group
}

func addWirelessRuleListFlags(cmd *cobra.Command) {
	cliapp.AddPaginationFlags(cmd)
	cmd.Flags().String("order", "", "Sort direction: asc|desc")
	cmd.Flags().String("order-by", "", "Sort field")
	_ = cmd.RegisterFlagCompletionFunc("order", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"asc", "desc"}, cobra.ShellCompDirectiveNoFileComp
	})
}

func acGroup(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "ac", Short: "AC management"}

	group.AddCommand(
		&cobra.Command{
			Use:   "get",
			Short: "Get AC config",
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := app.RequireAuth(); err != nil {
					return err
				}
				raw, err := app.APIClient.Get(cliapp.APIBase+"/network/ac/services", nil)
				if err != nil {
					return err
				}
				app.PrintRaw(raw)
				return nil
			},
		},
		&cobra.Command{
			Use:   "start",
			Short: "Start AC",
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := app.RequireAuth(); err != nil {
					return err
				}
				raw, err := app.APIClient.Post(cliapp.APIBase+"/network/ac/services:start", map[string]string{})
				if err != nil {
					return err
				}
				app.PrintRaw(raw)
				return nil
			},
		},
		&cobra.Command{
			Use:   "stop",
			Short: "Stop AC",
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := app.RequireAuth(); err != nil {
					return err
				}
				raw, err := app.APIClient.Post(cliapp.APIBase+"/network/ac/services:stop", map[string]string{})
				if err != nil {
					return err
				}
				app.PrintRaw(raw)
				return nil
			},
		},
		acAPListCmd(app),
		acAPGetCmd(app),
		acAPUpdateCmd(app),
	)

	return group
}

func acAPListCmd(app *cliapp.Runtime) *cobra.Command {
	c := &cobra.Command{
		Use:   "ap-list",
		Short: "List AP configs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = acAPDefaultColumns
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/ac/ap-config", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	return c
}

func acAPGetCmd(app *cliapp.Runtime) *cobra.Command {
	return &cobra.Command{
		Use:   "ap-get ID",
		Short: "Get AP config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = acAPDefaultColumns
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/ac/ap-config/"+args[0], nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
}

func acAPUpdateCmd(app *cliapp.Runtime) *cobra.Command {
	apFieldMap := map[string]string{
		"name":    "tagname",
		"ssid1":   "ssid1",
		"ssid3":   "ssid3",
		"enc1":    "enc1",
		"key1":    "key1",
		"enc3":    "enc3",
		"key3":    "key3",
		"channel": "channel",
		"comment": "comment",
	}
	return configCmdWithFlags(app, "ap-update ID", "Update an AP configuration", func(c *cobra.Command) {
		c.Flags().String("name", "", "AP name (tagname)")
		c.Flags().String("ssid1", "", "2.4G SSID1 name")
		c.Flags().String("ssid3", "", "5G SSID1 name")
		c.Flags().String("enc1", "", "2.4G SSID1 encryption (off/wpa/wpa2/wpa+wpa2)")
		c.Flags().String("key1", "", "2.4G SSID1 password")
		c.Flags().String("enc3", "", "5G SSID1 encryption")
		c.Flags().String("key3", "", "5G SSID1 password")
		c.Flags().String("channel", "", "2.4G channel (0=auto)")
		c.Flags().String("comment", "", "Comment")
	}, apFieldMap, func(body interface{}, id string) (json.RawMessage, error) {
		updates, _ := body.(map[string]interface{})
		fullBody, err := fullWirelessUpdateBody(app, "network/ac/ap-config", id, updates, nil)
		if err != nil {
			return nil, err
		}
		return app.APIClient.Put(cliapp.APIBase+"/network/ac/ap-config/"+id, fullBody)
	})
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

// writeCmd builds a create/update/toggle command with semantic flags + addrFields + defaults.
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
	if strings.Contains(use, "ID") {
		c.Args = cobra.ExactArgs(1)
	}
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
		id := ""
		if len(args) > 0 {
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

func fullWirelessUpdateBody(app *cliapp.Runtime, apiPath, id string, updates map[string]interface{}, inputFields []string) (map[string]interface{}, error) {
	readClient := app.APIClient
	if app.APIClient.DryRun {
		readClient = app.NewClient(app.Session.BaseURL, app.Session.Token)
	}
	raw, err := readClient.Get(cliapp.APIBase+"/"+apiPath+"/"+id, nil)
	if err != nil {
		return nil, err
	}
	current, err := wirelessInputBodyFromGet(raw, inputFields)
	if err != nil {
		return nil, err
	}
	for k, v := range updates {
		current[k] = v
	}
	return current, nil
}

func wirelessInputBodyFromGet(raw json.RawMessage, inputFields []string) (map[string]interface{}, error) {
	var value interface{}
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, err
	}
	body, ok := findWirelessObject(value)
	if !ok {
		return nil, &cliapp.ValidationError{Message: "unexpected wireless get response"}
	}
	if len(inputFields) == 0 {
		result := map[string]interface{}{}
		for k, v := range body {
			result[k] = v
		}
		return result, nil
	}
	result := map[string]interface{}{}
	for _, key := range inputFields {
		if value, ok := body[key]; ok {
			result[key] = value
		}
	}
	return result, nil
}

func findWirelessObject(value interface{}) (map[string]interface{}, bool) {
	switch typed := value.(type) {
	case map[string]interface{}:
		if _, ok := typed["tagname"]; ok {
			return typed, true
		}
		for _, key := range []string{"data", "results"} {
			if nested, ok := typed[key]; ok {
				if found, ok := findWirelessObject(nested); ok {
					return found, true
				}
			}
		}
	case []interface{}:
		if len(typed) == 0 {
			return nil, false
		}
		return findWirelessObject(typed[0])
	}
	return nil, false
}
