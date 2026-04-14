package advanced

import (
	"encoding/json"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

// ---------- field maps ----------

var ftpFieldMap = map[string]string{
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
	"enabled": "yes",
	"access":  1,
}

var sambaCreateDefaults = map[string]interface{}{
	"enabled":    "yes",
	"browseable": "",
}

// ---------- flag adders ----------

func addFTPFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Account name (tagname)")
	cmd.Flags().String("username", "", "FTP username")
	cmd.Flags().String("password", "", "FTP password")
	cmd.Flags().String("permission", "", "Permission level")
	cmd.Flags().String("home-dir", "", "Home directory")
	cmd.Flags().String("upload", "", "Upload bandwidth limit")
	cmd.Flags().String("download", "", "Download bandwidth limit")
	cliapp.AddEnabledFlag(cmd)
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

// ---------- top-level command ----------

func New(app *cliapp.Runtime) *cobra.Command {
	advancedCmd := &cobra.Command{
		Use:   "advanced",
		Short: "Advanced services",
	}

	advancedCmd.AddCommand(userGroup(app, "http", "HTTP server user management", "advanced-service/http-users", "", "", addHTTPFlags, httpFieldMap, httpCreateDefaults,
		[]string{"name", "port", "ssl", "autoindex", "download", "home-dir"},
		[]string{"id", "tagname", "http_port", "ssl_on", "home_dir", "download", "enabled"}))
	advancedCmd.AddCommand(serviceGroup(app, "ftp", "FTP server management", "advanced-service/ftp-config", "advanced-service/ftp-users", addFTPFlags, ftpFieldMap, addFTPConfigFlags, ftpConfigFieldMap, ftpCreateDefaults,
		[]string{"username", "password", "permission", "home-dir"},
		[]string{"id", "username", "permission", "home_dir", "upload", "download", "enabled"}))
	advancedCmd.AddCommand(serviceGroup(app, "samba", "Samba share management", "advanced-service/samba-config", "advanced-service/samba-users", addSambaFlags, sambaFieldMap, addSambaConfigFlags, sambaConfigFieldMap, sambaCreateDefaults,
		[]string{"name", "username", "password", "permission", "guest", "home-dir"},
		[]string{"id", "name", "username", "perm", "guest", "home_dir", "enabled"}))
	advancedCmd.AddCommand(snmpdGroup(app))
	return advancedCmd
}

func userGroup(app *cliapp.Runtime, use, short, userAPIPath, configGetPath, configSetPath string, addFlags func(*cobra.Command), fieldMap map[string]string, createDefaults map[string]interface{}, requiredCreateFlags []string, defaultColumns []string) *cobra.Command {
	return userGroupWithConfig(app, use, short, userAPIPath, configGetPath, configSetPath, addFlags, fieldMap, nil, nil, createDefaults, requiredCreateFlags, defaultColumns)
}

func userGroupWithConfig(app *cliapp.Runtime, use, short, userAPIPath, configGetPath, configSetPath string, addFlags func(*cobra.Command), fieldMap map[string]string, cfgAddFlags func(*cobra.Command), cfgFieldMap map[string]string, createDefaults map[string]interface{}, requiredCreateFlags []string, defaultColumns []string) *cobra.Command {
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
			dataCmd(app, "config-set", "Update "+use+" server configuration", cfgAddFlags, cfgFieldMap, nil, func(body interface{}, id string) (json.RawMessage, error) {
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
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/"+userAPIPath,
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(listCmd)

	createCmd := dataCmd(app, "create", "Create a "+use+" user", addFlags, fieldMap, createDefaults, func(body interface{}, id string) (json.RawMessage, error) {
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
		dataCmdWithID(app, "update ID", "Update a "+use+" user", addFlags, fieldMap, func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Put(cliapp.APIBase+"/"+userAPIPath+"/"+id, body)
		}),
		dataCmdWithID(app, "toggle ID", "Enable/disable a "+use+" user", cliapp.AddEnabledFlag, map[string]string{"enabled": "enabled"}, func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Patch(cliapp.APIBase+"/"+userAPIPath+"/"+id, body)
		}),
		deleteByIDCmd(app, "delete ID", "Delete a "+use+" user", "/"+userAPIPath+"/"),
	)

	return group
}

func serviceGroup(app *cliapp.Runtime, use, short, configAPIPath, userAPIPath string, addFlags func(*cobra.Command), fieldMap map[string]string, configAddFlags func(*cobra.Command), configFieldMap map[string]string, createDefaults map[string]interface{}, requiredCreateFlags []string, defaultColumns []string) *cobra.Command {
	return userGroupWithConfig(app, use, short, userAPIPath, configAPIPath, configAPIPath, addFlags, fieldMap, configAddFlags, configFieldMap, createDefaults, requiredCreateFlags, defaultColumns)
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
		dataCmd(app, "set", "Update SNMPD configuration", addSNMPDFlags, snmpdFieldMap, nil, func(body interface{}, id string) (json.RawMessage, error) {
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
