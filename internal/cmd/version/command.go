package version

import (
	"fmt"

	"github.com/ikuaidev/ikuai-cli/internal/buildinfo"
	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/ikuaidev/ikuai-cli/internal/output"
	"github.com/spf13/cobra"
)

func New(app *cliapp.Runtime) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			info := map[string]string{
				"name":    "ikuai-cli",
				"version": buildinfo.Version,
				"commit":  buildinfo.Commit,
				"date":    buildinfo.Date,
			}
			if app.Format == output.Table && len(app.UserColumns) == 0 && !app.WideMode {
				printHumanVersion(app)
				return nil
			}
			if app.Format == output.Table {
				app.PrintJSON([]map[string]string{info})
				return nil
			}
			app.PrintJSON(info)
			return nil
		},
	}
}

func printHumanVersion(app *cliapp.Runtime) {
	_, _ = fmt.Fprintf(app.Stdout, "ikuai-cli %s\n", buildinfo.Version)
	if buildinfo.Commit != "" && buildinfo.Commit != "none" {
		_, _ = fmt.Fprintf(app.Stdout, "commit: %s\n", buildinfo.Commit)
	}
	if buildinfo.Date != "" && buildinfo.Date != "unknown" {
		_, _ = fmt.Fprintf(app.Stdout, "built: %s\n", buildinfo.Date)
	}
}
