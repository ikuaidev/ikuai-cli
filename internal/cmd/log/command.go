package logcmd

import (
	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

func New(app *cliapp.Runtime) *cobra.Command {
	logCmd := &cobra.Command{
		Use:   "log",
		Short: "System logs",
	}

	logCmd.AddCommand(flatLogListCmd(app, "arp", "log/arp"))
	logCmd.AddCommand(flatLogClearCmd(app, "arp", "log/arp"))
	logCmd.AddCommand(flatLogListCmd(app, "auth", "log/auth"))
	logCmd.AddCommand(flatLogClearCmd(app, "auth", "log/auth"))
	logCmd.AddCommand(flatLogListCmd(app, "dhcp", "log/dhcp"))
	logCmd.AddCommand(flatLogClearCmd(app, "dhcp", "log/dhcp"))
	logCmd.AddCommand(flatLogListCmd(app, "pppoe", "log/pppoe"))
	logCmd.AddCommand(flatLogClearCmd(app, "pppoe", "log/pppoe"))
	logCmd.AddCommand(flatLogListCmd(app, "system", "log/system"))
	logCmd.AddCommand(flatLogClearCmd(app, "system", "log/system"))
	logCmd.AddCommand(flatLogListCmd(app, "web", "log/web_activity"))
	logCmd.AddCommand(flatLogClearCmd(app, "web", "log/web_activity"))
	logCmd.AddCommand(flatLogListCmd(app, "ddns", "log/ddns"))
	logCmd.AddCommand(flatLogClearCmd(app, "ddns", "log/ddns"))
	logCmd.AddCommand(flatLogListCmd(app, "notice", "log/notice"))
	logCmd.AddCommand(flatLogClearCmd(app, "notice", "log/notice"))
	logCmd.AddCommand(flatLogListCmd(app, "wireless", "log/wireless"))
	logCmd.AddCommand(flatLogClearCmd(app, "wireless", "log/wireless"))
	return logCmd
}

func flatLogListCmd(app *cliapp.Runtime, name, apiPath string) *cobra.Command {
	c := &cobra.Command{
		Use:   name,
		Short: "Show " + name + " log",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
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
	cliapp.AddListFlags(c)
	return c
}

func flatLogClearCmd(app *cliapp.Runtime, name, apiPath string) *cobra.Command {
	return &cobra.Command{
		Use:   name + "-clear",
		Short: "Clear " + name + " log",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Delete(cliapp.APIBase + "/" + apiPath)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
}
