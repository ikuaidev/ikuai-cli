package authserver

import (
	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

var authServerDefaultColumns = []string{
	"id", "enabled", "interface", "idle_time", "max_time", "user_auth",
	"coupon_auth", "phone_auth", "static_pwd", "nopasswd", "https_redirect",
}

func New(app *cliapp.Runtime) *cobra.Command {
	authServerCmd := &cobra.Command{
		Use:   "auth-server",
		Short: "Web auth server",
	}

	authServerGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get web auth config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = authServerDefaultColumns
			raw, err := app.APIClient.Get(cliapp.APIBase+"/auth/web/services", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	authServerCmd.AddCommand(authServerGetCmd)
	return authServerCmd
}
