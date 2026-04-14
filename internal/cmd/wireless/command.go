package wireless

import (
	"encoding/json"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

// --- Field maps, addr fields, and defaults for wireless subcommands ---

var (
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
		requiredCreateFlags: []string{"name", "ap"},
	}))
	wirelessCmd.AddCommand(ruleGroup(app, "vlan", "Wireless VLAN rules", "wireless/vlan/rules", ruleGroupOpts{
		fieldMap:            vlanFieldMap,
		addrFields:          vlanAddrFields,
		createDefaults:      vlanDefaults,
		defaultColumns:      []string{"id", "tagname", "vlanid", "lmac", "lssid", "enabled"},
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

	createCmd := writeCmd(app, "create", "Create a "+short, false, opts.fieldMap, opts.addrFields, opts.createDefaults,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Post(cliapp.APIBase+"/"+apiPath, body)
		})
	if len(opts.requiredCreateFlags) > 0 {
		origRunE := createCmd.RunE
		createCmd.RunE = func(cmd *cobra.Command, args []string) error {
			if err := cliapp.RequireFlags(cmd, opts.requiredCreateFlags...); err != nil {
				return err
			}
			return origRunE(cmd, args)
		}
	}
	updateCmd := writeCmd(app, "update ID", "Update a "+short, true, opts.fieldMap, opts.addrFields, nil,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Put(cliapp.APIBase+"/"+apiPath+"/"+id, body)
		})
	toggleFieldMap := map[string]string{"enabled": "enabled"}
	toggleCmd := writeCmd(app, "toggle ID", "Enable/disable a "+short, true, toggleFieldMap, nil, nil,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Patch(cliapp.APIBase+"/"+apiPath+"/"+id, body)
		})

	group.AddCommand(
		listCmd,
		createCmd,
		updateCmd,
		toggleCmd,
		deleteByIDCmd(app, "delete ID", "Delete a "+short, "/"+apiPath+"/"),
	)
	return group
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
		configCmdWithFlags(app, "set", "Update AC service configuration", func(c *cobra.Command) {
			c.Flags().String("ac-status", "", "AC service status (0=off, 1=on)")
		}, map[string]string{"ac-status": "ac_status"}, func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Put(cliapp.APIBase+"/network/ac/services", body)
		}),
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
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/network/ac/ap-config",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(c)
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
		return app.APIClient.Put(cliapp.APIBase+"/network/ac/ap-config/"+id, body)
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
