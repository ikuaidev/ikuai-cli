package objects

import (
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

func New(app *cliapp.Runtime) *cobra.Command {
	objectsCmd := &cobra.Command{
		Use:   "objects",
		Short: "Network objects",
	}

	fieldMap := map[string]string{
		"name":  "group_name",
		"value": "group_value",
	}

	objectsCmd.AddCommand(objectGroup(app, "ip", "ip-objects", fieldMap, "ip"))
	objectsCmd.AddCommand(objectGroup(app, "ip6", "ip6-objects", fieldMap, "ipv6"))
	objectsCmd.AddCommand(objectGroup(app, "mac", "mac-objects", fieldMap, "mac"))
	objectsCmd.AddCommand(objectGroup(app, "port", "port-objects", fieldMap, "port"))
	objectsCmd.AddCommand(objectGroup(app, "proto", "proto-objects", fieldMap, "proto"))
	objectsCmd.AddCommand(objectGroup(app, "domain", "domain-objects", fieldMap, "domain"))
	objectsCmd.AddCommand(objectGroup(app, "time", "time-objects", fieldMap, ""))
	return objectsCmd
}

// objectGroup builds a CRUD command group for network objects.
// valueKey is the JSON key used in group_value array items (e.g. "ip", "mac", "port").
// For time objects, valueKey is empty and --value is not available (use --data).
func objectGroup(app *cliapp.Runtime, name, apiPath string, fieldMap map[string]string, valueKey string) *cobra.Command {
	group := &cobra.Command{Use: name, Short: "Manage " + name + " objects"}

	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List " + name + " objects",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
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
		Short: "Get a " + name + " object",
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

	createCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a " + name + " object",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			required := []string{"name"}
			if valueKey != "" {
				required = append(required, "value")
			}
			if err := cliapp.RequireFlags(cmd, required...); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, map[string]string{"name": "group_name"})
			if err != nil {
				return err
			}
			// Convert --value to group_value JSON array.
			if f := cmd.Flags().Lookup("value"); f != nil && f.Changed && valueKey != "" {
				body["group_value"] = buildGroupValue(f.Value.String(), valueKey)
			}
			if valueKey == "" {
				if !hasGroupValue(body) {
					if err := cliapp.RequireFlags(cmd, requiredTimeFlags(cmd)...); err != nil {
						return err
					}
				}
				if timeFlagsChanged(cmd) {
					body["group_value"] = buildTimeGroupValue(cmd)
				}
			}
			raw, err := app.APIClient.Post(cliapp.APIBase+"/"+apiPath, body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	createCmd.Flags().String("data", "{}", "JSON body")
	createCmd.Flags().String("name", "", "Object group name")
	if valueKey != "" {
		createCmd.Flags().String("value", "", "Comma-separated values (e.g. 1.2.3.4,5.6.7.8/24)")
		cliapp.MarkFlagsRequired(createCmd, "name", "value")
	} else {
		addTimeFlags(createCmd)
		cliapp.MarkFlagsRequired(createCmd, "name")
	}

	updateCmd := &cobra.Command{
		Use:   "update ID",
		Short: "Update a " + name + " object",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, map[string]string{"name": "group_name"})
			if err != nil {
				return err
			}
			if f := cmd.Flags().Lookup("value"); f != nil && f.Changed && valueKey != "" {
				body["group_value"] = buildGroupValue(f.Value.String(), valueKey)
			}
			if valueKey == "" && timeFlagsChanged(cmd) {
				body["group_value"] = buildTimeGroupValue(cmd)
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/"+apiPath+"/"+args[0], body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	updateCmd.Flags().String("data", "{}", "JSON body")
	updateCmd.Flags().String("name", "", "Object group name")
	if valueKey != "" {
		updateCmd.Flags().String("value", "", "Comma-separated values")
	} else {
		addTimeFlags(updateCmd)
	}

	deleteCmd := &cobra.Command{
		Use:     "delete ID",
		Aliases: []string{"rm"},
		Short:   "Delete a " + name + " object",
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

	refsCmd := &cobra.Command{
		Use:   "refs",
		Short: "Show rules that reference a " + name + " object",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			groupName, _ := cmd.Flags().GetString("group-name")
			raw, err := app.APIClient.Get(cliapp.APIBase+"/"+apiPath+"/ref", map[string]string{
				"group_name": groupName,
			})
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	refsCmd.Flags().String("group-name", "", "Object group name to query references for (required)")
	_ = refsCmd.MarkFlagRequired("group-name")

	group.AddCommand(listCmd, getCmd, createCmd, updateCmd, deleteCmd, refsCmd)
	return group
}

// buildGroupValue converts a comma-separated string into the JSON array
// structure expected by the iKuai API: [{"key":"val1"},{"key":"val2"}].
func buildGroupValue(csv, key string) []interface{} {
	parts := strings.Split(csv, ",")
	result := make([]interface{}, 0, len(parts))
	for _, v := range parts {
		v = strings.TrimSpace(v)
		if v != "" {
			result = append(result, map[string]interface{}{key: v})
		}
	}
	return result
}

func addTimeFlags(cmd *cobra.Command) {
	cmd.Flags().String("type", "", "Time object type (weekly/date)")
	cmd.Flags().String("weekdays", "", "Weekly days, e.g. 12345 or 1234567")
	cmd.Flags().String("start-time", "", "Start time, e.g. 00:00 or 2026-05-01T08:00")
	cmd.Flags().String("end-time", "", "End time, e.g. 20:00 or 2026-05-10T08:00")
	cmd.Flags().String("comment", "", "Object value comment")
}

func requiredTimeFlags(cmd *cobra.Command) []string {
	required := []string{"type", "start-time", "end-time"}
	if typ, _ := cmd.Flags().GetString("type"); typ == "weekly" {
		required = append(required, "weekdays")
	}
	return required
}

func timeFlagsChanged(cmd *cobra.Command) bool {
	for _, name := range []string{"type", "weekdays", "start-time", "end-time", "comment"} {
		if f := cmd.Flags().Lookup(name); f != nil && f.Changed {
			return true
		}
	}
	return false
}

func hasGroupValue(body map[string]interface{}) bool {
	_, ok := body["group_value"]
	return ok
}

func buildTimeGroupValue(cmd *cobra.Command) []interface{} {
	value := map[string]interface{}{}
	for flag, field := range map[string]string{
		"type":       "type",
		"weekdays":   "weekdays",
		"start-time": "start_time",
		"end-time":   "end_time",
		"comment":    "comment",
	} {
		if f := cmd.Flags().Lookup(flag); f != nil && f.Changed {
			v := strings.TrimSpace(f.Value.String())
			if v != "" {
				value[field] = v
			}
		}
	}
	return []interface{}{value}
}
