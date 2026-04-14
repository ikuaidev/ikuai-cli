package system

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/ikuaidev/ikuai-cli/internal/session"
	ikuaissh "github.com/ikuaidev/ikuai-cli/internal/ssh"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func New(app *cliapp.Runtime) *cobra.Command {
	systemCmd := &cobra.Command{
		Use:   "system",
		Short: "System config",
	}

	systemGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get system config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "hostname", "language", "time_zone", "ntp_sync_cycle", "switch_ntp", "fast_nat", "lan_nat"}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/system/basic/config", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	systemSetFieldMap := map[string]string{
		"hostname":  "hostname",
		"language":  "language",
		"time-zone": "time_zone",
	}
	systemSetCmd := &cobra.Command{
		Use:   "set",
		Short: "Update system config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, systemSetFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/system/basic/config", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	systemNTPSyncCmd := &cobra.Command{
		Use:   "ntp-sync",
		Short: "Sync NTP",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Post(cliapp.APIBase+"/system/basic/ntp:sync", map[string]string{})
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	systemSchedulesCmd := &cobra.Command{Use: "schedules", Short: "Reboot schedules"}

	systemSchedulesListCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List schedules",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/system/reboot-schedules",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	schedulesFieldMap := map[string]string{
		"name":       "tagname",
		"event":      "event",
		"time":       "time",
		"strategy":   "strategy",
		"cycle-time": "cycle_time",
		"comment":    "comment",
		"enabled":    "enabled",
	}
	schedulesDefaults := map[string]interface{}{
		"enabled": "yes",
		"event":   "reboot",
		"comment": "",
	}
	systemSchedulesCreateCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create schedule",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "name", "time", "strategy", "cycle-time"); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, schedulesFieldMap)
			if err != nil {
				return err
			}
			for k, v := range schedulesDefaults {
				if _, exists := body[k]; !exists {
					body[k] = v
				}
			}
			raw, err := app.APIClient.Post(cliapp.APIBase+"/system/reboot-schedules", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	systemSchedulesUpdateCmd := &cobra.Command{
		Use:   "update ID",
		Short: "Update schedule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, schedulesFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/system/reboot-schedules/"+args[0], body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	toggleFieldMap := map[string]string{"enabled": "enabled"}
	systemSchedulesToggleCmd := &cobra.Command{
		Use:   "toggle ID",
		Short: "Toggle schedule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "enabled"); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, toggleFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Patch(cliapp.APIBase+"/system/reboot-schedules/"+args[0], body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	systemSchedulesDeleteCmd := &cobra.Command{
		Use:     "delete ID",
		Aliases: []string{"rm"},
		Short:   "Delete schedule",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			yes, _ := cmd.Flags().GetBool("yes")
			if err := cliapp.ConfirmDelete(app.Stdout, app.Stderr, "schedule", args[0], yes); err != nil {
				return err
			}
			raw, err := app.APIClient.Delete(cliapp.APIBase + "/system/reboot-schedules/" + args[0])
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	systemSchedulesDeleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	systemRemoteAccessCmd := &cobra.Command{Use: "remote-access", Short: "Remote access config"}

	systemRemoteAccessGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get remote access config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "open_sshd", "sshd_port", "open_telnetd", "open_wanweb", "http_port", "https_port", "force_https"}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/system/remote-access", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	remoteAccessFieldMap := map[string]string{
		"telnet":      "open_telnetd",
		"wan-web":     "open_wanweb",
		"ssh":         "open_sshd",
		"ssh-port":    "sshd_port",
		"http-port":   "http_port",
		"https-port":  "https_port",
		"force-https": "force_https",
	}
	systemRemoteAccessSetCmd := &cobra.Command{
		Use:   "set",
		Short: "Update remote access config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, remoteAccessFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/system/remote-access", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	systemVRRPCmd := &cobra.Command{Use: "vrrp", Short: "VRRP config"}

	systemVRRPGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get VRRP config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"enabled", "type", "prio", "gateway", "remote_addr", "interfaces", "virtual_ips", "auto_sync"}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/system/vrrp/config", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	vrrpFieldMap := map[string]string{
		"type":        "type",
		"priority":    "prio",
		"gateway":     "gateway",
		"remote-addr": "remote_addr",
		"enabled":     "enabled",
	}
	systemVRRPSetCmd := &cobra.Command{
		Use:   "set",
		Short: "Update VRRP config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, vrrpFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/system/vrrp/config", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	systemVRRPStartCmd := &cobra.Command{
		Use:   "start",
		Short: "Start VRRP",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Post(cliapp.APIBase+"/system/vrrp:start", map[string]string{})
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	systemVRRPStopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop VRRP",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Post(cliapp.APIBase+"/system/vrrp:stop", map[string]string{})
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	systemALGCmd := &cobra.Command{Use: "alg", Short: "ALG config"}

	systemALGGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get ALG config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/system/alg", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	algFieldMap := map[string]string{
		"ftp":        "support_ftp",
		"tftp":       "support_tftp",
		"sip":        "support_sip",
		"ftp-ports":  "ftp_ports",
		"sip-ports":  "sip_ports",
		"tftp-ports": "tftp_ports",
	}
	systemALGSetCmd := &cobra.Command{
		Use:   "set",
		Short: "Update ALG config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, algFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/system/alg", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	systemKernelCmd := &cobra.Command{Use: "kernel", Short: "Kernel params"}

	systemKernelGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get kernel params",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "bbr", "syn_recv_timeout", "established_timeout", "close_timeout", "fin_wait_timeout", "udp_timeout", "icmp_timeout"}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/system/kernel-params", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	kernelFieldMap := map[string]string{
		"bbr":                 "bbr",
		"syn-recv-timeout":    "syn_recv_timeout",
		"established-timeout": "established_timeout",
	}
	systemKernelSetCmd := &cobra.Command{
		Use:   "set",
		Short: "Update kernel params",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, kernelFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/system/kernel-params", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	systemCPUFreqCmd := &cobra.Command{Use: "cpufreq", Short: "CPU frequency config"}

	systemCPUFreqGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get CPU freq config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/system/cpufreq", nil)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	cpufreqFieldMap := map[string]string{
		"mode":  "mode",
		"turbo": "turbo",
	}
	systemCPUFreqSetCmd := &cobra.Command{
		Use:   "set",
		Short: "Update CPU freq",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, cpufreqFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/system/cpufreq", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	systemCPUFreqModeSetCmd := &cobra.Command{
		Use:   "mode-set",
		Short: "Set CPU freq mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, cpufreqFieldMap)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+"/system/cpufreq/mode", body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}

	systemWebPasswdCmd := &cobra.Command{Use: "web-passwd", Short: "Web password management"}

	systemWebPasswdResetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset web password via SSH",
		Long: `Connects to the router via SSH, selects option 7 (恢复WEB管理密码),
confirms, and returns the new credentials.

SSH credentials are resolved in order:
  1. --ssh-user / --ssh-password / --ssh-port flags
  2. Saved config (~/.ikuai-cli/config.json)
  3. Interactive prompt`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireURL(); err != nil {
				return err
			}

			skipConfirm, _ := cmd.Flags().GetBool("yes")
			saveFlag, _ := cmd.Flags().GetBool("save")

			host, err := hostFromURL(app.Session.BaseURL)
			if err != nil {
				return fmt.Errorf("cannot parse host from URL %q: %w", app.Session.BaseURL, err)
			}

			sshUser, _ := cmd.Flags().GetString("ssh-user")
			sshPassword, _ := cmd.Flags().GetString("ssh-password")
			sshPort, _ := cmd.Flags().GetInt("ssh-port")

			if sshUser == "" && app.Session.SSHUser != "" {
				sshUser = app.Session.SSHUser
			}
			if sshPassword == "" && app.Session.SSHPassword != "" {
				sshPassword = app.Session.SSHPassword
			}
			if sshPort == 22 && app.Session.SSHPort != 0 {
				sshPort = app.Session.SSHPort
			}

			if sshUser == "" || sshPassword == "" {
				if !term.IsTerminal(int(os.Stdin.Fd())) {
					return &cliapp.ValidationError{Message: "--ssh-user and --ssh-password are required in non-interactive mode"}
				}
				if sshUser == "" {
					sshUser = prompt("SSH username: ")
				}
				if sshPassword == "" {
					sshPassword = promptPassword("SSH password: ")
				}
			}

			if !skipConfirm {
				fmt.Printf("Reset Web admin password on %s via SSH as %s? [y/N] ", host, sshUser)
				reader := bufio.NewReader(os.Stdin)
				ans, _ := reader.ReadString('\n')
				if !strings.HasPrefix(strings.TrimSpace(strings.ToLower(ans)), "y") {
					fmt.Println("Aborted.")
					return nil
				}
			}

			username, passwd, err := ikuaissh.ResetWebPasswd(host, sshUser, sshPassword, sshPort)
			if err != nil {
				return err
			}

			if saveFlag {
				if err := session.SaveSSHCreds(sshUser, sshPassword, sshPort); err != nil {
					fmt.Fprintln(os.Stderr, "Warning: could not save SSH creds:", err)
				}
			}

			app.PrintJSON(map[string]interface{}{
				"code":    0,
				"message": "Web management password reset successful",
				"results": map[string]string{
					"username": username,
					"password": passwd,
				},
			})
			return nil
		},
	}

	systemCmd.AddCommand(systemGetCmd, systemSetCmd, systemNTPSyncCmd)

	systemCmd.AddCommand(systemSchedulesCmd)
	systemSchedulesCmd.AddCommand(
		systemSchedulesListCmd,
		systemSchedulesCreateCmd,
		systemSchedulesUpdateCmd,
		systemSchedulesToggleCmd,
		systemSchedulesDeleteCmd,
	)
	cliapp.AddListFlags(systemSchedulesListCmd)
	for _, c := range []*cobra.Command{
		systemSchedulesCreateCmd, systemSchedulesUpdateCmd, systemSchedulesToggleCmd,
	} {
		c.Flags().String("data", "{}", "JSON body")
	}
	cliapp.AddEnabledFlag(systemSchedulesToggleCmd)
	for _, c := range []*cobra.Command{systemSchedulesCreateCmd, systemSchedulesUpdateCmd} {
		c.Flags().String("name", "", "Schedule name (tagname)")
		c.Flags().String("event", "", "Event type (reboot/poweroff)")
		c.Flags().String("time", "", "Schedule time (HH:MM)")
		c.Flags().String("strategy", "", "Schedule strategy (week/month/day/once)")
		c.Flags().String("cycle-time", "", "Cycle time (e.g. 7 for weekly)")
		c.Flags().String("comment", "", "Comment")
		cliapp.AddEnabledFlag(c)
	}

	systemCmd.AddCommand(
		systemRemoteAccessCmd,
		systemVRRPCmd,
		systemALGCmd,
		systemKernelCmd,
		systemCPUFreqCmd,
		systemWebPasswdCmd,
	)

	systemRemoteAccessCmd.AddCommand(systemRemoteAccessGetCmd, systemRemoteAccessSetCmd)
	systemVRRPCmd.AddCommand(systemVRRPGetCmd, systemVRRPSetCmd, systemVRRPStartCmd, systemVRRPStopCmd)
	systemALGCmd.AddCommand(systemALGGetCmd, systemALGSetCmd)
	systemKernelCmd.AddCommand(systemKernelGetCmd, systemKernelSetCmd)
	systemCPUFreqCmd.AddCommand(systemCPUFreqGetCmd, systemCPUFreqSetCmd, systemCPUFreqModeSetCmd)
	systemWebPasswdCmd.AddCommand(systemWebPasswdResetCmd)

	for _, c := range []*cobra.Command{
		systemSetCmd, systemRemoteAccessSetCmd, systemVRRPSetCmd,
		systemALGSetCmd, systemKernelSetCmd, systemCPUFreqSetCmd, systemCPUFreqModeSetCmd,
	} {
		c.Flags().String("data", "{}", "JSON body")
	}

	// system set semantic flags
	systemSetCmd.Flags().String("hostname", "", "Router hostname")
	systemSetCmd.Flags().String("language", "", "UI language")
	systemSetCmd.Flags().String("time-zone", "", "Time zone")

	// remote-access set semantic flags
	systemRemoteAccessSetCmd.Flags().String("telnet", "", "Open Telnet (yes/no)")
	systemRemoteAccessSetCmd.Flags().String("wan-web", "", "Open WAN web access (yes/no)")
	systemRemoteAccessSetCmd.Flags().String("ssh", "", "Open SSH (yes/no)")
	systemRemoteAccessSetCmd.Flags().String("ssh-port", "", "SSH port")
	systemRemoteAccessSetCmd.Flags().String("http-port", "", "HTTP port")
	systemRemoteAccessSetCmd.Flags().String("https-port", "", "HTTPS port")
	systemRemoteAccessSetCmd.Flags().String("force-https", "", "Force HTTPS (yes/no)")

	// vrrp set semantic flags
	systemVRRPSetCmd.Flags().String("type", "", "VRRP type")
	systemVRRPSetCmd.Flags().String("priority", "", "VRRP priority")
	systemVRRPSetCmd.Flags().String("gateway", "", "Gateway address")
	systemVRRPSetCmd.Flags().String("remote-addr", "", "Remote address")
	cliapp.AddEnabledFlag(systemVRRPSetCmd)

	// alg set semantic flags
	systemALGSetCmd.Flags().String("ftp", "", "FTP ALG support (yes/no)")
	systemALGSetCmd.Flags().String("tftp", "", "TFTP ALG support (yes/no)")
	systemALGSetCmd.Flags().String("sip", "", "SIP ALG support (yes/no)")
	systemALGSetCmd.Flags().String("ftp-ports", "", "FTP ports")
	systemALGSetCmd.Flags().String("sip-ports", "", "SIP ports")
	systemALGSetCmd.Flags().String("tftp-ports", "", "TFTP ports")

	// kernel set semantic flags
	systemKernelSetCmd.Flags().String("bbr", "", "Enable BBR (yes/no)")
	systemKernelSetCmd.Flags().String("syn-recv-timeout", "", "SYN receive timeout")
	systemKernelSetCmd.Flags().String("established-timeout", "", "Established connection timeout")

	// cpufreq set semantic flags
	systemCPUFreqSetCmd.Flags().String("mode", "", "CPU frequency mode")
	systemCPUFreqSetCmd.Flags().String("turbo", "", "Turbo boost (yes/no)")

	// cpufreq mode-set semantic flags
	systemCPUFreqModeSetCmd.Flags().String("mode", "", "CPU frequency mode")
	systemCPUFreqModeSetCmd.Flags().String("turbo", "", "Turbo boost (yes/no)")

	systemWebPasswdResetCmd.Flags().String("ssh-user", "", "SSH username")
	systemWebPasswdResetCmd.Flags().String("ssh-password", "", "SSH password")
	systemWebPasswdResetCmd.Flags().Int("ssh-port", 22, "SSH port")
	systemWebPasswdResetCmd.Flags().Bool("save", false, "Save SSH credentials to session file")
	systemWebPasswdResetCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	return systemCmd
}

func hostFromURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL, err
	}
	if u.Hostname() != "" {
		return u.Hostname(), nil
	}
	return rawURL, nil
}

func prompt(label string) string {
	fmt.Print(label)
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	return strings.TrimSpace(s)
}

func promptPassword(label string) string {
	fmt.Print(label)
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return ""
	}
	return string(b)
}
