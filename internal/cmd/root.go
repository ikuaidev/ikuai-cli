// Package cmd contains all cobra commands for ikuai-cli.
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/api"
	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	advancedcmd "github.com/ikuaidev/ikuai-cli/internal/cmd/advanced"
	authcmd "github.com/ikuaidev/ikuai-cli/internal/cmd/auth"
	authservercmd "github.com/ikuaidev/ikuai-cli/internal/cmd/authserver"
	completioncmd "github.com/ikuaidev/ikuai-cli/internal/cmd/completion"
	logcmd "github.com/ikuaidev/ikuai-cli/internal/cmd/log"
	monitorcmd "github.com/ikuaidev/ikuai-cli/internal/cmd/monitor"
	networkcmd "github.com/ikuaidev/ikuai-cli/internal/cmd/network"
	objectscmd "github.com/ikuaidev/ikuai-cli/internal/cmd/objects"
	qoscmd "github.com/ikuaidev/ikuai-cli/internal/cmd/qos"
	routingcmd "github.com/ikuaidev/ikuai-cli/internal/cmd/routing"
	securitycmd "github.com/ikuaidev/ikuai-cli/internal/cmd/security"
	systemcmd "github.com/ikuaidev/ikuai-cli/internal/cmd/system"
	userscmd "github.com/ikuaidev/ikuai-cli/internal/cmd/users"
	versioncmd "github.com/ikuaidev/ikuai-cli/internal/cmd/version"
	vpncmd "github.com/ikuaidev/ikuai-cli/internal/cmd/vpn"
	wirelesscmd "github.com/ikuaidev/ikuai-cli/internal/cmd/wireless"
	"github.com/ikuaidev/ikuai-cli/internal/output"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	formatStr  string
	rawOutput  bool
	humanTime  bool
	dryRun     bool
	wideOutput bool
	columnsStr string
	stdout     io.Writer = os.Stdout
	stderr     io.Writer = os.Stderr
	app                  = cliapp.New(stdout, stderr)
)

var rootCmd = &cobra.Command{
	Use:          "ikuai-cli",
	Short:        "iKuai router local API v4.0 CLI",
	Long:         `CLI for managing an iKuai router via its local REST API (v4.0). Default output is table; use --format json/yaml for machine-friendly output.`,
	SilenceUsage: true,
	// RunE is set in init() to avoid initialization cycle with repl.go
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		app.Stdout = stdout
		app.Stderr = stderr

		// Validate --format.
		f, err := output.FormatFromString(formatStr)
		if err != nil {
			return err
		}

		// --raw and --format are mutually exclusive.
		if rawOutput && cmd.Flags().Changed("format") {
			return fmt.Errorf("--raw and --format are mutually exclusive")
		}

		// --wide and --columns are mutually exclusive.
		if wideOutput && cmd.Flags().Changed("columns") {
			return fmt.Errorf("--wide and --columns are mutually exclusive")
		}

		app.Format = f
		app.RawMode = rawOutput
		app.HumanTime = humanTime
		app.DryRun = dryRun
		app.WideMode = wideOutput
		app.DefaultColumns = nil // Reset per-command defaults to prevent leaking in REPL mode.
		app.UserColumns = nil

		// Parse --columns (only when explicitly passed).
		if cmd.Flags().Changed("columns") && columnsStr != "" {
			parts := strings.Split(columnsStr, ",")
			cols := parts[:0]
			for _, c := range parts {
				c = strings.TrimSpace(c)
				if c != "" {
					cols = append(cols, c)
				}
			}
			app.UserColumns = cols
		}

		// TTY detection: auto-switch to JSON when stdout is not a terminal.
		// Check the actual output writer (not hardcoded os.Stdout) so that
		// programmatic callers that redirect stdout also get correct behavior.
		isTTY := os.Getenv("IKUAI_FORCE_TTY") == "1" || isTerminalWriter(stdout)
		if !cmd.Flags().Changed("format") && !rawOutput {
			if !isTTY {
				app.Format = output.JSON
			}
		}

		// Terminal width for auto-fit (only when TTY and not --wide).
		if isTTY && !wideOutput {
			if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
				app.TermWidth = w
			}
		}

		return app.SyncSession()
	},
}

// Exit codes for structured error reporting.
const (
	ExitOK         = 0
	ExitGeneral    = 1
	ExitValidation = 2
	ExitAuth       = 3
	ExitNetwork    = 4
	ExitAPI        = 5
)

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		code := classifyError(err)
		// Non-TTY: emit JSON error envelope to stderr for machine consumers.
		if !isTerminalWriter(stderr) {
			envelope := map[string]interface{}{
				"ok":    false,
				"error": buildErrorPayload(err, code),
			}
			enc := json.NewEncoder(stderr)
			enc.SetEscapeHTML(false)
			_ = enc.Encode(envelope)
		}
		os.Exit(code)
	}
}

// classifyError maps an error to an exit code.
func classifyError(err error) int {
	var authErr *cliapp.AuthError
	if errors.As(err, &authErr) {
		return ExitAuth
	}
	var valErr *cliapp.ValidationError
	if errors.As(err, &valErr) {
		return ExitValidation
	}
	var netErr *api.NetworkError
	if errors.As(err, &netErr) {
		return ExitNetwork
	}
	var apiErr *api.APIError
	if errors.As(err, &apiErr) {
		// Auth-related API codes → exit 3
		if apiErr.Code == 3007 || apiErr.Code == 1008 {
			return ExitAuth
		}
		return ExitAPI
	}
	return ExitGeneral
}

// buildErrorPayload creates a structured error map for JSON output.
func buildErrorPayload(err error, code int) map[string]interface{} {
	payload := map[string]interface{}{
		"message":   err.Error(),
		"exit_code": code,
	}
	var apiErr *api.APIError
	if errors.As(err, &apiErr) {
		payload["type"] = "api_error"
		payload["api_code"] = apiErr.Code
		if len(apiErr.Details) > 0 {
			payload["details"] = apiErr.Details
		}
	} else {
		var authErr *cliapp.AuthError
		var valErr *cliapp.ValidationError
		var netErr *api.NetworkError
		switch {
		case errors.As(err, &authErr):
			payload["type"] = "auth_error"
		case errors.As(err, &valErr):
			payload["type"] = "validation_error"
		case errors.As(err, &netErr):
			payload["type"] = "network_error"
		default:
			payload["type"] = "general_error"
		}
	}
	return payload
}

func init() {
	// Set REPL as the default action when no subcommand is given.
	// Done here (not in var declaration) to avoid initialization cycle.
	rootCmd.RunE = runREPL

	rootCmd.AddCommand(replCmd)
	rootCmd.AddCommand(completioncmd.New(app))
	rootCmd.AddCommand(authcmd.New(app))
	rootCmd.AddCommand(monitorcmd.New(app))
	rootCmd.AddCommand(networkcmd.New(app))
	rootCmd.AddCommand(securitycmd.New(app))
	rootCmd.AddCommand(objectscmd.New(app))
	rootCmd.AddCommand(qoscmd.New(app))
	rootCmd.AddCommand(routingcmd.New(app))
	rootCmd.AddCommand(vpncmd.New(app))
	rootCmd.AddCommand(userscmd.New(app))
	rootCmd.AddCommand(logcmd.New(app))
	rootCmd.AddCommand(systemcmd.New(app))
	rootCmd.AddCommand(authservercmd.New(app))
	rootCmd.AddCommand(wirelesscmd.New(app))
	rootCmd.AddCommand(advancedcmd.New(app))
	rootCmd.AddCommand(versioncmd.New(app))
	rootCmd.PersistentFlags().StringVarP(&formatStr, "format", "f", "table", "Output format: table, json, yaml")
	rootCmd.PersistentFlags().BoolVar(&rawOutput, "raw", false, "Output full API envelope (debug)")
	rootCmd.PersistentFlags().BoolVar(&humanTime, "human-time", false, "Convert timestamp columns to human-readable local time")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Preview API request without executing")
	rootCmd.PersistentFlags().BoolVarP(&wideOutput, "wide", "w", false, "Show all table columns")
	rootCmd.PersistentFlags().StringVar(&columnsStr, "columns", "", "Comma-separated list of columns to display")

	// Enum completion for global flags.
	_ = rootCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"table", "json", "yaml"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Reformat cobra's required-flag error to match project convention:
	// "required flag(s) "x" not set" → "missing required flag: --x"
	rootCmd.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		msg := err.Error()
		const prefix = `required flag(s) `
		const suffix = ` not set`
		if strings.HasPrefix(msg, prefix) && strings.HasSuffix(msg, suffix) {
			inner := strings.TrimPrefix(msg, prefix)
			inner = strings.TrimSuffix(inner, suffix)
			// inner looks like: `"name"` or `"name1", "name2"`
			parts := strings.Split(inner, ", ")
			flags := make([]string, 0, len(parts))
			for _, p := range parts {
				flags = append(flags, "--"+strings.Trim(p, `"`))
			}
			if len(flags) == 1 {
				return &cliapp.ValidationError{Message: "missing required flag: " + flags[0]}
			}
			return &cliapp.ValidationError{Message: "missing required flags: " + strings.Join(flags, ", ")}
		}
		return &cliapp.ValidationError{Message: msg}
	})
}

// isTerminalWriter checks whether w is connected to a terminal.
// If w is an *os.File, its file descriptor is checked; otherwise non-terminal.
func isTerminalWriter(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return term.IsTerminal(int(f.Fd()))
	}
	return false
}
