package logcmd

import (
	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

type logResource struct {
	name           string
	apiPath        string
	defaultColumns []string
}

var logResources = []logResource{
	{name: "arp", apiPath: "log/arp", defaultColumns: []string{"id", "timestamp", "content"}},
	{name: "auth", apiPath: "log/auth", defaultColumns: []string{"id", "timestamp", "username", "macip", "ppptype", "result", "ip_addr", "event", "interface"}},
	{name: "dhcp", apiPath: "log/dhcp", defaultColumns: []string{"id", "timestamp", "msgtype", "interface", "ip_addr", "mac", "hostname"}},
	{name: "pppoe", apiPath: "log/pppoe", defaultColumns: []string{"id", "timestamp", "interface", "content"}},
	{name: "system", apiPath: "log/system", defaultColumns: []string{"id", "timestamp", "content"}},
	{name: "web", apiPath: "log/web_activity", defaultColumns: []string{"id", "timestamp", "username", "ip_addr", "function", "event"}},
	{name: "ddns", apiPath: "log/ddns", defaultColumns: []string{"id", "timestamp", "domain", "status", "ip_addr", "message"}},
	{name: "notice", apiPath: "log/notice", defaultColumns: []string{"id", "timestamp", "type", "title", "content"}},
	{name: "wireless", apiPath: "log/wireless", defaultColumns: []string{"id", "timestamp", "action", "mac", "apmac", "ssid", "errmsg", "signal"}},
}

func New(app *cliapp.Runtime) *cobra.Command {
	logCmd := &cobra.Command{
		Use:   "log",
		Short: "System logs",
	}

	for _, resource := range logResources {
		logCmd.AddCommand(logGroup(app, resource))
	}
	return logCmd
}

func logGroup(app *cliapp.Runtime, resource logResource) *cobra.Command {
	group := &cobra.Command{
		Use:   resource.name,
		Short: resource.name + " log",
	}
	group.AddCommand(logListCmd(app, resource))
	group.AddCommand(logDeleteCmd(app, resource))
	return group
}

func logListCmd(app *cliapp.Runtime, resource logResource) *cobra.Command {
	c := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List " + resource.name + " logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = resource.defaultColumns
			raw, err := app.APIClient.Get(cliapp.APIBase+"/"+resource.apiPath, logListParams(cmd))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	addLogListFlags(c)
	return c
}

func logDeleteCmd(app *cliapp.Runtime, resource logResource) *cobra.Command {
	c := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"clear"},
		Short:   "Clear " + resource.name + " logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			yes, _ := cmd.Flags().GetBool("yes")
			if err := cliapp.ConfirmDelete(app.Stdout, app.Stderr, resource.name+" logs", "all", yes); err != nil {
				return err
			}
			raw, err := app.APIClient.Delete(cliapp.APIBase + "/" + resource.apiPath)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	c.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	return c
}

func addLogListFlags(cmd *cobra.Command) {
	cliapp.AddPaginationFlags(cmd)
	cmd.Flags().String("filter", "", "Filter expression")
	cmd.Flags().String("key", "", "Fields for fuzzy search")
	cmd.Flags().String("pattern", "", "Fuzzy search pattern")
	cmd.Flags().String("order", "", "Sort direction: asc|desc")
	cmd.Flags().String("order-by", "", "Sort field")
	_ = cmd.RegisterFlagCompletionFunc("order", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"asc", "desc"}, cobra.ShellCompDirectiveNoFileComp
	})
}

func logListParams(cmd *cobra.Command) map[string]string {
	page, _ := cmd.Flags().GetInt("page")
	pageSize, _ := cmd.Flags().GetInt("page-size")
	params := cliapp.ListParamsWithPageSizeKey(page, pageSize, "", "", "", "limit")
	for _, name := range []string{"filter", "key", "pattern", "order", "order-by"} {
		value, _ := cmd.Flags().GetString(name)
		if value == "" {
			continue
		}
		apiName := name
		if name == "order-by" {
			apiName = "order_by"
		}
		params[apiName] = value
	}
	return params
}
