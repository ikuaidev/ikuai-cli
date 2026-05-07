package authserver

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

var authServerDefaultColumns = []string{
	"id", "enabled", "interface", "idle_time", "max_time", "user_auth",
	"coupon_auth", "phone_auth", "static_pwd", "nopasswd", "https_redirect",
}

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

var authServerIntegerFields = map[string]bool{
	"max_time":       true,
	"idle_time":      true,
	"user_auth":      true,
	"coupon_auth":    true,
	"phone_auth":     true,
	"static_pwd":     true,
	"nopasswd":       true,
	"weixin":         true,
	"https_redirect": true,
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

	authServerSetCmd := &cobra.Command{
		Use:   "set",
		Short: "Update web auth config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			changes, err := mergeAuthServerChanges(data, cmd)
			if err != nil {
				return err
			}
			if len(changes) == 0 {
				return &cliapp.ValidationError{Message: "at least one config field is required"}
			}
			body, err := buildAuthServerSetBody(app, changes)
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

func mergeAuthServerChanges(data string, cmd *cobra.Command) (map[string]interface{}, error) {
	changes, err := cliapp.MergeDataWithFlags(data, cmd, authServerFieldMap)
	if err != nil {
		return nil, err
	}
	for key, value := range changes {
		if !authServerIntegerFields[key] {
			continue
		}
		switch typed := value.(type) {
		case string:
			if typed == "" {
				continue
			}
			n, err := strconv.ParseInt(typed, 10, 64)
			if err != nil {
				return nil, &cliapp.ValidationError{Message: fmt.Sprintf("invalid integer for %s: %s", key, typed)}
			}
			changes[key] = n
		}
	}
	return changes, nil
}

func buildAuthServerSetBody(app *cliapp.Runtime, changes map[string]interface{}) (map[string]interface{}, error) {
	readClient := app.APIClient
	if app.APIClient.DryRun {
		readClient = app.NewClient(app.Session.BaseURL, app.Session.Token)
	}
	raw, err := readClient.Get(cliapp.APIBase+"/auth/web/services", nil)
	if err != nil {
		return nil, err
	}
	current, err := extractAuthServerConfig(raw)
	if err != nil {
		return nil, err
	}
	for key, value := range changes {
		current[key] = value
	}
	return current, nil
}

func extractAuthServerConfig(raw json.RawMessage) (map[string]interface{}, error) {
	var value interface{}
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, err
	}
	if obj, ok := findAuthServerConfigObject(value); ok {
		return obj, nil
	}
	return nil, &cliapp.ValidationError{Message: "empty auth-server config response"}
}

func findAuthServerConfigObject(value interface{}) (map[string]interface{}, bool) {
	switch typed := value.(type) {
	case []interface{}:
		if len(typed) == 0 {
			return nil, false
		}
		return findAuthServerConfigObject(typed[0])
	case map[string]interface{}:
		for _, key := range []string{"data", "results"} {
			if inner, ok := typed[key]; ok {
				if obj, found := findAuthServerConfigObject(inner); found {
					return obj, true
				}
			}
		}
		return typed, true
	default:
		return nil, false
	}
}
