package users

import (
	"encoding/json"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

var (
	accountFieldMap = map[string]string{
		"username": "username",
		"password": "passwd",
		"ppptype":  "ppptype",
		"upload":   "upload",
		"download": "download",
		"share":    "share",
		"comment":  "comment",
		"enabled":  "enabled",
	}
	accountDefaults = map[string]interface{}{
		"comment":     "",
		"ppptype":     "any",
		"packages":    0,
		"upload":      0,
		"download":    0,
		"start_time":  0,
		"expires":     0,
		"share":       1,
		"ip_type":     0,
		"auto_mac":    0,
		"auto_vlanid": 0,
		"bind_vlanid": "0",
		"bind_ifname": "any",
		"src_addr":    map[string]interface{}{"custom": map[string]interface{}{}},
	}

	packageFieldMap = map[string]string{
		"name":       "packname",
		"time":       "packtime",
		"price":      "price",
		"up-speed":   "up_speed",
		"down-speed": "down_speed",
		"comment":    "comment",
	}
	packageDefaults = map[string]interface{}{
		"comment": "",
	}
)

func New(app *cliapp.Runtime) *cobra.Command {
	usersCmd := &cobra.Command{
		Use:   "users",
		Short: "User management",
	}

	usersCmd.AddCommand(accountsGroup(app))

	onlineCmd := &cobra.Command{
		Use:   "online",
		Short: "List online sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/auth/users",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(onlineCmd)
	usersCmd.AddCommand(onlineCmd)

	kickCmd := &cobra.Command{
		Use:   "kick [SESSION_ID]",
		Short: "Kick user offline",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			sid, _ := cmd.Flags().GetString("session-id")
			if len(args) > 0 && sid == "" {
				sid = args[0]
			}
			if sid == "" {
				return &cliapp.ValidationError{Message: "SESSION_ID is required: use --session-id or pass as argument"}
			}
			raw, err := app.APIClient.Delete(cliapp.APIBase + "/auth/users/" + sid)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	kickCmd.Flags().String("session-id", "", "Session ID to kick")
	usersCmd.AddCommand(kickCmd)

	usersCmd.AddCommand(packagesGroup(app))
	return usersCmd
}

func accountsGroup(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "accounts", Short: "Auth accounts"}

	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "username", "ppptype", "upload", "download", "share", "expires", "enabled"}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/auth/users",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(listCmd)

	group.AddCommand(
		listCmd,
		getByIDCmd(app, "get ID", "Get a single user account", "/auth/users/"),
		writeCmd(app, "create", "Create a user account", false, accountFieldMap, nil, accountDefaults,
			func(body interface{}, id string) (json.RawMessage, error) {
				return app.APIClient.Post(cliapp.APIBase+"/auth/users", body)
			}),
		writeCmd(app, "update ID", "Update a user account", true, accountFieldMap, nil, nil,
			func(body interface{}, id string) (json.RawMessage, error) {
				return app.APIClient.Put(cliapp.APIBase+"/auth/users/"+id, body)
			}),
		deleteByIDCmd(app, "delete ID", "Delete a user account", "/auth/users/"),
	)
	return group
}

func packagesGroup(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "packages", Short: "Auth packages"}

	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/auth/packages",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(listCmd)

	group.AddCommand(
		listCmd,
		getByIDCmd(app, "get ID", "Get a single package", "/auth/packages/"),
		writeCmd(app, "create", "Create a package", false, packageFieldMap, nil, packageDefaults,
			func(body interface{}, id string) (json.RawMessage, error) {
				return app.APIClient.Post(cliapp.APIBase+"/auth/packages", body)
			}),
		writeCmd(app, "update ID", "Update a package", true, packageFieldMap, nil, nil,
			func(body interface{}, id string) (json.RawMessage, error) {
				return app.APIClient.Put(cliapp.APIBase+"/auth/packages/"+id, body)
			}),
		deleteByIDCmd(app, "delete ID", "Delete a package", "/auth/packages/"),
	)
	return group
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
