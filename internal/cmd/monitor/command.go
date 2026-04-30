package monitor

import (
	"fmt"
	"strconv"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

func New(app *cliapp.Runtime) *cobra.Command {
	monitorCmd := &cobra.Command{
		Use:   "monitor",
		Short: "System monitoring",
		Long:  `Real-time system monitoring: CPU, memory, interfaces, online clients, protocols, and traffic statistics.`,
		Example: `  ikuai-cli monitor system
  ikuai-cli monitor cpu --time-range hour --start-time 1773300000 --end-time 1773303600 --aggregate avg
  ikuai-cli monitor memory --time-range hour --start-time 1773300000 --end-time 1773303600 --aggregate avg
  ikuai-cli monitor clients-online`,
	}

	monitorCmd.AddCommand(simpleGetCmd(app, "system", "System overview", "/monitoring/system", nil))
	monitorCmd.AddCommand(loadCmd(app, "cpu", "CPU load history", "/monitoring/cpu"))
	monitorCmd.AddCommand(loadCmd(app, "memory", "Memory usage history", "/monitoring/memory"))
	monitorCmd.AddCommand(loadCmd(app, "disk", "Disk usage history", "/monitoring/disk"))
	monitorCmd.AddCommand(loadCmd(app, "temp", "CPU temperature history", "/monitoring/cputemp"))
	monitorCmd.AddCommand(loadCmd(app, "terminals", "Terminal count history", "/monitoring/terminals"))
	monitorCmd.AddCommand(loadCmd(app, "connections", "Connection count history", "/monitoring/connections"))
	monitorCmd.AddCommand(loadCmd(app, "network-load", "Network load history", "/monitoring/network"))
	monitorCmd.AddCommand(legacyMonitorListCmd(app, "downstream", "Downstream traffic", "/monitoring/downstream",
		[]string{"id", "name", "ip_addr", "device", "method", "status", "port", "protocol", "enabled"}))

	onlineCols := []string{"id", "ip_addr", "mac", "hostname", "interface", "upload", "download", "connect_num", "uptime", "client_type", "comment"}
	offlineCols := []string{"id", "ip_addr", "mac", "hostname", "interface", "total_up", "total_down", "uptime", "downtime", "client_type"}

	for _, c := range []*cobra.Command{
		pagedCmd(app, "clients-online", "Online IPv4 clients", "/monitoring/clients-online", onlineCols),
		pagedCmd(app, "clients-offline", "Offline IPv4 clients", "/monitoring/clients-offline", offlineCols),
		pagedCmd(app, "clients-ip6-online", "Online IPv6 clients", "/monitoring/clients-ip6-online", onlineCols),
		pagedCmd(app, "clients-ip6-offline", "Offline IPv6 clients", "/monitoring/clients-ip6-offline", offlineCols),
		monitorTrafficSummaryCmd(app),
		monitorTrafficLoadCmd(app),
		monitorClientProtocolsCmd(app),
		monitorClientProtocolsHistoryCmd(app),
		monitorClientAppProtocolsCmd(app),
		monitorAppTrafficSummaryCmd(app),
		monitorProtocolsCmd(app),
		monitorProtocolsHistoryCmd(app),
		monitorAppProtocolsLoadCmd(app),
		monitorAppProtocolsHistoryCmd(app),
		monitorAppProtocolsTerminalsCmd(app),
		simpleGetCmd(app, "interfaces", "WAN interface status", "/monitoring/interfaces-status", nil),
		simpleGetCmd(app, "interfaces-traffic", "Per-interface traffic", "/monitoring/interfaces-traffic", nil),
		simpleGetCmd(app, "interfaces-config", "Interface config", "/monitoring/interfaces-config", nil),
		simpleGetCmd(app, "interfaces-physical", "Physical NIC info", "/monitoring/interfaces-physical", nil),
		simpleGetCmd(app, "interfaces-traffic-v6", "IPv6 interface traffic", "/monitoring/interfaces-traffic-v6", nil),
		simpleGetCmd(app, "wireless-stats", "Wireless statistics", "/monitoring/wireless-statistics", nil),
		simpleGetCmd(app, "wireless-score", "Wireless quality score", "/monitoring/wireless-score", nil),
		monitorWirelessTrafficCmd(app),
		monitorSSIDClientsCmd(app),
		monitorChannelClientsCmd(app),
		legacyMonitorListCmd(app, "cameras", "IP camera list", "/monitoring/cameras",
			[]string{"id", "tagname", "name", "vendor", "ip_addr", "mac", "flag", "status", "enabled"}),
		monitorFlowShuntingCmd(app),
		legacyMonitorListCmd(app, "switch", "Switch port monitoring", "/monitoring/switch",
			[]string{"id", "name", "ip_addr", "mac", "device", "version", "type", "status"}),
	} {
		monitorCmd.AddCommand(c)
	}

	return monitorCmd
}

// loadCmd creates a monitor command for load-type endpoints.
func loadCmd(app *cliapp.Runtime, use, short, apiPath string) *cobra.Command {
	c := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "time-range", "start-time", "end-time", "aggregate"); err != nil {
				return err
			}
			timeRange, _ := cmd.Flags().GetString("time-range")
			switch timeRange {
			case "hour", "day", "week", "month":
			default:
				return fmt.Errorf("--time-range must be one of: hour, day, week, month")
			}
			agg, _ := cmd.Flags().GetString("aggregate")
			if agg != "avg" && agg != "max" {
				return fmt.Errorf("--aggregate must be one of: avg, max")
			}
			end, _ := cmd.Flags().GetInt64("end-time")
			start, _ := cmd.Flags().GetInt64("start-time")
			if start >= end {
				return fmt.Errorf("--start-time must be less than --end-time")
			}
			p := map[string]string{
				"datetype":   timeRange,
				"start_time": strconv.FormatInt(start, 10),
				"end_time":   strconv.FormatInt(end, 10),
				"math":       agg,
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+apiPath, p)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	c.Flags().String("time-range", "", "Time range: hour, day, week, month (required)")
	c.Flags().Int64("start-time", 0, "Start Unix timestamp (required)")
	c.Flags().Int64("end-time", 0, "End Unix timestamp (required)")
	c.Flags().String("aggregate", "", "Aggregation: avg, max (required)")
	cliapp.MarkFlagsRequired(c, "time-range", "start-time", "end-time", "aggregate")
	return c
}

func simpleGetCmd(app *cliapp.Runtime, use, short, apiPath string, cols []string) *cobra.Command {
	return simpleGetWithParamsCmd(app, use, short, apiPath, nil, cols)
}

func simpleGetWithParamsCmd(app *cliapp.Runtime, use, short, apiPath string, params map[string]string, defaultCols []string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if len(defaultCols) > 0 {
				app.DefaultColumns = defaultCols
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+apiPath, params)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
}

func pagedCmd(app *cliapp.Runtime, use, short, apiPath string, defaultCols []string) *cobra.Command {
	c := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if len(defaultCols) > 0 {
				app.DefaultColumns = defaultCols
			}
			page, pageSize, filter, _, _ := cliapp.GetListParams(cmd)
			p := cliapp.ListParamsWithPageSizeKey(page, pageSize, filter, "", "", "limit")
			if key, _ := cmd.Flags().GetString("key"); key != "" {
				p["key"] = key
			}
			if pattern, _ := cmd.Flags().GetString("pattern"); pattern != "" {
				p["pattern"] = pattern
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+apiPath, p)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddPaginationFlags(c)
	c.Flags().String("filter", "", "Filter: field==value, & for AND, comma for OR")
	c.Flags().String("key", "", "Fuzzy match fields, comma-separated")
	c.Flags().String("pattern", "", "Fuzzy match pattern")
	return c
}

func legacyMonitorListCmd(app *cliapp.Runtime, use, short, apiPath string, defaultCols []string) *cobra.Command {
	c := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if len(defaultCols) > 0 {
				app.DefaultColumns = defaultCols
			}
			page, pageSize, _, order, orderBy := cliapp.GetListParams(cmd)
			if page < 1 {
				return fmt.Errorf("--page must be greater than 0")
			}
			if pageSize < 1 {
				return fmt.Errorf("--page-size must be greater than 0")
			}
			offset := (page - 1) * pageSize
			p := map[string]string{
				"TYPE":   "data",
				"LIMIT":  fmt.Sprintf("%d,%d", offset, pageSize),
				"OFFSET": intStr(offset),
			}
			if order != "" {
				p["ORDER"] = order
			}
			if orderBy != "" {
				p["ORDER_BY"] = orderBy
			}
			if status, _ := cmd.Flags().GetString("status"); status != "" {
				if status != "0" && status != "1" {
					return fmt.Errorf("--status must be 0 or 1")
				}
				p["status"] = status
			}
			if deviceFlag := cmd.Flags().Lookup("device"); deviceFlag != nil {
				if device, _ := cmd.Flags().GetString("device"); device != "" {
					p["device"] = device
				}
			}
			if keywordFlag := cmd.Flags().Lookup("keyword"); keywordFlag != nil {
				if keyword, _ := cmd.Flags().GetString("keyword"); keyword != "" {
					p["KEYWORDS"] = keyword
				}
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+apiPath, p)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddPaginationFlags(c)
	c.Flags().String("order", "", "Sort direction: asc|desc")
	c.Flags().String("order-by", "", "Sort field")
	c.Flags().String("status", "", "Status filter: 0|1")
	if use == "downstream" {
		c.Flags().String("device", "", "Device type: router|switches|firewall|server|camera|printer|SmartDevices|ap|other")
	} else {
		c.Flags().String("keyword", "", "Keyword search")
	}
	return c
}

func monitorTrafficLoadCmd(app *cliapp.Runtime) *cobra.Command {
	c := &cobra.Command{
		Use:   "traffic-load",
		Short: "Per-client 5-min traffic load (requires ip and mac)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "ip", "mac"); err != nil {
				return err
			}
			ip, _ := cmd.Flags().GetString("ip")
			mac, _ := cmd.Flags().GetString("mac")
			app.DefaultColumns = []string{"timestamp", "upload", "download", "conn_num"}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/clients-traffic-load", map[string]string{
				"ip": ip, "mac": mac,
			})
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	c.Flags().String("ip", "", "Client IP address (required)")
	c.Flags().String("mac", "", "Client MAC address (required)")
	cliapp.MarkFlagsRequired(c, "ip", "mac")
	return c
}

func monitorTrafficSummaryCmd(app *cliapp.Runtime) *cobra.Command {
	c := &cobra.Command{
		Use:   "traffic-summary",
		Short: "Client traffic summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			page, pageSize, _, _, _ := cliapp.GetListParams(cmd)
			app.DefaultColumns = []string{"id", "ip_addr", "mac", "username", "sum_total_up", "sum_total_down", "sum_total", "comment"}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/clients-traffic-summary", map[string]string{
				"page":  intStr(page),
				"limit": intStr(pageSize),
			})
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddPaginationFlags(c)
	return c
}

func monitorClientProtocolsCmd(app *cliapp.Runtime) *cobra.Command {
	c := &cobra.Command{
		Use:   "client-protocols",
		Short: "Per-client protocol breakdown (requires ip and mac)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "ip", "mac"); err != nil {
				return err
			}
			ip, _ := cmd.Flags().GetString("ip")
			mac, _ := cmd.Flags().GetString("mac")
			app.DefaultColumns = []string{"id", "proto", "proto_name", "total"}
			params := map[string]string{"ip": ip, "mac": mac}
			if starttime, _ := cmd.Flags().GetInt64("starttime"); starttime > 0 {
				params["starttime"] = strconv.FormatInt(starttime, 10)
			}
			if stoptime, _ := cmd.Flags().GetInt64("stoptime"); stoptime > 0 {
				params["stoptime"] = strconv.FormatInt(stoptime, 10)
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/clients/protocols", params)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	c.Flags().String("ip", "", "Client IP address (required)")
	c.Flags().String("mac", "", "Client MAC address (required)")
	c.Flags().Int64("starttime", 0, "Start Unix timestamp")
	c.Flags().Int64("stoptime", 0, "Stop Unix timestamp")
	cliapp.MarkFlagsRequired(c, "ip", "mac")
	return c
}

func monitorClientProtocolsHistoryCmd(app *cliapp.Runtime) *cobra.Command {
	c := &cobra.Command{
		Use:   "client-protocols-history",
		Short: "Per-client protocol rate history (requires ip and mac)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "ip", "mac"); err != nil {
				return err
			}
			ip, _ := cmd.Flags().GetString("ip")
			mac, _ := cmd.Flags().GetString("mac")
			app.DefaultColumns = []string{"id", "proto_name", "timestamp", "upload", "download", "total"}
			params := map[string]string{"ip": ip, "mac": mac}
			if starttime, _ := cmd.Flags().GetInt64("starttime"); starttime > 0 {
				params["starttime"] = strconv.FormatInt(starttime, 10)
			}
			if stoptime, _ := cmd.Flags().GetInt64("stoptime"); stoptime > 0 {
				params["stoptime"] = strconv.FormatInt(stoptime, 10)
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/clients/protocols/history-load", params)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	c.Flags().String("ip", "", "Client IP address (required)")
	c.Flags().String("mac", "", "Client MAC address (required)")
	c.Flags().Int64("starttime", 0, "Start Unix timestamp")
	c.Flags().Int64("stoptime", 0, "Stop Unix timestamp")
	cliapp.MarkFlagsRequired(c, "ip", "mac")
	return c
}

func monitorClientAppProtocolsCmd(app *cliapp.Runtime) *cobra.Command {
	c := &cobra.Command{
		Use:   "client-app-protocols",
		Short: "Per-client app-protocol rate (requires ip and mac)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "ip", "mac"); err != nil {
				return err
			}
			ip, _ := cmd.Flags().GetString("ip")
			mac, _ := cmd.Flags().GetString("mac")
			app.DefaultColumns = []string{"id", "appid", "appname", "conn_cnt", "upload", "download", "total"}
			params := map[string]string{"ip": ip, "mac": mac}
			if limit, _ := cmd.Flags().GetInt("page-size"); limit > 0 {
				params["limit"] = intStr(limit)
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/clients/app-protocols/load", params)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	c.Flags().String("ip", "", "Client IP address (required)")
	c.Flags().String("mac", "", "Client MAC address (required)")
	c.Flags().Int("page-size", 20, "Items per page")
	cliapp.MarkFlagsRequired(c, "ip", "mac")
	return c
}

func monitorAppTrafficSummaryCmd(app *cliapp.Runtime) *cobra.Command {
	c := &cobra.Command{
		Use:   "app-traffic-summary",
		Short: "App traffic summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			page, pageSize, _, _, _ := cliapp.GetListParams(cmd)
			app.DefaultColumns = []string{"id", "appname", "appname_level1", "appname_level2", "total_up", "total_down", "total", "appid"}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/app-traffic-summary",
				map[string]string{"page": intStr(page), "limit": intStr(pageSize)})
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddPaginationFlags(c)
	return c
}

func monitorProtocolsCmd(app *cliapp.Runtime) *cobra.Command {
	c := &cobra.Command{
		Use:   "protocols",
		Short: "Protocol distribution",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			params := timeRangeParams(cmd, "starttime", "stoptime")
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/protocols", params)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	addStartStopFlags(c)
	return c
}

func monitorProtocolsHistoryCmd(app *cliapp.Runtime) *cobra.Command {
	c := &cobra.Command{
		Use:   "protocols-history",
		Short: "Protocol history",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			app.DefaultColumns = []string{"proto_name", "timestamp", "upload", "download", "total"}
			params := timeRangeParams(cmd, "starttime", "stoptime")
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/protocols/history-load", params)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	addStartStopFlags(c)
	return c
}

func monitorAppProtocolsLoadCmd(app *cliapp.Runtime) *cobra.Command {
	c := &cobra.Command{
		Use:   "app-protocols-load",
		Short: "App-protocol load",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			page, pageSize, _, order, orderBy := cliapp.GetListParams(cmd)
			app.DefaultColumns = []string{"id", "appname", "conn_cnt", "upload", "download", "total_up", "total_down", "total"}
			p := map[string]string{"page": intStr(page), "limit": intStr(pageSize)}
			if order != "" {
				p["order"] = order
			}
			if orderBy != "" {
				p["order_by"] = orderBy
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/app-protocols/load", p)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddPaginationFlags(c)
	c.Flags().String("order", "", "Sort direction: asc|desc")
	c.Flags().String("order-by", "", "Sort field")
	return c
}

func monitorAppProtocolsHistoryCmd(app *cliapp.Runtime) *cobra.Command {
	c := &cobra.Command{
		Use:   "app-protocols-history",
		Short: "App-protocol rate history (requires appids)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "appids"); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "appname", "timestamp", "upload", "download", "total"}
			appids, _ := cmd.Flags().GetString("appids")
			params := map[string]string{"appids": appids}
			if starttime, _ := cmd.Flags().GetInt64("starttime"); starttime > 0 {
				params["starttime"] = strconv.FormatInt(starttime, 10)
			}
			if stoptime, _ := cmd.Flags().GetInt64("stoptime"); stoptime > 0 {
				params["stoptime"] = strconv.FormatInt(stoptime, 10)
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/app-protocols/history-load", params)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	c.Flags().String("appids", "", "App protocol IDs, comma-separated, e.g. 2580003,2580004")
	c.Flags().Int64("starttime", 0, "Start Unix timestamp")
	c.Flags().Int64("stoptime", 0, "Stop Unix timestamp")
	cliapp.MarkFlagsRequired(c, "appids")
	return c
}

func monitorAppProtocolsTerminalsCmd(app *cliapp.Runtime) *cobra.Command {
	c := &cobra.Command{
		Use:   "app-protocols-terminals",
		Short: "Terminals accessing an app protocol (requires appid)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if err := cliapp.RequireFlags(cmd, "appid"); err != nil {
				return err
			}
			app.DefaultColumns = []string{"id", "ipaddr", "mac", "hostname", "client_type", "conn_cnt", "upload", "download", "total_up"}
			appid, _ := cmd.Flags().GetInt("appid")
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/app-protocols/terminal-load", map[string]string{
				"appid": intStr(appid),
			})
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	c.Flags().Int("appid", 0, "App protocol ID, e.g. 2580003")
	cliapp.MarkFlagsRequired(c, "appid")
	return c
}

func monitorFlowShuntingCmd(app *cliapp.Runtime) *cobra.Command {
	c := &cobra.Command{
		Use:   "flow-shunting",
		Short: "Traffic shunting",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			value, _ := cmd.Flags().GetString("type")
			if value != "data" {
				return fmt.Errorf("--type must be data")
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/flow-shunting", map[string]string{
				"TYPE": value,
			})
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	c.Flags().String("type", "data", "Response type: data")
	return c
}

func monitorWirelessTrafficCmd(app *cliapp.Runtime) *cobra.Command {
	return simpleGetWithOptionalParamsCmd(app, "wireless-traffic", "Per-AP traffic", "/monitoring/wireless-traffic", nil, func(c *cobra.Command) {
		c.Flags().String("apmac", "", "AP MAC address")
	})
}

func monitorSSIDClientsCmd(app *cliapp.Runtime) *cobra.Command {
	return simpleGetWithOptionalParamsCmd(app, "ssid-clients", "SSID client distribution", "/monitoring/ssid-clients", nil, func(c *cobra.Command) {
		c.Flags().String("ssid", "", "SSID name")
	})
}

func monitorChannelClientsCmd(app *cliapp.Runtime) *cobra.Command {
	return simpleGetWithOptionalParamsCmd(app, "channel-clients", "Channel client distribution", "/monitoring/channel-clients", nil, func(c *cobra.Command) {
		c.Flags().Int("channel", 0, "Wireless channel")
	})
}

func simpleGetWithOptionalParamsCmd(app *cliapp.Runtime, use, short, apiPath string, defaultCols []string, addFlags func(*cobra.Command)) *cobra.Command {
	c := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			if len(defaultCols) > 0 {
				app.DefaultColumns = defaultCols
			}
			params := map[string]string{}
			for _, name := range []string{"apmac", "ssid"} {
				if flag := cmd.Flags().Lookup(name); flag != nil {
					if val, _ := cmd.Flags().GetString(name); val != "" {
						params[name] = val
					}
				}
			}
			if flag := cmd.Flags().Lookup("channel"); flag != nil {
				if val, _ := cmd.Flags().GetInt("channel"); val > 0 {
					params["channel"] = intStr(val)
				}
			}
			if len(params) == 0 {
				params = nil
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+apiPath, params)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	addFlags(c)
	return c
}

func addStartStopFlags(c *cobra.Command) {
	c.Flags().Int64("starttime", 0, "Start Unix timestamp")
	c.Flags().Int64("stoptime", 0, "Stop Unix timestamp")
}

func timeRangeParams(cmd *cobra.Command, startName, stopName string) map[string]string {
	params := map[string]string{}
	if start, _ := cmd.Flags().GetInt64(startName); start > 0 {
		params[startName] = strconv.FormatInt(start, 10)
	}
	if stop, _ := cmd.Flags().GetInt64(stopName); stop > 0 {
		params[stopName] = strconv.FormatInt(stop, 10)
	}
	if len(params) == 0 {
		return nil
	}
	return params
}

func intStr(n int) string {
	return fmt.Sprint(n)
}
