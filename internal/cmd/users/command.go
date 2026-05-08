package users

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

// flagDescs provides human-readable descriptions for CLI flags.
var flagDescs = map[string]string{
	"username":    "Login username",
	"password":    "Login password",
	"ppptype":     "PPP type (any/pppoe/pptp/l2tp/ovpn/web/pppoe_relay/ike)",
	"packages":    "Package ID (0=custom)",
	"upload":      "Upload bandwidth limit (KB/s, 0=unlimited)",
	"download":    "Download bandwidth limit (KB/s, 0=unlimited)",
	"start-time":  "Start Unix timestamp",
	"expires":     "Expire Unix timestamp (0=never)",
	"share":       "Max concurrent sessions",
	"ip-type":     "IP type (0=fixed IP, 1=address pool)",
	"auto-mac":    "Auto-bind MAC (0=no, 1=yes)",
	"auto-vlanid": "Auto-bind VLAN (0=no, 1=yes)",
	"bind-vlanid": "Bound VLAN ID",
	"pppname":     "PPP relay WAN interface",
	"pppoev6-wan": "PPPoE IPv6 WAN prefix source",
	"bind-ifname": "Bound interface name",
	"mac":         "Bound MAC address",
	"address":     "Address",
	"real-name":   "Real name",
	"phone":       "Phone number",
	"cardid":      "ID card number",
	"comment":     "Comment",
	"name":        "Package name",
	"time":        "Duration (hours)",
	"price":       "Price",
	"up-speed":    "Upload speed limit (KB/s)",
	"down-speed":  "Download speed limit (KB/s)",
}

var (
	onlineDefaultColumns  = []string{"id", "username", "ppptype", "ip_addr", "mac", "auth_time", "session", "interface", "expires", "packages"}
	accountDefaultColumns = []string{"id", "username", "ppptype", "upload", "download", "share", "expires", "enabled"}
	packageDefaultColumns = []string{"id", "packname", "packtime", "price", "up_speed", "down_speed", "comment"}

	accountFieldMap = map[string]string{
		"username":    "username",
		"password":    "passwd",
		"enabled":     "enabled",
		"ppptype":     "ppptype",
		"packages":    "packages",
		"upload":      "upload",
		"download":    "download",
		"start-time":  "start_time",
		"expires":     "expires",
		"share":       "share",
		"ip-type":     "ip_type",
		"auto-mac":    "auto_mac",
		"auto-vlanid": "auto_vlanid",
		"bind-vlanid": "bind_vlanid",
		"pppname":     "pppname",
		"pppoev6-wan": "pppoev6_wan",
		"bind-ifname": "bind_ifname",
		"mac":         "mac",
		"address":     "address",
		"real-name":   "name",
		"phone":       "phone",
		"cardid":      "cardid",
		"comment":     "comment",
	}
	accountDefaults = map[string]interface{}{
		"enabled":     "yes",
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
	accountUpdateInputFields = []string{
		"username", "passwd", "enabled", "ppptype", "packages", "upload", "download",
		"start_time", "expires", "share", "ip_type", "auto_mac", "auto_vlanid",
		"bind_vlanid", "pppname", "pppoev6_wan", "bind_ifname", "src_addr", "mac",
		"address", "name", "phone", "cardid", "comment",
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
	packageUpdateInputFields = []string{
		"packname", "packtime", "price", "up_speed", "down_speed", "comment",
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
			app.DefaultColumns = onlineDefaultColumns
			page, pageSize, _, _, _ := cliapp.GetListParams(cmd)
			params := onlineListParams(cmd, page, pageSize)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/auth/online-users", params)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	addOnlineListFlags(onlineCmd)
	onlineCmd.AddCommand(getByIDCmd(app, "get ID", "Get a single online session", "/auth/online-users/", onlineDefaultColumns))
	usersCmd.AddCommand(onlineCmd)

	kickCmd := &cobra.Command{
		Use:   "kick ID",
		Short: "Kick an online user offline",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			yes, _ := cmd.Flags().GetBool("yes")
			if err := cliapp.ConfirmDelete(app.Stdout, app.Stderr, "online user", args[0], yes); err != nil {
				return err
			}
			raw, err := app.APIClient.Delete(cliapp.APIBase + "/auth/online-users/" + args[0])
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	kickCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
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
			app.DefaultColumns = accountDefaultColumns
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/auth/users",
				cliapp.ListParamsWithPageSizeKey(page, pageSize, filter, order, orderBy, "limit"))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(listCmd)

	createCmd := writeCmd(app, "create", "Create a user account", false, accountFieldMap, nil, accountDefaults,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Post(cliapp.APIBase+"/auth/users", body)
		})
	cliapp.MarkFlagsRequired(createCmd, "username", "password")
	{
		origRunE := createCmd.RunE
		createCmd.RunE = func(cmd *cobra.Command, args []string) error {
			if err := cliapp.RequireFlags(cmd, "username", "password"); err != nil {
				return err
			}
			return origRunE(cmd, args)
		}
	}
	group.AddCommand(
		listCmd,
		getByIDCmd(app, "get ID", "Get a single user account", "/auth/users/", accountDefaultColumns),
		createCmd,
		updateByIDCmd(app, "update ID", "Update a user account", accountFieldMap, accountUpdateInputFields, cliapp.APIBase+"/auth/users/"),
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
			app.DefaultColumns = packageDefaultColumns
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/auth/packages",
				cliapp.ListParamsWithPageSizeKey(page, pageSize, filter, order, orderBy, "limit"))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(listCmd)

	pkgCreateCmd := writeCmd(app, "create", "Create a package", false, packageFieldMap, nil, packageDefaults,
		func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Post(cliapp.APIBase+"/auth/packages", body)
		})
	cliapp.MarkFlagsRequired(pkgCreateCmd, "name", "time", "price", "up-speed", "down-speed")
	{
		origRunE := pkgCreateCmd.RunE
		pkgCreateCmd.RunE = func(cmd *cobra.Command, args []string) error {
			if err := cliapp.RequireFlags(cmd, "name", "time", "price", "up-speed", "down-speed"); err != nil {
				return err
			}
			return origRunE(cmd, args)
		}
	}
	group.AddCommand(
		listCmd,
		getByIDCmd(app, "get ID", "Get a single package", "/auth/packages/", packageDefaultColumns),
		pkgCreateCmd,
		updateByIDCmd(app, "update ID", "Update a package", packageFieldMap, packageUpdateInputFields, cliapp.APIBase+"/auth/packages/"),
		deleteByIDCmd(app, "delete ID", "Delete a package", "/auth/packages/"),
	)
	return group
}

func updateByIDCmd(app *cliapp.Runtime, use, short string, fieldMap map[string]string, inputFields []string, apiPathPrefix string) *cobra.Command {
	c := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
	}
	addBodyFlags(c, fieldMap)
	c.RunE = func(cmd *cobra.Command, args []string) error {
		if err := app.RequireAuth(); err != nil {
			return err
		}
		data, _ := cmd.Flags().GetString("data")
		body, err := buildFullBody(app, cmd, data, fieldMap, inputFields, apiPathPrefix+args[0])
		if err != nil {
			return err
		}
		raw, err := app.APIClient.Put(apiPathPrefix+args[0], body)
		if err != nil {
			return err
		}
		app.PrintRaw(raw)
		return nil
	}
	return c
}

func addOnlineListFlags(cmd *cobra.Command) {
	cliapp.AddPaginationFlags(cmd)
	cmd.Flags().String("keywords", "", "Search keyword")
	cmd.Flags().String("finds", "", "Search fields, comma-separated")
	cmd.Flags().String("order", "", "Sort direction: asc|desc")
	cmd.Flags().String("order-by", "", "Sort field")
	_ = cmd.RegisterFlagCompletionFunc("order", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"asc", "desc"}, cobra.ShellCompDirectiveNoFileComp
	})
}

func onlineListParams(cmd *cobra.Command, page, pageSize int) map[string]string {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	params := map[string]string{
		"page":  strconv.Itoa(page),
		"limit": strconv.Itoa(pageSize),
	}
	if keywords, _ := cmd.Flags().GetString("keywords"); keywords != "" {
		params["KEYWORDS"] = keywords
	}
	if finds, _ := cmd.Flags().GetString("finds"); finds != "" {
		params["FINDS"] = finds
	}
	if order, _ := cmd.Flags().GetString("order"); order != "" {
		params["ORDER"] = order
	}
	if orderBy, _ := cmd.Flags().GetString("order-by"); orderBy != "" {
		params["ORDER_BY"] = orderBy
	}
	return params
}

func getByIDCmd(app *cliapp.Runtime, use, short, apiPath string, defaultCols ...[]string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if len(defaultCols) > 0 {
				app.DefaultColumns = defaultCols[0]
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

func addBodyFlags(c *cobra.Command, fieldMap map[string]string) {
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
	if _, ok := fieldMap["enabled"]; ok {
		cliapp.AddEnabledFlag(c)
	}
}

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
	if _, ok := fieldMap["enabled"]; ok {
		cliapp.AddEnabledFlag(c)
	}
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

func buildFullBody(app *cliapp.Runtime, cmd *cobra.Command, data string, fieldMap map[string]string, inputFields []string, getPath string) (map[string]interface{}, error) {
	changes, err := cliapp.MergeDataWithFlags(data, cmd, fieldMap)
	if err != nil {
		return nil, err
	}
	readClient := app.APIClient
	if app.APIClient.DryRun {
		readClient = app.NewClient(app.Session.BaseURL, app.Session.Token)
	}
	raw, err := readClient.Get(getPath, nil)
	if err != nil {
		if hasAllInputFields(changes, inputFields) {
			return changes, nil
		}
		return nil, err
	}
	current, err := extractUsersInputObject(raw, inputFields)
	if err != nil {
		if hasAllInputFields(changes, inputFields) {
			return changes, nil
		}
		return nil, err
	}
	for k, v := range changes {
		current[k] = v
	}
	return current, nil
}

func extractUsersInputObject(raw json.RawMessage, inputFields []string) (map[string]interface{}, error) {
	var v interface{}
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, err
	}
	obj, err := findFirstUsersObject(v)
	if err != nil {
		return nil, err
	}
	out := map[string]interface{}{}
	for _, key := range inputFields {
		if val, ok := obj[key]; ok {
			out[key] = val
		}
	}
	return out, nil
}

func findFirstUsersObject(v interface{}) (map[string]interface{}, error) {
	switch data := v.(type) {
	case map[string]interface{}:
		if rows, ok := data["data"].([]interface{}); ok {
			return findFirstUsersObject(rows)
		}
		if rows, ok := data["results"].([]interface{}); ok {
			return findFirstUsersObject(rows)
		}
		return data, nil
	case []interface{}:
		if len(data) == 0 {
			return nil, &cliapp.ValidationError{Message: "empty users response"}
		}
		return findFirstUsersObject(data[0])
	default:
		return nil, &cliapp.ValidationError{Message: "unexpected users response"}
	}
}

func hasAllInputFields(body map[string]interface{}, inputFields []string) bool {
	for _, key := range inputFields {
		if _, ok := body[key]; !ok {
			return false
		}
	}
	return true
}
