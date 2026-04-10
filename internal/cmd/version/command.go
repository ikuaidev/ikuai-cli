package version

import (
	"github.com/ikuaidev/ikuai-cli/internal/buildinfo"
	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

func New(app *cliapp.Runtime) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version",
		RunE: func(cmd *cobra.Command, args []string) error {
			app.PrintJSON(map[string]string{
				"name":    "ikuai-cli",
				"version": buildinfo.Version,
				"commit":  buildinfo.Commit,
				"date":    buildinfo.Date,
			})
			return nil
		},
	}
}
