package completion

import (
	"fmt"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

func New(app *cliapp.Runtime) *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate completion scripts",
		Long: `Generate shell completion scripts.

Supported shells:
  bash
  zsh
  fish
  powershell`,
	}

	completionCmd.AddCommand(
		shellCmd(app, "bash", "Generate Bash completion script", func(cmd *cobra.Command) error {
			return cmd.Root().GenBashCompletionV2(app.Stdout, true)
		}),
		shellCmd(app, "zsh", "Generate Zsh completion script", func(cmd *cobra.Command) error {
			return cmd.Root().GenZshCompletion(app.Stdout)
		}),
		shellCmd(app, "fish", "Generate Fish completion script", func(cmd *cobra.Command) error {
			return cmd.Root().GenFishCompletion(app.Stdout, true)
		}),
		shellCmd(app, "powershell", "Generate PowerShell completion script", func(cmd *cobra.Command) error {
			return cmd.Root().GenPowerShellCompletionWithDesc(app.Stdout)
		}),
	)

	return completionCmd
}

func shellCmd(app *cliapp.Runtime, use, short string, run func(cmd *cobra.Command) error) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if app.Stdout == nil {
				return fmt.Errorf("stdout is not configured")
			}
			return run(cmd)
		},
	}
}
