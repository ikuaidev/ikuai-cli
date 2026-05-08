package advanced

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

// ---------- field maps ----------

var ftpCreateFieldMap = map[string]string{
	"username":   "username",
	"password":   "passwd",
	"permission": "permission",
	"home-dir":   "home_dir",
	"upload":     "upload",
	"download":   "download",
	"enabled":    "enabled",
}

var ftpUpdateFieldMap = map[string]string{
	"name":       "tagname",
	"username":   "username",
	"password":   "passwd",
	"permission": "permission",
	"home-dir":   "home_dir",
	"upload":     "upload",
	"download":   "download",
	"enabled":    "enabled",
}

var httpFieldMap = map[string]string{
	"name":        "tagname",
	"port":        "http_port",
	"server-name": "server_name",
	"ssl":         "ssl_on",
	"autoindex":   "autoindex",
	"download":    "download",
	"home-dir":    "home_dir",
	"access":      "access",
	"enabled":     "enabled",
}

var sambaFieldMap = map[string]string{
	"name":       "name",
	"username":   "username",
	"password":   "passwd",
	"permission": "perm",
	"guest":      "guest",
	"home-dir":   "home_dir",
	"enabled":    "enabled",
}

// Config fieldMaps for service config-set commands.
var ftpConfigFieldMap = map[string]string{
	"open-ftp":   "open_ftp",
	"ftp-port":   "ftp_port",
	"ftp-access": "ftp_access",
}

var sambaConfigFieldMap = map[string]string{
	"enabled":   "enabled",
	"workgroup": "workgroup",
	"wsdd2":     "wsdd2",
	"access":    "access",
}

var snmpdFieldMap = map[string]string{
	"listen-port": "listen_port",
	"syslocation": "syslocation",
	"syscontact":  "syscontact",
	"sysname":     "sysname",
	"version":     "version",
	"community":   "community",
	"source":      "source",
	"rw":          "rw",
	"username":    "username",
	"security":    "security",
	"auth-proto":  "auth_proto",
	"auth-pass":   "auth_pass",
	"priv-proto":  "priv_proto",
	"priv-pass":   "priv_pass",
	"enabled":     "enabled",
}

// ---------- create defaults ----------

var ftpCreateDefaults = map[string]interface{}{
	"enabled":  "yes",
	"upload":   0,
	"download": 0,
}

var httpCreateDefaults = map[string]interface{}{
	"enabled":     "yes",
	"server_name": "",
	"access":      1,
}

var sambaCreateDefaults = map[string]interface{}{
	"enabled":    "yes",
	"browseable": "",
}

var (
	ftpConfigInputFields = []string{"open_ftp", "ftp_port", "ftp_access"}
	ftpUserInputFields   = []string{"enabled", "username", "passwd", "permission", "home_dir", "upload", "download", "tagname"}
	httpUserInputFields  = []string{"enabled", "tagname", "http_port", "server_name", "ssl_on", "autoindex", "download", "home_dir", "access"}
	sambaConfigFields    = []string{"enabled", "workgroup", "wsdd2", "access"}
	sambaUserInputFields = []string{"enabled", "username", "passwd", "name", "perm", "guest", "home_dir", "browseable", "tagname"}
	snmpdConfigFields    = []string{
		"enabled", "listen_port", "syslocation", "syscontact", "sysname", "version", "community",
		"source", "rw", "username", "security", "auth_proto", "auth_pass", "priv_proto", "priv_pass",
	}
)

// ---------- flag adders ----------

func addFTPFlags(cmd *cobra.Command) {
	cmd.Flags().String("username", "", "FTP username")
	cmd.Flags().String("password", "", "FTP password")
	cmd.Flags().String("permission", "", "Permission level")
	cmd.Flags().String("home-dir", "", "Home directory")
	cmd.Flags().String("upload", "", "Upload bandwidth limit")
	cmd.Flags().String("download", "", "Download bandwidth limit")
	cliapp.AddEnabledFlag(cmd)
}

func addFTPUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Account name (tagname)")
	addFTPFlags(cmd)
}

func addHTTPFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Account name (tagname)")
	cmd.Flags().String("port", "", "HTTP listen port")
	cmd.Flags().String("server-name", "", "Server name")
	cmd.Flags().String("ssl", "", "SSL on/off (0=off, 1=on)")
	cmd.Flags().String("autoindex", "", "Enable directory listing (0=off, 1=on)")
	cmd.Flags().String("download", "", "Download bandwidth limit (KByte/s, 0=unlimited)")
	cmd.Flags().String("home-dir", "", "Home directory path")
	cmd.Flags().String("access", "", "WAN access (0=deny, 1=allow)")
	cliapp.AddEnabledFlag(cmd)
}

func addSambaFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Share name")
	cmd.Flags().String("username", "", "Samba username")
	cmd.Flags().String("password", "", "Samba password")
	cmd.Flags().String("permission", "", "Permission (rw/ro)")
	cmd.Flags().String("guest", "", "Guest access (yes/no)")
	cmd.Flags().String("home-dir", "", "Share directory path")
	cliapp.AddEnabledFlag(cmd)
}

func addFTPConfigFlags(cmd *cobra.Command) {
	cmd.Flags().String("open-ftp", "", "FTP service switch (0=off, 1=on)")
	cmd.Flags().String("ftp-port", "", "FTP port (1-65535)")
	cmd.Flags().String("ftp-access", "", "WAN access (0=deny, 1=allow)")
}

func addSambaConfigFlags(cmd *cobra.Command) {
	cmd.Flags().String("enabled", "", "Service enabled (yes/no)")
	cmd.Flags().String("workgroup", "", "Workgroup name")
	cmd.Flags().String("wsdd2", "", "Network discovery (0=off, 1=on)")
	cmd.Flags().String("access", "", "WAN access (0=deny, 1=allow)")
}

func addSNMPDFlags(cmd *cobra.Command) {
	cmd.Flags().String("listen-port", "", "SNMPD listen port (2-65535)")
	cmd.Flags().String("syslocation", "", "System location")
	cmd.Flags().String("syscontact", "", "System contact")
	cmd.Flags().String("sysname", "", "System name")
	cmd.Flags().String("version", "", "SNMP version (2 or 3)")
	cmd.Flags().String("community", "", "Community string (v2)")
	cmd.Flags().String("source", "", "Allowed IP/subnet")
	cmd.Flags().String("rw", "", "Read/write permission (ro/rw)")
	cmd.Flags().String("username", "", "Username (v3)")
	cmd.Flags().String("security", "", "Security level (authNoPriv/authPriv)")
	cmd.Flags().String("auth-proto", "", "Auth protocol (MD5/SHA)")
	cmd.Flags().String("auth-pass", "", "Auth password (8-30 chars)")
	cmd.Flags().String("priv-proto", "", "Privacy protocol (DES/AES)")
	cmd.Flags().String("priv-pass", "", "Privacy password (8-30 chars)")
	cliapp.AddEnabledFlag(cmd)
}

func addAdvancedListFlags(cmd *cobra.Command) {
	cliapp.AddPaginationFlags(cmd)
	cmd.Flags().String("key", "", "Search field")
	cmd.Flags().String("pattern", "", "Search pattern")
	cmd.Flags().String("order", "", "Sort direction: asc|desc")
	cmd.Flags().String("order-by", "", "Sort field")
	_ = cmd.RegisterFlagCompletionFunc("order", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"asc", "desc"}, cobra.ShellCompDirectiveNoFileComp
	})
}

func advancedListParams(cmd *cobra.Command) map[string]string {
	page, _ := cmd.Flags().GetInt("page")
	pageSize, _ := cmd.Flags().GetInt("page-size")
	params := map[string]string{
		"page":  intString(page),
		"limit": intString(pageSize),
	}
	if key, _ := cmd.Flags().GetString("key"); key != "" {
		params["key"] = key
	}
	if pattern, _ := cmd.Flags().GetString("pattern"); pattern != "" {
		params["pattern"] = pattern
	}
	if orderBy, _ := cmd.Flags().GetString("order-by"); orderBy != "" {
		params["order"] = orderBy
	}
	if order, _ := cmd.Flags().GetString("order"); order != "" {
		params["order_by"] = order
	}
	return params
}

func intString(v int) string {
	return strconv.Itoa(v)
}

// ---------- top-level command ----------

func New(app *cliapp.Runtime) *cobra.Command {
	advancedCmd := &cobra.Command{
		Use:   "advanced",
		Short: "Advanced services",
	}

	advancedCmd.AddCommand(userGroup(app, "http", "HTTP server user management", "advanced-service/http-users", "", "", addHTTPFlags, httpFieldMap, httpCreateDefaults,
		[]string{"name", "port", "ssl", "autoindex", "download", "home-dir"},
		httpUserInputFields,
		[]string{"id", "tagname", "http_port", "ssl_on", "home_dir", "download", "enabled"}))
	advancedCmd.AddCommand(serviceGroup(app, "ftp", "FTP server management", "advanced-service/ftp-config", "advanced-service/ftp-users", addFTPFlags, addFTPUpdateFlags, ftpCreateFieldMap, ftpUpdateFieldMap, addFTPConfigFlags, ftpConfigFieldMap, ftpCreateDefaults,
		[]string{"username", "password", "permission", "home-dir"},
		ftpUserInputFields,
		ftpConfigInputFields,
		[]string{"id", "username", "permission", "home_dir", "upload", "download", "enabled"}))
	advancedCmd.AddCommand(serviceGroup(app, "samba", "Samba share management", "advanced-service/samba-config", "advanced-service/samba-users", addSambaFlags, addSambaFlags, sambaFieldMap, sambaFieldMap, addSambaConfigFlags, sambaConfigFieldMap, sambaCreateDefaults,
		[]string{"name", "username", "password", "permission", "guest", "home-dir"},
		sambaUserInputFields,
		sambaConfigFields,
		[]string{"id", "name", "username", "perm", "guest", "home_dir", "enabled"}))
	advancedCmd.AddCommand(snmpdGroup(app))
	return advancedCmd
}

func userGroup(app *cliapp.Runtime, use, short, userAPIPath, configGetPath, configSetPath string, addFlags func(*cobra.Command), fieldMap map[string]string, createDefaults map[string]interface{}, requiredCreateFlags, userInputFields []string, defaultColumns []string) *cobra.Command {
	return userGroupWithConfig(app, use, short, userAPIPath, configGetPath, configSetPath, addFlags, fieldMap, nil, nil, createDefaults, requiredCreateFlags, userInputFields, nil, defaultColumns)
}

func userGroupWithConfig(app *cliapp.Runtime, use, short, userAPIPath, configGetPath, configSetPath string, addFlags func(*cobra.Command), fieldMap map[string]string, cfgAddFlags func(*cobra.Command), cfgFieldMap map[string]string, createDefaults map[string]interface{}, requiredCreateFlags, userInputFields, configInputFields []string, defaultColumns []string) *cobra.Command {
	return userGroupWithSeparateUpdate(app, use, short, userAPIPath, configGetPath, configSetPath, addFlags, addFlags, fieldMap, fieldMap, cfgAddFlags, cfgFieldMap, createDefaults, requiredCreateFlags, userInputFields, configInputFields, defaultColumns)
}

func userGroupWithSeparateUpdate(app *cliapp.Runtime, use, short, userAPIPath, configGetPath, configSetPath string, addCreateFlags, addUpdateFlags func(*cobra.Command), createFieldMap, updateFieldMap map[string]string, cfgAddFlags func(*cobra.Command), cfgFieldMap map[string]string, createDefaults map[string]interface{}, requiredCreateFlags, userInputFields, configInputFields []string, defaultColumns []string) *cobra.Command {
	group := &cobra.Command{Use: use, Short: short}

	if configGetPath != "" && configSetPath != "" {
		group.AddCommand(
			&cobra.Command{
				Use:   "config-get",
				Short: "Get " + use + " server configuration",
				RunE: func(cmd *cobra.Command, args []string) error {
					if err := app.RequireAuth(); err != nil {
						return err
					}
					raw, err := app.APIClient.Get(cliapp.APIBase+"/"+configGetPath, nil)
					if err != nil {
						return err
					}
					app.PrintRaw(raw)
					return nil
				},
			},
			configSetCmd(app, "config-set", "Update "+use+" server configuration", cfgAddFlags, cfgFieldMap, configInputFields, cliapp.APIBase+"/"+configGetPath, func(body interface{}, id string) (json.RawMessage, error) {
				return app.APIClient.Put(cliapp.APIBase+"/"+configSetPath, body)
			}),
		)
	}

	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List " + use + " users",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if len(defaultColumns) > 0 {
				app.DefaultColumns = defaultColumns
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/"+userAPIPath, advancedListParams(cmd))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	addAdvancedListFlags(listCmd)

	createCmd := dataCmd(app, "create", "Create a "+use+" user", addCreateFlags, createFieldMap, createDefaults, func(body interface{}, id string) (json.RawMessage, error) {
		return app.APIClient.Post(cliapp.APIBase+"/"+userAPIPath, body)
	})
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
	group.AddCommand(
		listCmd,
		createCmd,
		updateCmd(app, "update ID", "Update a "+use+" user", addUpdateFlags, updateFieldMap, userInputFields, cliapp.APIBase+"/"+userAPIPath, func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Put(cliapp.APIBase+"/"+userAPIPath+"/"+id, body)
		}),
		dataCmdWithID(app, "toggle ID", "Enable/disable a "+use+" user", cliapp.AddEnabledFlag, map[string]string{"enabled": "enabled"}, func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Patch(cliapp.APIBase+"/"+userAPIPath+"/"+id, body)
		}),
		deleteByIDCmd(app, "delete ID", "Delete a "+use+" user", "/"+userAPIPath+"/"),
	)

	return group
}

func serviceGroup(app *cliapp.Runtime, use, short, configAPIPath, userAPIPath string, addFlags, addUpdateFlags func(*cobra.Command), createFieldMap, updateFieldMap map[string]string, configAddFlags func(*cobra.Command), configFieldMap map[string]string, createDefaults map[string]interface{}, requiredCreateFlags, userInputFields, configInputFields []string, defaultColumns []string) *cobra.Command {
	return userGroupWithSeparateUpdate(app, use, short, userAPIPath, configAPIPath, configAPIPath, addFlags, addUpdateFlags, createFieldMap, updateFieldMap, configAddFlags, configFieldMap, createDefaults, requiredCreateFlags, userInputFields, configInputFields, defaultColumns)
}

func snmpdGroup(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "snmpd", Short: "SNMPD config"}
	group.AddCommand(
		&cobra.Command{
			Use:   "get",
			Short: "Get SNMPD config",
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := app.RequireAuth(); err != nil {
					return err
				}
				app.DefaultColumns = []string{"id", "enabled", "listen_port", "version", "community", "rw"}
				raw, err := app.APIClient.Get(cliapp.APIBase+"/advanced-service/snmpd-config", nil)
				if err != nil {
					return err
				}
				app.PrintRaw(raw)
				return nil
			},
		},
		configSetCmd(app, "set", "Update SNMPD configuration", addSNMPDFlags, snmpdFieldMap, snmpdConfigFields, cliapp.APIBase+"/advanced-service/snmpd-config", func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Put(cliapp.APIBase+"/advanced-service/snmpd-config", body)
		}),
	)
	return group
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

func dataCmd(app *cliapp.Runtime, use, short string, addFlags func(*cobra.Command), fieldMap map[string]string, defaults map[string]interface{}, fn callWithBody) *cobra.Command {
	return dataCmdImpl(app, use, short, false, addFlags, fieldMap, defaults, fn)
}

func dataCmdWithID(app *cliapp.Runtime, use, short string, addFlags func(*cobra.Command), fieldMap map[string]string, fn callWithBody) *cobra.Command {
	return dataCmdImpl(app, use, short, true, addFlags, fieldMap, nil, fn)
}

func configSetCmd(app *cliapp.Runtime, use, short string, addFlags func(*cobra.Command), fieldMap map[string]string, inputFields []string, getPath string, fn callWithBody) *cobra.Command {
	c := dataCmdImpl(app, use, short, false, addFlags, fieldMap, nil, func(body interface{}, id string) (json.RawMessage, error) {
		return fn(body, id)
	})
	c.RunE = func(cmd *cobra.Command, args []string) error {
		if err := app.RequireAuth(); err != nil {
			return err
		}
		data, _ := cmd.Flags().GetString("data")
		body, err := buildFullBody(app, cmd, data, fieldMap, inputFields, getPath)
		if err != nil {
			return err
		}
		raw, err := fn(body, "")
		if err != nil {
			return err
		}
		app.PrintRaw(raw)
		return nil
	}
	return c
}

func updateCmd(app *cliapp.Runtime, use, short string, addFlags func(*cobra.Command), fieldMap map[string]string, inputFields []string, listPath string, fn callWithBody) *cobra.Command {
	c := dataCmdImpl(app, use, short, true, addFlags, fieldMap, nil, func(body interface{}, id string) (json.RawMessage, error) {
		return fn(body, id)
	})
	c.RunE = func(cmd *cobra.Command, args []string) error {
		if err := app.RequireAuth(); err != nil {
			return err
		}
		data, _ := cmd.Flags().GetString("data")
		id := args[0]
		body, err := buildFullBodyFromList(app, cmd, data, fieldMap, inputFields, listPath, id)
		if err != nil {
			return err
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

func dataCmdImpl(app *cliapp.Runtime, use, short string, withID bool, addFlags func(*cobra.Command), fieldMap map[string]string, defaults map[string]interface{}, fn callWithBody) *cobra.Command {
	c := &cobra.Command{Use: use, Short: short}
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
	current, err := extractInputObject(raw, inputFields, "")
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

func buildFullBodyFromList(app *cliapp.Runtime, cmd *cobra.Command, data string, fieldMap map[string]string, inputFields []string, listPath, id string) (map[string]interface{}, error) {
	changes, err := cliapp.MergeDataWithFlags(data, cmd, fieldMap)
	if err != nil {
		return nil, err
	}
	readClient := app.APIClient
	if app.APIClient.DryRun {
		readClient = app.NewClient(app.Session.BaseURL, app.Session.Token)
	}
	raw, err := readClient.Get(listPath, map[string]string{"page": "1", "limit": "500"})
	if err != nil {
		if hasAllInputFields(changes, inputFields) {
			return changes, nil
		}
		return nil, err
	}
	current, err := extractInputObject(raw, inputFields, id)
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

func extractInputObject(raw json.RawMessage, inputFields []string, targetID string) (map[string]interface{}, error) {
	var value interface{}
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, err
	}
	obj, ok := findAPIObject(value, targetID)
	if !ok {
		return nil, &cliapp.ValidationError{Message: "unexpected advanced API response"}
	}
	if len(inputFields) == 0 {
		result := map[string]interface{}{}
		for k, v := range obj {
			result[k] = v
		}
		delete(result, "id")
		return result, nil
	}
	result := map[string]interface{}{}
	for _, key := range inputFields {
		if v, ok := obj[key]; ok {
			result[key] = v
		}
	}
	return result, nil
}

func findAPIObject(value interface{}, targetID string) (map[string]interface{}, bool) {
	switch v := value.(type) {
	case map[string]interface{}:
		if targetID != "" && matchesID(v, targetID) {
			return v, true
		}
		for _, key := range []string{"data", "dir_data", "results"} {
			if nested, ok := v[key]; ok {
				if obj, found := findAPIObject(nested, targetID); found {
					return obj, true
				}
			}
		}
		if targetID != "" {
			return nil, false
		}
		return v, true
	case []interface{}:
		if len(v) == 0 {
			return nil, false
		}
		if targetID != "" {
			for _, item := range v {
				if obj, ok := item.(map[string]interface{}); ok && matchesID(obj, targetID) {
					return obj, true
				}
			}
			return nil, false
		}
		if obj, ok := v[0].(map[string]interface{}); ok {
			return obj, true
		}
		return findAPIObject(v[0], targetID)
	default:
		return nil, false
	}
}

func matchesID(obj map[string]interface{}, targetID string) bool {
	id, ok := obj["id"]
	if !ok {
		return false
	}
	if fmt.Sprint(id) == targetID {
		return true
	}
	if n, ok := id.(float64); ok {
		return strconv.FormatInt(int64(n), 10) == targetID
	}
	return false
}

func hasAllInputFields(body map[string]interface{}, inputFields []string) bool {
	for _, key := range inputFields {
		if _, ok := body[key]; !ok {
			return false
		}
	}
	return true
}
