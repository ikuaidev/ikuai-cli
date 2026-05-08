package system

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

var (
	systemBasicFields = []string{
		"hostname", "language", "time_zone", "time_zone_full",
		"switch_nat", "switch_ntp", "switch_ntpd", "switch_ntpserver",
		"ntpserver_list", "ntp_sync_cycle", "link_mode", "lan_nat",
		"backport", "listenport", "fast_nat",
	}
	systemSetFieldMap = map[string]string{
		"hostname":         "hostname",
		"language":         "language",
		"time-zone":        "time_zone",
		"time-zone-full":   "time_zone_full",
		"switch-nat":       "switch_nat",
		"switch-ntp":       "switch_ntp",
		"switch-ntpd":      "switch_ntpd",
		"switch-ntpserver": "switch_ntpserver",
		"ntpserver-list":   "ntpserver_list",
		"ntp-sync-cycle":   "ntp_sync_cycle",
		"link-mode":        "link_mode",
		"lan-nat":          "lan_nat",
		"backport":         "backport",
		"listenport":       "listenport",
		"fast-nat":         "fast_nat",
	}

	schedulesFields   = []string{"enabled", "tagname", "event", "strategy", "cycle_time", "time", "comment"}
	schedulesFieldMap = map[string]string{
		"name":       "tagname",
		"event":      "event",
		"time":       "time",
		"strategy":   "strategy",
		"cycle-time": "cycle_time",
		"comment":    "comment",
		"enabled":    "enabled",
	}
	schedulesDefaults = map[string]interface{}{
		"enabled": "yes",
		"event":   "reboot",
		"comment": "",
	}

	remoteAccessFields = []string{
		"open_telnetd", "open_wanweb", "open_sshd", "sshd_port",
		"sshd_passwd", "http_port", "https_port", "force_https",
	}
	remoteAccessFieldMap = map[string]string{
		"telnet":      "open_telnetd",
		"wan-web":     "open_wanweb",
		"ssh":         "open_sshd",
		"ssh-port":    "sshd_port",
		"http-port":   "http_port",
		"https-port":  "https_port",
		"force-https": "force_https",
	}

	vrrpFields = []string{
		"enabled", "type", "prio", "method", "domain", "dns", "gateway",
		"interval", "ifnames", "auto_sync", "single_line", "ignore_wanstatus",
		"interfaces", "virtual_ips", "ht_iface", "remote_addr",
	}
	vrrpDefaults = map[string]interface{}{
		"method":   2,
		"domain":   "www.baidu.com",
		"dns":      "114.114.114.114",
		"interval": 3,
	}
	vrrpFieldMap = map[string]string{
		"type":             "type",
		"priority":         "prio",
		"method":           "method",
		"domain":           "domain",
		"dns":              "dns",
		"gateway":          "gateway",
		"interval":         "interval",
		"ifnames":          "ifnames",
		"auto-sync":        "auto_sync",
		"single-line":      "single_line",
		"ignore-wanstatus": "ignore_wanstatus",
		"interfaces":       "interfaces",
		"virtual-ips":      "virtual_ips",
		"heartbeat-iface":  "ht_iface",
		"remote-addr":      "remote_addr",
		"enabled":          "enabled",
	}

	algFields   = []string{"support_ftp", "support_tftp", "support_sip", "support_h323", "ftp_ports", "sip_ports", "tftp_ports"}
	algFieldMap = map[string]string{
		"ftp":        "support_ftp",
		"tftp":       "support_tftp",
		"sip":        "support_sip",
		"h323":       "support_h323",
		"ftp-ports":  "ftp_ports",
		"sip-ports":  "sip_ports",
		"tftp-ports": "tftp_ports",
	}

	kernelFields = []string{
		"bbr", "syn_recv_timeout", "syn_send_timeout", "established_timeout",
		"fin_wait_timeout", "last_ack_timeout", "close_wait_timeout",
		"time_wait_timeout", "close_timeout", "udp_timeout",
		"udp_stream_timeout", "icmp_timeout",
	}
	kernelFieldMap = map[string]string{
		"bbr":                 "bbr",
		"syn-recv-timeout":    "syn_recv_timeout",
		"syn-send-timeout":    "syn_send_timeout",
		"established-timeout": "established_timeout",
		"fin-wait-timeout":    "fin_wait_timeout",
		"last-ack-timeout":    "last_ack_timeout",
		"close-wait-timeout":  "close_wait_timeout",
		"time-wait-timeout":   "time_wait_timeout",
		"close-timeout":       "close_timeout",
		"udp-timeout":         "udp_timeout",
		"udp-stream-timeout":  "udp_stream_timeout",
		"icmp-timeout":        "icmp_timeout",
	}

	cpufreqModeFields = []string{"mode", "turbo"}
	cpufreqFieldMap   = map[string]string{"mode": "mode", "turbo": "turbo"}

	backupAutoFields = []string{"id", "enabled", "strategy", "time", "cycle_time", "valid_days"}
	backupAutoMap    = map[string]string{
		"id":         "id",
		"enabled":    "enabled",
		"strategy":   "strategy",
		"time":       "time",
		"cycle-time": "cycle_time",
		"valid-days": "valid_days",
	}

	webAdminGroupFields = []string{"group_name", "ip_addr", "perm_config"}
	webAdminGroupMap    = map[string]string{
		"name":        "group_name",
		"ip-addr":     "ip_addr",
		"perm-config": "perm_config",
	}

	webAdminAccountFields = []string{"username", "passwd", "enabled", "group_id", "force", "interval", "sesstimeout", "comment"}
	webAdminAccountMap    = map[string]string{
		"username":        "username",
		"passwd-md5":      "passwd",
		"enabled":         "enabled",
		"group-id":        "group_id",
		"force":           "force",
		"interval":        "interval",
		"session-timeout": "sesstimeout",
		"comment":         "comment",
	}
	webAdminAccountDefaults = map[string]interface{}{
		"enabled":     "yes",
		"force":       0,
		"interval":    30,
		"sesstimeout": 120,
		"comment":     "",
	}
)

func New(app *cliapp.Runtime) *cobra.Command {
	systemCmd := &cobra.Command{
		Use:   "system",
		Short: "System config",
	}

	systemCmd.AddCommand(
		systemGetCmd(app),
		systemSetCmd(app),
		systemNTPSyncCmd(app),
		systemSchedulesCmd(app),
		systemRemoteAccessCmd(app),
		systemVRRPCmd(app),
		systemALGCmd(app),
		systemKernelCmd(app),
		systemCPUFreqCmd(app),
		systemDisksCmd(app),
		systemFilesCmd(app),
		systemBackupCmd(app),
		systemBackupAutoCmd(app),
		systemUpgradeCmd(app),
		systemWebAdminCmd(app),
	)
	return systemCmd
}

func systemGetCmd(app *cliapp.Runtime) *cobra.Command {
	return readCommand(app, "get", "Get system config", "/system/basic/config",
		[]string{"id", "hostname", "language", "time_zone", "ntp_sync_cycle", "switch_ntp", "fast_nat", "lan_nat"}, nil)
}

func systemSetCmd(app *cliapp.Runtime) *cobra.Command {
	cmd := fullUpdateCommand(app, "set", "Update system config", "/system/basic/config", "/system/basic/config", systemSetFieldMap, systemBasicFields, nil)
	cmd.Flags().String("hostname", "", "Router hostname")
	cmd.Flags().String("language", "", "UI language ID")
	cmd.Flags().String("time-zone", "", "Time zone offset")
	cmd.Flags().String("time-zone-full", "", "Full time zone code")
	cmd.Flags().String("switch-nat", "", "NAT mode")
	cmd.Flags().String("switch-ntp", "", "Auto NTP sync switch")
	cmd.Flags().String("switch-ntpd", "", "NTP daemon switch")
	cmd.Flags().String("switch-ntpserver", "", "NTP server switch")
	cmd.Flags().String("ntpserver-list", "", "NTP server list")
	cmd.Flags().String("ntp-sync-cycle", "", "NTP sync cycle in minutes")
	cmd.Flags().String("link-mode", "", "Link mode")
	cmd.Flags().String("lan-nat", "", "LAN NAT switch")
	cmd.Flags().String("backport", "", "Backport interface")
	cmd.Flags().String("listenport", "", "Listen interface")
	cmd.Flags().String("fast-nat", "", "Fast NAT mode")
	return cmd
}

func systemNTPSyncCmd(app *cliapp.Runtime) *cobra.Command {
	return postActionCommand(app, "ntp-sync", "Sync NTP", "/system/basic/ntp:sync", map[string]string{})
}

func systemSchedulesCmd(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "schedules", Short: "Reboot schedules"}
	listCmd := readCommand(app, "list", "List schedules", "/system/reboot-schedules",
		[]string{"id", "tagname", "event", "strategy", "cycle_time", "time", "enabled"}, scheduleListParams)
	listCmd.Aliases = []string{"ls"}
	addPageAndOrderFlags(listCmd)

	getCmd := readCommand(app, "get ID", "Get a single schedule", "/system/reboot-schedules",
		[]string{"id", "tagname", "event", "strategy", "cycle_time", "time", "enabled"}, nil)

	createCmd := writeCommand(app, "create", "Create schedule", schedulesFieldMap, schedulesDefaults, func(body interface{}, _ string) (json.RawMessage, error) {
		return app.APIClient.Post(cliapp.APIBase+"/system/reboot-schedules", body)
	})
	createCmd.Aliases = []string{"new"}
	addScheduleFlags(createCmd)
	cliapp.MarkFlagsRequired(createCmd, "name", "time", "strategy", "cycle-time")
	wrapRequireFlags(createCmd, "name", "time", "strategy", "cycle-time")

	updateCmd := fullUpdateCommand(app, "update ID", "Update schedule", "/system/reboot-schedules", "/system/reboot-schedules", schedulesFieldMap, schedulesFields, nil)
	addScheduleFlags(updateCmd)

	toggleCmd := writeCommand(app, "toggle ID", "Toggle schedule", map[string]string{"enabled": "enabled"}, nil, func(body interface{}, id string) (json.RawMessage, error) {
		return app.APIClient.Patch(cliapp.APIBase+"/system/reboot-schedules/"+id, body)
	})
	cliapp.AddEnabledFlag(toggleCmd)
	cliapp.MarkFlagsRequired(toggleCmd, "enabled")
	wrapRequireFlags(toggleCmd, "enabled")

	deleteCmd := deleteIDCommand(app, "delete ID", "Delete schedule", "schedule", "/system/reboot-schedules")
	deleteCmd.Aliases = []string{"rm"}

	group.AddCommand(listCmd, getCmd, createCmd, updateCmd, toggleCmd, deleteCmd)
	return group
}

func systemRemoteAccessCmd(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "remote-access", Short: "Remote access config"}
	getCmd := readCommand(app, "get", "Get remote access config", "/system/remote-access",
		[]string{"id", "open_sshd", "sshd_port", "open_telnetd", "open_wanweb", "http_port", "https_port", "force_https"}, nil)
	setCmd := fullUpdateCommand(app, "set", "Update remote access config", "/system/remote-access", "/system/remote-access", remoteAccessFieldMap, remoteAccessFields, nil)
	setCmd.Flags().String("telnet", "", "Open Telnet (0/1)")
	setCmd.Flags().String("wan-web", "", "WAN web access mode")
	setCmd.Flags().String("ssh", "", "Open SSH (0/1)")
	setCmd.Flags().String("ssh-port", "", "SSH port")
	setCmd.Flags().String("http-port", "", "HTTP port")
	setCmd.Flags().String("https-port", "", "HTTPS port")
	setCmd.Flags().String("force-https", "", "Force HTTPS (0/1)")
	group.AddCommand(getCmd, setCmd)
	return group
}

func systemVRRPCmd(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "vrrp", Short: "VRRP config"}
	getCmd := readCommand(app, "get", "Get VRRP config", "/system/vrrp/config",
		[]string{"enabled", "type", "prio", "gateway", "remote_addr", "interfaces", "virtual_ips", "auto_sync"}, nil)
	setCmd := fullUpdateCommand(app, "set", "Update VRRP config", "/system/vrrp/config", "/system/vrrp/config", vrrpFieldMap, vrrpFields, vrrpDefaults)
	setCmd.Flags().String("type", "", "VRRP type")
	setCmd.Flags().String("priority", "", "VRRP priority")
	setCmd.Flags().String("method", "", "Detect method")
	setCmd.Flags().String("domain", "", "DNS detect domain")
	setCmd.Flags().String("dns", "", "DNS server")
	setCmd.Flags().String("gateway", "", "Gateway address")
	setCmd.Flags().String("interval", "", "Heartbeat interval")
	setCmd.Flags().String("ifnames", "", "WAN detect interfaces")
	setCmd.Flags().String("auto-sync", "", "Auto sync (0/1)")
	setCmd.Flags().String("single-line", "", "Single line mode (0/1)")
	setCmd.Flags().String("ignore-wanstatus", "", "Ignore WAN status (0/1)")
	setCmd.Flags().String("interfaces", "", "Transport interfaces")
	setCmd.Flags().String("virtual-ips", "", "Virtual IPs")
	setCmd.Flags().String("heartbeat-iface", "", "Heartbeat interface")
	setCmd.Flags().String("remote-addr", "", "Remote address")
	cliapp.AddEnabledFlag(setCmd)
	group.AddCommand(
		getCmd,
		setCmd,
		postActionCommand(app, "start", "Start VRRP", "/system/vrrp:start", map[string]string{}),
		postActionCommand(app, "stop", "Stop VRRP", "/system/vrrp:stop", map[string]string{}),
	)
	return group
}

func systemALGCmd(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "alg", Short: "ALG config"}
	getCmd := readCommand(app, "get", "Get ALG config", "/system/alg",
		[]string{"support_ftp", "support_tftp", "support_sip", "support_h323", "ftp_ports", "sip_ports", "tftp_ports"}, nil)
	setCmd := fullUpdateCommand(app, "set", "Update ALG config", "/system/alg", "/system/alg", algFieldMap, algFields, nil)
	setCmd.Flags().String("ftp", "", "FTP ALG support (0/1)")
	setCmd.Flags().String("tftp", "", "TFTP ALG support (0/1)")
	setCmd.Flags().String("sip", "", "SIP ALG support (0/1)")
	setCmd.Flags().String("h323", "", "H323 ALG support (0/1)")
	setCmd.Flags().String("ftp-ports", "", "FTP ports")
	setCmd.Flags().String("sip-ports", "", "SIP ports")
	setCmd.Flags().String("tftp-ports", "", "TFTP ports")
	group.AddCommand(getCmd, setCmd)
	return group
}

func systemKernelCmd(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "kernel", Short: "Kernel params"}
	getCmd := readCommand(app, "get", "Get kernel params", "/system/kernel-params",
		[]string{"id", "bbr", "syn_recv_timeout", "established_timeout", "close_timeout", "fin_wait_timeout", "udp_timeout", "icmp_timeout"}, nil)
	setCmd := fullUpdateCommand(app, "set", "Update kernel params", "/system/kernel-params", "/system/kernel-params", kernelFieldMap, kernelFields, nil)
	setCmd.Flags().String("bbr", "", "Enable BBR (0/1)")
	setCmd.Flags().String("syn-recv-timeout", "", "SYN_RECV timeout")
	setCmd.Flags().String("syn-send-timeout", "", "SYN_SEND timeout")
	setCmd.Flags().String("established-timeout", "", "Established timeout")
	setCmd.Flags().String("fin-wait-timeout", "", "FIN_WAIT timeout")
	setCmd.Flags().String("last-ack-timeout", "", "LAST_ACK timeout")
	setCmd.Flags().String("close-wait-timeout", "", "CLOSE_WAIT timeout")
	setCmd.Flags().String("time-wait-timeout", "", "TIME_WAIT timeout")
	setCmd.Flags().String("close-timeout", "", "CLOSE timeout")
	setCmd.Flags().String("udp-timeout", "", "UDP timeout")
	setCmd.Flags().String("udp-stream-timeout", "", "UDP stream timeout")
	setCmd.Flags().String("icmp-timeout", "", "ICMP timeout")
	group.AddCommand(getCmd, setCmd)
	return group
}

func systemCPUFreqCmd(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "cpufreq", Short: "CPU frequency config"}
	listCmd := readCommand(app, "list", "List CPU frequencies", "/system/cpufreq",
		[]string{"cpuid", "freq", "phyid", "coreid", "used", "softirq", "hardirq"}, nil)
	listCmd.Aliases = []string{"get"}
	modeCmd := &cobra.Command{Use: "mode", Short: "CPU frequency mode"}
	modeGetCmd := readCommand(app, "get", "Get CPU freq mode", "/system/cpufreq/mode",
		[]string{"cpufreq_support", "current_cpufreq", "current_turbo", "turbo_support"}, nil)
	modeSetCmd := fullUpdateCommand(app, "set", "Set CPU freq mode", "/system/cpufreq/mode", "/system/cpufreq/mode", cpufreqFieldMap, cpufreqModeFields, nil)
	modeSetCmd.Flags().String("mode", "", "CPU frequency mode")
	modeSetCmd.Flags().String("turbo", "", "Turbo boost (0/1)")
	modeCmd.AddCommand(modeGetCmd, modeSetCmd)
	group.AddCommand(listCmd, modeCmd)
	return group
}

func systemDisksCmd(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "disks", Short: "System disks"}
	group.AddCommand(readCommand(app, "list", "List system disks", "/system/disks",
		[]string{"disk", "model", "size", "type", "system", "block_size", "creating"}, nil))
	return group
}

func systemFilesCmd(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "files", Short: "System files"}
	listCmd := readCommand(app, "list", "List files", "/system/files",
		[]string{"f_name", "st_type", "st_size", "st_mtime", "st_inode"}, func(cmd *cobra.Command) map[string]string {
			path, _ := cmd.Flags().GetString("path")
			return map[string]string{"path": path}
		})
	listCmd.Flags().String("path", "", "File path (required)")
	cliapp.MarkFlagsRequired(listCmd, "path")
	wrapRequireFlags(listCmd, "path")
	group.AddCommand(listCmd)
	return group
}

func systemBackupCmd(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "backup", Short: "System backup snapshots"}
	listCmd := readCommand(app, "list", "List backup snapshots", "/system/backup",
		[]string{"id", "timestamp", "filename", "backtype", "version", "filesize"}, nil)
	createCmd := postActionCommand(app, "create", "Create backup snapshot", "/system/backup", map[string]string{})
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete backup snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			srcfile, _ := cmd.Flags().GetString("srcfile")
			if err := cliapp.RequireFlags(cmd, "srcfile"); err != nil {
				return err
			}
			yes, _ := cmd.Flags().GetBool("yes")
			if err := cliapp.ConfirmDelete(app.Stdout, app.Stderr, "backup", srcfile, yes); err != nil {
				return err
			}
			raw, err := app.APIClient.Delete(cliapp.APIBase + "/system/backup?srcfile=" + srcfile)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	deleteCmd.Flags().String("srcfile", "", "Backup filename (required)")
	deleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	cliapp.MarkFlagsRequired(deleteCmd, "srcfile")
	restoreCmd := writeCommand(app, "restore", "Restore backup snapshot", map[string]string{
		"srcfile":         "srcfile",
		"restore-type":    "restore_type",
		"sync-bind-cloud": "sync_bind_cloud",
		"cloud-comment":   "cloud_comment",
	}, map[string]interface{}{"restore_type": 0, "sync_bind_cloud": 0, "cloud_comment": ""}, func(body interface{}, _ string) (json.RawMessage, error) {
		return app.APIClient.Post(cliapp.APIBase+"/system/backup:restore", body)
	})
	restoreCmd.Flags().String("srcfile", "", "Backup filename (required)")
	restoreCmd.Flags().String("restore-type", "", "Restore type")
	restoreCmd.Flags().String("sync-bind-cloud", "", "Sync cloud binding (0/1)")
	restoreCmd.Flags().String("cloud-comment", "", "Cloud comment")
	restoreCmd.Flags().BoolP("yes", "y", false, "Confirm restore")
	cliapp.MarkFlagsRequired(restoreCmd, "srcfile")
	wrapRequireFlags(restoreCmd, "srcfile")
	wrapRequireYesUnlessDryRun(restoreCmd, app, "backup restore requires --yes")
	group.AddCommand(listCmd, createCmd, deleteCmd, restoreCmd)
	return group
}

func systemBackupAutoCmd(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "backup-auto", Short: "Automatic backup policy"}
	getCmd := readCommand(app, "get", "Get automatic backup policy", "/system/backup-auto",
		[]string{"id", "enabled", "strategy", "time", "cycle_time", "valid_days"}, nil)
	setCmd := fullUpdateCommand(app, "set", "Set automatic backup policy", "/system/backup-auto", "/system/backup-auto", backupAutoMap, backupAutoFields, nil)
	setCmd.Flags().String("id", "", "Policy ID")
	setCmd.Flags().String("enabled", "", "Enabled (yes/no)")
	setCmd.Flags().String("strategy", "", "Strategy (one/week/month)")
	setCmd.Flags().String("time", "", "Run time HH:MM")
	setCmd.Flags().String("cycle-time", "", "Cycle time")
	setCmd.Flags().String("valid-days", "", "Retention days")
	group.AddCommand(getCmd, setCmd)
	return group
}

func systemUpgradeCmd(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "upgrade", Short: "System upgrade"}
	group.AddCommand(
		readCommand(app, "get", "Get upgrade info", "/system/upgrade",
			[]string{"system_ver", "build_date", "new_system_ver", "new_build_date", "firmware_channel"}, nil),
		postActionCommand(app, "check", "Check for upgrade", "/system/upgrade:check", map[string]string{}),
		readCommand(app, "status", "Get upgrade status", "/system/upgrade:status",
			[]string{"auto_upgrade_status", "auto_upgrade_status_msg"}, nil),
	)
	startCmd := postActionCommand(app, "start", "Start system upgrade", "/system/upgrade:start", map[string]string{"type": "system"})
	startCmd.Flags().BoolP("yes", "y", false, "Confirm upgrade")
	wrapRequireYesUnlessDryRun(startCmd, app, "system upgrade start requires --yes")
	group.AddCommand(startCmd)
	return group
}

func systemWebAdminCmd(app *cliapp.Runtime) *cobra.Command {
	group := &cobra.Command{Use: "web-admin", Short: "Web admin accounts and groups"}
	groups := &cobra.Command{Use: "groups", Short: "Web admin groups"}
	groups.AddCommand(webAdminCRUD(app, "group", "/system/web-admin/groups", webAdminGroupMap, webAdminGroupFields,
		[]string{"id", "group_name", "ip_addr", "perm_config"}, []string{"name", "ip-addr", "perm-config"}, addWebAdminGroupFlags)...)

	accounts := &cobra.Command{Use: "accounts", Short: "Web admin accounts"}
	accounts.AddCommand(webAdminCRUD(app, "account", "/system/web-admin/accounts", webAdminAccountMap, webAdminAccountFields,
		[]string{"id", "username", "enabled", "group_id", "force", "interval", "sesstimeout", "comment"}, []string{"username", "passwd-md5", "group-id"}, addWebAdminAccountFlags)...)

	passwordStatus := readCommand(app, "password-status", "Get password status", "/system/web-admin/password-status",
		[]string{"mod_passwd"}, func(cmd *cobra.Command) map[string]string {
			username, _ := cmd.Flags().GetString("username")
			return map[string]string{"username": username}
		})
	passwordStatus.Flags().String("username", "", "Username (required)")
	cliapp.MarkFlagsRequired(passwordStatus, "username")
	wrapRequireFlags(passwordStatus, "username")

	passwordCmd := &cobra.Command{Use: "password", Short: "Web admin password"}
	passwordSet := writeCommand(app, "set", "Set current web admin password", map[string]string{"passwd-md5": "passwd"}, nil, func(body interface{}, _ string) (json.RawMessage, error) {
		return app.APIClient.Put(cliapp.APIBase+"/system/web-admin/password", body)
	})
	passwordSet.Flags().String("passwd-md5", "", "New password MD5 (required)")
	passwordSet.Flags().BoolP("yes", "y", false, "Confirm password update")
	cliapp.MarkFlagsRequired(passwordSet, "passwd-md5")
	wrapRequireFlags(passwordSet, "passwd-md5")
	wrapRequireYesUnlessDryRun(passwordSet, app, "web admin password set requires --yes")
	passwordCmd.AddCommand(passwordSet)

	group.AddCommand(groups, accounts, passwordStatus, passwordCmd)
	return group
}

func readCommand(app *cliapp.Runtime, use, short, apiPath string, columns []string, params func(*cobra.Command) map[string]string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if len(columns) > 0 {
				app.DefaultColumns = columns
			}
			path := cliapp.APIBase + apiPath
			if strings.Contains(use, "ID") {
				path += "/" + args[0]
			}
			raw, err := app.APIClient.Get(path, paramsFor(cmd, params))
			if err != nil {
				return err
			}
			if strings.Contains(use, "ID") {
				raw = filterRawByID(raw, args[0])
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	if strings.Contains(use, "ID") {
		cmd.Args = cobra.ExactArgs(1)
	}
	return cmd
}

func postActionCommand(app *cliapp.Runtime, use, short, apiPath string, body map[string]string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			raw, err := app.APIClient.Post(cliapp.APIBase+apiPath, body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
}

type writeCall func(body interface{}, id string) (json.RawMessage, error)

func writeCommand(app *cliapp.Runtime, use, short string, fieldMap map[string]string, defaults map[string]interface{}, call writeCall) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body, err := cliapp.MergeDataWithFlags(data, cmd, fieldMap)
			if err != nil {
				return err
			}
			for k, v := range defaults {
				if _, ok := body[k]; !ok {
					body[k] = v
				}
			}
			id := ""
			if strings.Contains(use, "ID") {
				id = args[0]
			}
			raw, err := call(body, id)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	if strings.Contains(use, "ID") {
		cmd.Args = cobra.ExactArgs(1)
	}
	cmd.Flags().String("data", "{}", "JSON body")
	return cmd
}

func fullUpdateCommand(app *cliapp.Runtime, use, short, getPath, putPath string, fieldMap map[string]string, allowedFields []string, defaults map[string]interface{}) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			overrides, err := cliapp.MergeDataWithFlags(data, cmd, fieldMap)
			if err != nil {
				return err
			}
			pathSuffix := ""
			targetID := ""
			if strings.Contains(use, "ID") {
				targetID = args[0]
				pathSuffix = "/" + targetID
			}
			body, err := fullUpdateBody(app, cliapp.APIBase+getPath+pathSuffix, overrides, allowedFields, defaults, targetID)
			if err != nil {
				return err
			}
			raw, err := app.APIClient.Put(cliapp.APIBase+putPath+pathSuffix, body)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	if strings.Contains(use, "ID") {
		cmd.Args = cobra.ExactArgs(1)
	}
	cmd.Flags().String("data", "{}", "JSON body")
	return cmd
}

func deleteIDCommand(app *cliapp.Runtime, use, short, resource, apiPath string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			yes, _ := cmd.Flags().GetBool("yes")
			if err := cliapp.ConfirmDelete(app.Stdout, app.Stderr, resource, args[0], yes); err != nil {
				return err
			}
			raw, err := app.APIClient.Delete(cliapp.APIBase + apiPath + "/" + args[0])
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	return cmd
}

func webAdminCRUD(app *cliapp.Runtime, resource, apiPath string, fieldMap map[string]string, allowedFields, columns, requiredCreate []string, addFlags func(*cobra.Command)) []*cobra.Command {
	listCmd := readCommand(app, "list", "List web admin "+resource+"s", apiPath, columns, webAdminListParams)
	listCmd.Aliases = []string{"ls"}
	cliapp.AddPaginationFlags(listCmd)

	getCmd := readCommand(app, "get ID", "Get web admin "+resource, apiPath, columns, nil)

	createDefaults := map[string]interface{}{}
	if resource == "account" {
		createDefaults = webAdminAccountDefaults
	}
	createCmd := writeCommand(app, "create", "Create web admin "+resource, fieldMap, createDefaults, func(body interface{}, _ string) (json.RawMessage, error) {
		return app.APIClient.Post(cliapp.APIBase+apiPath, body)
	})
	createCmd.Aliases = []string{"new"}
	addFlags(createCmd)
	cliapp.MarkFlagsRequired(createCmd, requiredCreate...)
	wrapRequireFlags(createCmd, requiredCreate...)

	updateCmd := fullUpdateCommand(app, "update ID", "Update web admin "+resource, apiPath, apiPath, fieldMap, allowedFields, nil)
	addFlags(updateCmd)

	deleteCmd := deleteIDCommand(app, "delete ID", "Delete web admin "+resource, resource, apiPath)
	deleteCmd.Aliases = []string{"rm"}
	return []*cobra.Command{listCmd, getCmd, createCmd, updateCmd, deleteCmd}
}

func addScheduleFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Schedule name (tagname)")
	cmd.Flags().String("event", "", "Event type (reboot/poweroff)")
	cmd.Flags().String("time", "", "Schedule time (HH:MM)")
	cmd.Flags().String("strategy", "", "Schedule strategy (one/day/week/month)")
	cmd.Flags().String("cycle-time", "", "Cycle time")
	cmd.Flags().String("comment", "", "Comment")
	cliapp.AddEnabledFlag(cmd)
}

func addWebAdminGroupFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Group name")
	cmd.Flags().String("ip-addr", "", "Allowed source IP range")
	cmd.Flags().String("perm-config", "", "Permission config")
}

func addWebAdminAccountFlags(cmd *cobra.Command) {
	cmd.Flags().String("username", "", "Username")
	cmd.Flags().String("passwd-md5", "", "Password MD5")
	cliapp.AddEnabledFlag(cmd)
	cmd.Flags().String("group-id", "", "Group ID")
	cmd.Flags().String("force", "", "Force password rotation (0/1)")
	cmd.Flags().String("interval", "", "Password rotation interval")
	cmd.Flags().String("session-timeout", "", "Session timeout")
	cmd.Flags().String("comment", "", "Comment")
}

func addPageAndOrderFlags(cmd *cobra.Command) {
	cliapp.AddPaginationFlags(cmd)
	cmd.Flags().String("order", "", "Sort direction: asc|desc")
	cmd.Flags().String("order-by", "", "Sort field")
}

func scheduleListParams(cmd *cobra.Command) map[string]string {
	page, pageSize, _, order, orderBy := cliapp.GetListParams(cmd)
	return cliapp.ListParamsWithPageSizeKey(page, pageSize, "", order, orderBy, "limit")
}

func webAdminListParams(cmd *cobra.Command) map[string]string {
	page, pageSize, _, _, _ := cliapp.GetListParams(cmd)
	return cliapp.ListParamsWithPageSizeKey(page, pageSize, "", "", "", "limit")
}

func paramsFor(cmd *cobra.Command, fn func(*cobra.Command) map[string]string) map[string]string {
	if fn == nil {
		return nil
	}
	return fn(cmd)
}

func wrapRequireFlags(cmd *cobra.Command, flags ...string) {
	origRunE := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if err := cliapp.RequireFlags(cmd, flags...); err != nil {
			return err
		}
		return origRunE(cmd, args)
	}
}

func wrapRequireYesUnlessDryRun(cmd *cobra.Command, app *cliapp.Runtime, message string) {
	origRunE := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		yes, _ := cmd.Flags().GetBool("yes")
		if !app.DryRun && !yes {
			return &cliapp.ValidationError{Message: message}
		}
		return origRunE(cmd, args)
	}
}

func fullUpdateBody(app *cliapp.Runtime, getPath string, overrides map[string]interface{}, allowedFields []string, defaults map[string]interface{}, targetID string) (map[string]interface{}, error) {
	raw, err := getBaseline(app, getPath)
	if err != nil {
		return nil, err
	}
	base, err := firstObject(raw, targetID)
	if err != nil {
		return nil, err
	}
	allowed := make(map[string]bool, len(allowedFields))
	for _, field := range allowedFields {
		allowed[field] = true
	}
	body := make(map[string]interface{}, len(allowedFields))
	for _, field := range allowedFields {
		if v, ok := base[field]; ok {
			body[field] = v
		}
	}
	for k, v := range defaults {
		if allowed[k] {
			if _, exists := body[k]; !exists {
				body[k] = v
			}
		}
	}
	for k, v := range overrides {
		if allowed[k] {
			body[k] = v
		}
	}
	return body, nil
}

func getBaseline(app *cliapp.Runtime, path string) (json.RawMessage, error) {
	wasDryRun := app.APIClient.DryRun
	app.APIClient.DryRun = false
	raw, err := app.APIClient.Get(path, nil)
	app.APIClient.DryRun = wasDryRun
	return raw, err
}

func firstObject(raw json.RawMessage, targetID string) (map[string]interface{}, error) {
	var v interface{}
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, err
	}
	if obj, ok := firstObjectFrom(v, targetID); ok {
		return obj, nil
	}
	return nil, fmt.Errorf("could not find object in API response")
}

func firstObjectFrom(v interface{}, targetID string) (map[string]interface{}, bool) {
	switch data := v.(type) {
	case []interface{}:
		if len(data) == 0 {
			return nil, false
		}
		if targetID != "" {
			for _, item := range data {
				if obj, ok := item.(map[string]interface{}); ok && matchesID(obj, targetID) {
					return obj, true
				}
			}
		}
		return firstObjectFrom(data[0], targetID)
	case map[string]interface{}:
		if targetID != "" && matchesID(data, targetID) {
			return data, true
		}
		for _, key := range []string{"data", "items", "groups_data", "accounts_data", "backup_info", "auto_backup"} {
			if inner, ok := data[key]; ok {
				if obj, found := firstObjectFrom(inner, targetID); found {
					return obj, true
				}
			}
		}
		return data, true
	default:
		return nil, false
	}
}

func filterRawByID(raw json.RawMessage, targetID string) json.RawMessage {
	var v interface{}
	if err := json.Unmarshal(raw, &v); err != nil {
		return raw
	}
	filtered := filterValueByID(v, targetID)
	out, err := json.Marshal(filtered)
	if err != nil {
		return raw
	}
	return out
}

func filterValueByID(v interface{}, targetID string) interface{} {
	switch data := v.(type) {
	case []interface{}:
		return filterArrayByID(data, targetID)
	case map[string]interface{}:
		out := make(map[string]interface{}, len(data))
		for k, val := range data {
			switch k {
			case "data", "items", "groups_data", "accounts_data", "backup_info":
				if arr, ok := val.([]interface{}); ok {
					out[k] = filterArrayByID(arr, targetID)
					continue
				}
			}
			out[k] = filterValueByID(val, targetID)
		}
		return out
	default:
		return v
	}
}

func filterArrayByID(items []interface{}, targetID string) []interface{} {
	filtered := make([]interface{}, 0, len(items))
	for _, item := range items {
		if obj, ok := item.(map[string]interface{}); ok && matchesID(obj, targetID) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func matchesID(obj map[string]interface{}, targetID string) bool {
	id, ok := obj["id"]
	if !ok {
		return false
	}
	return fmt.Sprint(id) == targetID
}
