package qos

import (
	"encoding/json"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

// qosIPFieldMap maps CLI flags to API fields for QoS IP rules.
var qosIPFieldMap = map[string]string{
	"name":      "tagname",
	"interface": "interface",
	"protocol":  "protocol",
	"upload":    "upload",
	"download":  "download",
	"comment":   "comment",
	"enabled":   "enabled",
}
var qosIPAddrFields = map[string]string{
	"ip-addr":  "ip_addr",
	"src-port": "src_port",
	"dst-port": "dst_port",
}
var qosIPCreateDefaults = map[string]interface{}{
	"enabled":  "yes",
	"comment":  "",
	"type":     0,
	"protocol": "any",
	"ip_addr":  map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
	"src_port": map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
	"dst_port": map[string]interface{}{"custom": []interface{}{}, "object": []interface{}{}},
	"time": map[string]interface{}{
		"custom": []interface{}{
			map[string]interface{}{
				"type": "weekly", "weekdays": "1234567",
				"start_time": "00:00", "end_time": "23:59", "comment": "",
			},
		},
		"object": []interface{}{},
	},
}

// qosMACFieldMap maps CLI flags to API fields for QoS MAC rules.
var qosMACFieldMap = map[string]string{
	"name":      "tagname",
	"interface": "interface",
	"upload":    "upload",
	"download":  "download",
	"comment":   "comment",
	"enabled":   "enabled",
}
var qosMACAddrFields = map[string]string{
	"mac-addr": "mac_addr",
}
var qosMACCreateDefaults = map[string]interface{}{
	"enabled": "yes",
	"comment": "",
	"type":    0,
	"ip_type": "4",
	"time": map[string]interface{}{
		"custom": []interface{}{
			map[string]interface{}{
				"type": "weekly", "weekdays": "1234567",
				"start_time": "00:00", "end_time": "23:59", "comment": "",
			},
		},
		"object": []interface{}{},
	},
}

func addQoSIPFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Rule name (tagname)")
	cmd.Flags().String("ip-addr", "", "IP address or range")
	cmd.Flags().String("interface", "", "Network interface")
	cmd.Flags().String("protocol", "", "Protocol (any/tcp/udp/tcp+udp)")
	cmd.Flags().String("src-port", "", "Source port")
	cmd.Flags().String("dst-port", "", "Destination port")
	cmd.Flags().String("upload", "", "Upload bandwidth limit in Kbps")
	cmd.Flags().String("download", "", "Download bandwidth limit in Kbps")
	cmd.Flags().String("comment", "", "Comment")
	cliapp.AddEnabledFlag(cmd)
}

func addQoSMACFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Rule name (tagname)")
	cmd.Flags().String("mac-addr", "", "MAC address")
	cmd.Flags().String("interface", "", "Network interface")
	cmd.Flags().String("upload", "", "Upload bandwidth limit in Kbps")
	cmd.Flags().String("download", "", "Download bandwidth limit in Kbps")
	cmd.Flags().String("comment", "", "Comment")
	cliapp.AddEnabledFlag(cmd)
}

func New(app *cliapp.Runtime) *cobra.Command {
	qosCmd := &cobra.Command{
		Use:   "qos",
		Short: "QoS bandwidth control",
	}

	qosCmd.AddCommand(qosGroup(app, "ip", "network/qos/ip", addQoSIPFlags, qosIPFieldMap, qosIPAddrFields, qosIPCreateDefaults,
		[]string{"id", "tagname", "ip_addr", "interface", "protocol", "upload", "download", "enabled"},
		[]string{"name", "interface", "upload", "download"}))
	qosCmd.AddCommand(qosGroup(app, "mac", "network/qos/mac", addQoSMACFlags, qosMACFieldMap, qosMACAddrFields, qosMACCreateDefaults,
		[]string{"id", "tagname", "mac_addr", "interface", "upload", "download", "enabled"},
		[]string{"name", "interface", "upload", "download"}))
	return qosCmd
}

func qosGroup(app *cliapp.Runtime, name, apiPath string, addFlags func(*cobra.Command), fieldMap map[string]string, addrFields map[string]string, defaults map[string]interface{}, defaultColumns []string, requiredCreateFlags []string) *cobra.Command {
	group := &cobra.Command{Use: name, Short: "QoS rules based on " + name}

	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List " + name + " QoS rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if len(defaultColumns) > 0 {
				app.DefaultColumns = defaultColumns
			}
			page, _ := cmd.Flags().GetInt("page")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			raw, err := app.APIClient.Get(cliapp.APIBase+"/"+apiPath,
				cliapp.ListParamsWithPageSizeKey(page, pageSize, "", "", "", "limit"))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddPaginationFlags(listCmd)

	getCmd := &cobra.Command{
		Use:   "get ID",
		Short: "Get a single " + name + " QoS rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/"+apiPath+"/"+args[0], nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	createCmd := dataCommandImpl(app, "create", "Create a "+name+" QoS rule", false, addFlags, fieldMap, addrFields, defaults, func(body interface{}, id string) (json.RawMessage, error) {
		return app.APIClient.Post(cliapp.APIBase+"/"+apiPath, body)
	}, "")
	if len(requiredCreateFlags) > 0 {
		cliapp.MarkFlagsRequired(createCmd, requiredCreateFlags...)
		origRunE := createCmd.RunE
		createCmd.RunE = func(cmd *cobra.Command, args []string) error {
			if err := cliapp.RequireFlags(cmd, requiredCreateFlags...); err != nil {
				return err
			}
			return origRunE(cmd, args)
		}
	}
	updateCmd := dataCommandImpl(app, "update ID", "Update a "+name+" QoS rule", true, addFlags, fieldMap, addrFields, nil, func(body interface{}, id string) (json.RawMessage, error) {
		return app.APIClient.Put(cliapp.APIBase+"/"+apiPath+"/"+id, body)
	}, apiPath)
	toggleFieldMap := map[string]string{"enabled": "enabled"}
	toggleCmd := dataCommandImpl(app, "toggle ID", "Enable/disable a "+name+" QoS rule", true, cliapp.AddEnabledFlag, toggleFieldMap, nil, nil, func(body interface{}, id string) (json.RawMessage, error) {
		return app.APIClient.Patch(cliapp.APIBase+"/"+apiPath+"/"+id, body)
	}, "")

	deleteCmd := &cobra.Command{
		Use:     "delete ID",
		Aliases: []string{"rm"},
		Short:   "Delete a " + name + " QoS rule",
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
			raw, err := app.APIClient.Delete(cliapp.APIBase + "/" + apiPath + "/" + args[0])
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	deleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	group.AddCommand(listCmd, getCmd, createCmd, updateCmd, toggleCmd, deleteCmd)
	return group
}

type callWithBody func(body interface{}, id string) (json.RawMessage, error)

func dataCommandImpl(app *cliapp.Runtime, use, short string, withID bool, addFlags func(*cobra.Command), fieldMap map[string]string, addrFields map[string]string, defaults map[string]interface{}, fn callWithBody, fullUpdatePath string) *cobra.Command {
	c := &cobra.Command{
		Use:   use,
		Short: short,
	}
	if use == "create" {
		c.Aliases = []string{"new"}
	}
	if withID {
		c.Args = cobra.ExactArgs(1)
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
		if withID && fullUpdatePath != "" {
			body, err = qosFullUpdateBody(app, fullUpdatePath, args[0], body)
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
	c.Flags().String("data", "{}", "JSON body (escape hatch)")
	if addFlags != nil {
		addFlags(c)
	}
	if strings.HasPrefix(use, "toggle") {
		cliapp.MarkFlagsRequired(c, "enabled")
	}
	return c
}

func qosFullUpdateBody(app *cliapp.Runtime, apiPath, id string, updates map[string]interface{}) (map[string]interface{}, error) {
	readClient := app.APIClient
	if app.APIClient.DryRun {
		readClient = app.NewClient(app.Session.BaseURL, app.Session.Token)
	}
	raw, err := readClient.Get(cliapp.APIBase+"/"+apiPath+"/"+id, nil)
	if err != nil {
		return nil, err
	}
	current, err := qosInputFromGet(raw)
	if err != nil {
		return nil, err
	}
	for k, v := range updates {
		current[k] = v
	}
	return current, nil
}

func qosInputFromGet(raw json.RawMessage) (map[string]interface{}, error) {
	var value interface{}
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, err
	}
	body, ok := findQoSObject(value)
	if !ok {
		return nil, &cliapp.ValidationError{Message: "unexpected QoS get response"}
	}
	result := map[string]interface{}{}
	for k, v := range body {
		if k == "id" {
			continue
		}
		result[k] = v
	}
	normalizeQoSBody(result)
	return result, nil
}

func findQoSObject(value interface{}) (map[string]interface{}, bool) {
	switch typed := value.(type) {
	case map[string]interface{}:
		if _, ok := typed["tagname"]; ok {
			return typed, true
		}
		for _, key := range []string{"data", "results"} {
			if nested, ok := typed[key]; ok {
				if found, ok := findQoSObject(nested); ok {
					return found, true
				}
			}
		}
	case []interface{}:
		if len(typed) == 0 {
			return nil, false
		}
		return findQoSObject(typed[0])
	}
	return nil, false
}

func normalizeQoSBody(value interface{}) {
	switch typed := value.(type) {
	case map[string]interface{}:
		for k, child := range typed {
			if (k == "custom" || k == "object") && emptyQoSMap(child) {
				typed[k] = []interface{}{}
				continue
			}
			normalizeQoSBody(child)
		}
	case []interface{}:
		for _, child := range typed {
			normalizeQoSBody(child)
		}
	}
}

func emptyQoSMap(value interface{}) bool {
	if typed, ok := value.(map[string]interface{}); ok {
		return len(typed) == 0
	}
	return false
}
