package auth

import (
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/ikuaidev/ikuai-cli/internal/session"
	"github.com/spf13/cobra"
)

type setURLResult struct {
	Message    string `json:"message"`
	BaseURL    string `json:"base_url"`
	Normalized bool   `json:"normalized"`
	InputURL   string `json:"input_url,omitempty"`
}

func New(app *cliapp.Runtime) *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication",
		Long:  `Manage authentication credentials for the iKuai router API.`,
		Example: `  ikuai-cli auth set-url https://192.168.1.1
  ikuai-cli auth set-token MGFjYzg1ZjMt...
  ikuai-cli auth status
  ikuai-cli auth clear`,
	}

	authSetURLCmd := &cobra.Command{
		Use:   "set-url [URL]",
		Short: "Set router base URL",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			urlVal, _ := cmd.Flags().GetString("url")
			if len(args) > 0 && urlVal == "" {
				urlVal = args[0]
			}
			if strings.TrimSpace(urlVal) == "" {
				return &cliapp.ValidationError{Message: "URL is required: use --url or pass as argument"}
			}
			inputURL := urlVal
			baseURL, normalized, err := session.NormalizeBaseURL(urlVal)
			if err != nil {
				return &cliapp.ValidationError{Message: invalidRouterURLMessage(err)}
			}
			if err := session.SaveBaseURL(baseURL); err != nil {
				return err
			}
			result := setURLResult{
				Message:    "Base URL saved",
				BaseURL:    baseURL,
				Normalized: normalized,
			}
			if normalized {
				result.InputURL = safeInputURLForOutput(inputURL)
			}
			app.PrintJSON(result)
			return nil
		},
	}
	authSetURLCmd.Flags().String("url", "", "Router base URL (e.g. https://192.168.1.1)")

	authSetTokenCmd := &cobra.Command{
		Use:   "set-token [TOKEN]",
		Short: "Set API Bearer token",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tokenVal, _ := cmd.Flags().GetString("token")
			if len(args) > 0 && tokenVal == "" {
				tokenVal = args[0]
			}
			if tokenVal == "" {
				return &cliapp.ValidationError{Message: "TOKEN is required: use --token or pass as argument"}
			}
			if err := session.SaveToken(tokenVal); err != nil {
				return err
			}
			app.PrintJSON(map[string]string{
				"message": "Token saved",
			})
			return nil
		},
	}
	authSetTokenCmd.Flags().String("token", "", "API Bearer token")

	authClearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear saved base URL and token (SSH credentials preserved)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := session.Clear(); err != nil {
				return err
			}
			app.PrintJSON(map[string]string{
				"message": "Cleared",
			})
			return nil
		},
	}

	authStatusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show session info",
		RunE: func(cmd *cobra.Command, args []string) error {
			app.PrintJSON(map[string]string{
				"base_url": app.Session.BaseURL,
				"source":   app.CredSource,
			})
			return nil
		},
	}

	authCmd.AddCommand(authSetURLCmd, authSetTokenCmd, authClearCmd, authStatusCmd)
	return authCmd
}

func invalidRouterURLMessage(err error) string {
	return "Invalid router URL.\n\n" + err.Error() + "\n\nTry:\n  ikuai-cli auth set-url 192.168.1.1\n  ikuai-cli auth set-url https://192.168.1.1"
}

func safeInputURLForOutput(input string) string {
	input = strings.TrimSpace(input)
	if strings.ContainsAny(input, "?#@") {
		return ""
	}
	return input
}
