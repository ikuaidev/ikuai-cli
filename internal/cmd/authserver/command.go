package authserver

import (
	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

var authServerFieldMap = map[string]string{
	"enabled":         "enabled",
	"max-time":        "max_time",
	"idle-time":       "idle_time",
	"user-auth":       "user_auth",
	"coupon-auth":     "coupon_auth",
	"phone-auth":      "phone_auth",
	"static-pwd":      "static_pwd",
	"nopasswd":        "nopasswd",
	"weixin":          "weixin",
	"interface":       "interface",
	"passwd":          "passwd",
	"whitelist":       "whitelist",
	"whitelist-https": "whitelist_https",
	"whiteip":         "whiteip",
	"noauth-mac":      "noauth_mac",
	"radius-ip":       "radius_ip",
	"radius-key":      "radius_key",
	"https-redirect":  "https_redirect",
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
			raw, err := app.APIClient.Get(cliapp.APIBase+"/auth/web/services", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	authServerSetCmd := &cobra.Command{
		Use:   "set",
		Short: "Update web auth config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, authServerFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/auth/web/services", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	authServerSetCmd.Flags().String("data", "{}", "JSON body")
	authServerSetCmd.Flags().String("enabled", "", "Service enabled (yes/no)")
	authServerSetCmd.Flags().String("max-time", "", "Re-auth timeout (minutes, 0=unlimited)")
	authServerSetCmd.Flags().String("idle-time", "", "Idle timeout (seconds, 0=unlimited)")
	authServerSetCmd.Flags().String("user-auth", "", "User/password auth (0/1)")
	authServerSetCmd.Flags().String("coupon-auth", "", "Coupon auth (0/1)")
	authServerSetCmd.Flags().String("phone-auth", "", "Phone auth (0/1)")
	authServerSetCmd.Flags().String("static-pwd", "", "Static password auth (0/1)")
	authServerSetCmd.Flags().String("nopasswd", "", "One-click auth (0/1)")
	authServerSetCmd.Flags().String("weixin", "", "WeChat auth (0/1)")
	authServerSetCmd.Flags().String("interface", "", "Auth interface")
	authServerSetCmd.Flags().String("passwd", "", "Static password (MD5)")
	authServerSetCmd.Flags().String("whitelist", "", "HTTP whitelist domains (comma-separated)")
	authServerSetCmd.Flags().String("whitelist-https", "", "HTTPS whitelist domains (comma-separated)")
	authServerSetCmd.Flags().String("whiteip", "", "Whitelist IPs (comma-separated)")
	authServerSetCmd.Flags().String("noauth-mac", "", "No-auth MAC addresses (comma-separated)")
	authServerSetCmd.Flags().String("radius-ip", "", "RADIUS server IP")
	authServerSetCmd.Flags().String("radius-key", "", "RADIUS shared secret")
	authServerSetCmd.Flags().String("https-redirect", "", "HTTPS redirect to portal (0/1)")

	authServerCmd.AddCommand(authServerGetCmd, authServerSetCmd)
	return authServerCmd
}
