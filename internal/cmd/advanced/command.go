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
	"enabled":     "enabled",
}

var sambaFieldMap = map[string]string{
	"name":       "tagname",
	"username":   "username",
	"password":   "passwd",
	"permission": "perm",
	"guest":      "guest",
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
	"version":     "version",
	"community":   "community",
	"username":    "username",
	"enabled":     "enabled",
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
	cmd.Flags().String("ssl", "", "SSL on/off")
	cmd.Flags().String("autoindex", "", "Enable directory listing")
	cmd.Flags().String("download", "", "Download bandwidth limit")
	cmd.Flags().String("home-dir", "", "Home directory")
	cliapp.AddEnabledFlag(cmd)
}

func addSambaFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Account name (tagname)")
	cmd.Flags().String("username", "", "Samba username")
	cmd.Flags().String("password", "", "Samba password")
	cmd.Flags().String("permission", "", "Permission level")
	cmd.Flags().String("guest", "", "Guest access (yes/no)")
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
	cmd.Flags().String("listen-port", "", "SNMPD listen port")
	cmd.Flags().String("version", "", "SNMP version")
	cmd.Flags().String("community", "", "SNMP community string")
	cmd.Flags().String("username", "", "SNMP username")
	cliapp.AddEnabledFlag(cmd)
}

// ---------- top-level command ----------

func New(app *cliapp.Runtime) *cobra.Command {
	advancedCmd := &cobra.Command{
		Use:   "advanced",
		Short: "Advanced services",
	}

	advancedCmd.AddCommand(userGroup(app, "http", "HTTP server user management", "advanced-service/http-users", "", "", addHTTPFlags, httpFieldMap,
		[]string{"name", "port", "ssl", "autoindex", "download", "home-dir"}))
	advancedCmd.AddCommand(serviceGroup(app, "ftp", "FTP server management", "advanced-service/ftp-config", "advanced-service/ftp-users", addFTPFlags, ftpFieldMap, addFTPConfigFlags, ftpConfigFieldMap,
		[]string{"username", "password", "permission", "home-dir"}))
	advancedCmd.AddCommand(serviceGroup(app, "samba", "Samba share management", "advanced-service/samba-config", "advanced-service/samba-users", addSambaFlags, sambaFieldMap, addSambaConfigFlags, sambaConfigFieldMap,
		[]string{"name", "username", "password", "permission", "guest"}))
	advancedCmd.AddCommand(snmpdGroup(app))
	return advancedCmd
}

func userGroup(app *cliapp.Runtime, use, short, userAPIPath, configGetPath, configSetPath string, addFlags func(*cobra.Command), fieldMap map[string]string, requiredCreateFlags []string) *cobra.Command {
	return userGroupWithConfig(app, use, short, userAPIPath, configGetPath, configSetPath, addFlags, fieldMap, nil, nil, requiredCreateFlags)
}

func userGroupWithConfig(app *cliapp.Runtime, use, short, userAPIPath, configGetPath, configSetPath string, addFlags func(*cobra.Command), fieldMap map[string]string, cfgAddFlags func(*cobra.Command), cfgFieldMap map[string]string, requiredCreateFlags []string) *cobra.Command {
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
			dataCmd(app, "config-set", "Update "+use+" server configuration", cfgAddFlags, cfgFieldMap, func(body interface{}, id string) (json.RawMessage, error) {
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

	createCmd := dataCmd(app, "create", "Create a "+use+" user", addFlags, fieldMap, func(body interface{}, id string) (json.RawMessage, error) {
		return app.APIClient.Post(cliapp.APIBase+"/"+userAPIPath, body)
	})
	if len(requiredCreateFlags) > 0 {
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
		getByIDCmd(app, "get ID", "Get a "+use+" user", "/"+userAPIPath+"/"),
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

func serviceGroup(app *cliapp.Runtime, use, short, configAPIPath, userAPIPath string, addFlags func(*cobra.Command), fieldMap map[string]string, configAddFlags func(*cobra.Command), configFieldMap map[string]string, requiredCreateFlags []string) *cobra.Command {
	return userGroupWithConfig(app, use, short, userAPIPath, configAPIPath, configAPIPath, addFlags, fieldMap, configAddFlags, configFieldMap, requiredCreateFlags)
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
				raw, err := app.APIClient.Get(cliapp.APIBase+"/advanced-service/snmpd-config", nil)
				if err != nil {
					return err
				}
				app.PrintRaw(raw)
				return nil
			},
		},
		dataCmd(app, "set", "Update SNMPD configuration", addSNMPDFlags, snmpdFieldMap, func(body interface{}, id string) (json.RawMessage, error) {
			return app.APIClient.Put(cliapp.APIBase+"/advanced-service/snmpd-config", body)
		}),
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

func dataCmd(app *cliapp.Runtime, use, short string, addFlags func(*cobra.Command), fieldMap map[string]string, fn callWithBody) *cobra.Command {
	return dataCmdImpl(app, use, short, false, addFlags, fieldMap, fn)
}

func dataCmdWithID(app *cliapp.Runtime, use, short string, addFlags func(*cobra.Command), fieldMap map[string]string, fn callWithBody) *cobra.Command {
	return dataCmdImpl(app, use, short, true, addFlags, fieldMap, fn)
}

func dataCmdImpl(app *cliapp.Runtime, use, short string, withID bool, addFlags func(*cobra.Command), fieldMap map[string]string, fn callWithBody) *cobra.Command {
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
	return c
}
