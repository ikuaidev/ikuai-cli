package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/peterh/liner"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Enter interactive REPL (default when no subcommand given)",
	RunE:  runREPL,
}

var helpGroups = []struct{ name, desc string }{
	{"auth", "Authentication (login, set-url, set-token, clear, status)"},
	{"monitor", "System monitoring (system, cpu, memory, clients-online, client-protocols ...)"},
	{"network", "Network config (wan, lan, dhcp, dns, dnat, vlan ...)"},
	{"security", "Security rules (acl, mac, l7, url, domain-blacklist ...)"},
	{"objects", "Network objects (ip, mac, port, domain, time ...)"},
	{"qos", "QoS bandwidth control (ip, mac)"},
	{"routing", "Routing & traffic shunting (static, stream, protocols)"},
	{"vpn", "VPN (pptp, l2tp, openvpn, ikev2, ipsec, wireguard)"},
	{"users", "User management (accounts, online, kick, packages)"},
	{"log", "System logs (arp, auth, dhcp, system, web, wireless ...)"},
	{"system", "System config (get, schedules, remote-access, vrrp, backup, upgrade, web-admin)"},
	{"auth-server", "Web auth server (get)"},
	{"wireless", "Wireless control (blacklist, vlan, ac)"},
	{"advanced", "Advanced services (ftp, http, samba, snmpd)"},
	{"completion", "Generate shell completion scripts"},
	{"version", "Show build version info"},
	{"help", "Show this command list"},
	{"quit", "Exit REPL"},
}

// ANSI color helpers — disabled when NO_COLOR is set or stdout is not a TTY.
var useColor bool

func initColor() {
	useColor = isTerminalWriter(stdout) && os.Getenv("NO_COLOR") == ""
}

func ansi(code, s string) string {
	if !useColor {
		return s
	}
	return code + s + "\x1b[0m"
}

func cyan(s string) string   { return ansi("\x1b[1;36m", s) }
func green(s string) string  { return ansi("\x1b[32m", s) }
func yellow(s string) string { return ansi("\x1b[33m", s) }
func dim(s string) string    { return ansi("\x1b[2m", s) }

// bannerWidth is the inner content width between box borders.
const bannerWidth = 48

// bannerLine builds a line padded to bannerWidth visible chars inside "│ ... │".
// visibleLen is the on-screen character count of content (excluding ANSI escapes).
func bannerLine(w io.Writer, content string, visibleLen int) {
	pad := bannerWidth - visibleLen
	if pad < 0 {
		pad = 0
	}
	_, _ = fmt.Fprintf(w, "│ %s%s │\n", content, strings.Repeat(" ", pad))
}

func printBanner() {
	// Skip banner entirely when stdout is not a terminal.
	if !isTerminalWriter(stdout) && os.Getenv("IKUAI_FORCE_TTY") != "1" {
		return
	}
	initColor()
	w := stdout

	authed := app.Session != nil && app.Session.Token != ""
	border := strings.Repeat("─", bannerWidth+2)

	_, _ = fmt.Fprintf(w, "╭%s╮\n", border)
	brand := " ◆  " + cyan("ikuai-cli") + "  ·  iKuai Router CLI"
	bannerLine(w, brand, len(" ◆  ikuai-cli  ·  iKuai Router CLI"))
	if authed {
		url := trunc(app.Session.BaseURL, 34)
		routerLine := "    Router:  " + url
		bannerLine(w, routerLine, len(routerLine))

		src := app.CredSource
		statusText := "    Status:  " + green("✔") + " authenticated  " + dim("(via "+src+")")
		statusVisible := len("    Status:  ✔ authenticated  (via " + src + ")")
		bannerLine(w, statusText, statusVisible)
	} else {
		bannerLine(w, "", 0)
		warnLine := " " + yellow("⚠") + "  Not authenticated"
		bannerLine(w, warnLine, len(" ⚠  Not authenticated"))
		bannerLine(w, "", 0)
		getStarted := " " + dim("Get started:")
		bannerLine(w, getStarted, len(" Get started:"))
		urlLine := "   " + green("$") + " auth set-url <URL>"
		bannerLine(w, urlLine, len("   $ auth set-url <URL>"))
		tokLine := "   " + green("$") + " auth set-token <TOKEN>"
		bannerLine(w, tokLine, len("   $ auth set-token <TOKEN>"))
		bannerLine(w, "", 0)
	}
	helpLine := " " + dim("Type") + " help" + dim(",") + " quit"
	bannerLine(w, helpLine, len(" Type help, quit"))
	_, _ = fmt.Fprintf(w, "╰%s╯\n", border)
	_, _ = fmt.Fprintln(w)
}

func trunc(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
}

func printHelp() {
	fmt.Println()
	fmt.Println("  Commands:")
	for _, g := range helpGroups {
		fmt.Printf("    %-14s  %s\n", g.name, g.desc)
	}
	fmt.Println()
}

func runREPL(_ *cobra.Command, _ []string) error {
	// Re-load session so banner shows current state
	if s, err := app.LoadSession(); err == nil {
		app.Session = s
		if app.Session != nil && app.Session.BaseURL != "" {
			app.APIClient = app.NewClient(app.Session.BaseURL, app.Session.Token)
		} else {
			app.APIClient = nil
		}
	}

	printBanner()

	line := liner.NewLiner()
	defer func() {
		_ = line.Close()
	}()
	line.SetCtrlCAborts(true)

	// Tab completion: walk the cobra tree to support multi-level commands.
	line.SetWordCompleter(func(line string, pos int) (head string, completions []string, tail string) {
		pre := string([]rune(line)[:pos])
		tail = string([]rune(line)[pos:])
		tokens, partial := splitForCompletion(pre)

		// R2: no flag completion.
		if strings.HasPrefix(partial, "-") {
			return pre, nil, tail
		}

		head = pre[:len(pre)-len(partial)]
		completions = completeFromTree(rootCmd, tokens, partial)
		return
	})

	for {
		input, err := line.Prompt("◆ ikuai ❯ ")
		if err != nil {
			// Ctrl+D or Ctrl+C
			fmt.Println()
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		line.AppendHistory(input)

		lower := strings.ToLower(input)
		if lower == "quit" || lower == "exit" || lower == "q" {
			break
		}
		if lower == "help" || lower == "?" || lower == "h" {
			printHelp()
			continue
		}

		// Parse into args (handles quoted strings and escapes)
		args, parseErr := splitArgs(input)
		if parseErr != nil {
			fmt.Fprintln(os.Stderr, "✗", parseErr)
			continue
		}
		if len(args) == 0 {
			continue
		}

		// Reset ALL flags before each REPL command to prevent stale values
		// from leaking across iterations. MergeDataWithFlags checks f.Changed,
		// so stale Changed=true bits would silently inject previous values.
		resetAllFlags(rootCmd)

		// Execute via rootCmd
		rootCmd.SetArgs(args)
		if err := rootCmd.Execute(); err != nil {
			fmt.Fprintln(os.Stderr, "✗", err)
		}
		rootCmd.SetArgs(nil)
	}

	fmt.Println("Goodbye!")
	return nil
}

// splitForCompletion splits the text before the cursor into fully-typed
// tokens and a trailing partial word. If the text ends with whitespace,
// partial is empty (user pressed Tab after a space — wants all children).
func splitForCompletion(pre string) (tokens []string, partial string) {
	fields := strings.Fields(pre)
	if len(fields) == 0 {
		return nil, ""
	}
	// If the raw text ends with a space, all fields are complete tokens.
	if pre[len(pre)-1] == ' ' {
		return fields, ""
	}
	// Otherwise the last field is the partial word being typed.
	return fields[:len(fields)-1], fields[len(fields)-1]
}

// completeFromTree walks the cobra command tree along tokens, then returns
// child command names matching the partial prefix.
func completeFromTree(root *cobra.Command, tokens []string, partial string) []string {
	node := root
	for _, tok := range tokens {
		found := false
		for _, child := range node.Commands() {
			name := strings.Fields(child.Use)[0]
			if name == tok {
				node = child
				found = true
				break
			}
		}
		if !found {
			return nil // unknown token — no completions
		}
	}
	var out []string
	for _, child := range node.Commands() {
		name := strings.Fields(child.Use)[0]
		if strings.HasPrefix(name, partial) {
			out = append(out, name)
		}
	}
	return out
}

// resetAllFlags traverses the command tree and resets every flag (persistent
// and local) to its default value with Changed=false. This prevents stale
// flag values from leaking between REPL iterations.
func resetAllFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	for _, child := range cmd.Commands() {
		resetAllFlags(child)
	}
}

// splitArgs splits a command line respecting single and double quotes.
// Returns an error for unterminated quotes or trailing backslash.
func splitArgs(s string) ([]string, error) {
	var args []string
	var cur strings.Builder
	inSingle, inDouble := false, false

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\\' && !inSingle:
			if i+1 >= len(s) {
				return nil, fmt.Errorf("trailing backslash in input")
			}
			cur.WriteByte(s[i+1])
			i++
		case c == '\'' && !inDouble:
			inSingle = !inSingle
		case c == '"' && !inSingle:
			inDouble = !inDouble
		case c == ' ' && !inSingle && !inDouble:
			if cur.Len() > 0 {
				args = append(args, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteByte(c)
		}
	}
	if inSingle || inDouble {
		return nil, fmt.Errorf("unterminated quote in input")
	}
	if cur.Len() > 0 {
		args = append(args, cur.String())
	}
	return args, nil
}
