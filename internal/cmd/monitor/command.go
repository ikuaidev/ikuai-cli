package monitor

import (
	"fmt"

	"github.com/ikuaidev/ikuai-cli/internal/cliapp"
	"github.com/spf13/cobra"
)

func New(app *cliapp.Runtime) *cobra.Command {
	monitorCmd := &cobra.Command{
		Use:   "monitor",
		Short: "System monitoring",
		Long:  `Real-time system monitoring: CPU, memory, interfaces, online clients, protocols, and traffic statistics.`,
		Example: `  ikuai-cli monitor system
  ikuai-cli monitor cpu
  ikuai-cli monitor memory
  ikuai-cli monitor clients-online`,
	}

	monitorCmd.AddCommand(simpleGetCmd(app, "system", "System overview", "/monitoring/system"))
	monitorCmd.AddCommand(loadCmd(app, "cpu", "CPU load history", "/monitoring/cpu"))
	monitorCmd.AddCommand(loadCmd(app, "memory", "Memory usage history", "/monitoring/memory"))
	monitorCmd.AddCommand(loadCmd(app, "disk", "Disk usage history", "/monitoring/disk"))
	monitorCmd.AddCommand(loadCmd(app, "temp", "CPU temperature history", "/monitoring/cputemp"))
	monitorCmd.AddCommand(loadCmd(app, "terminals", "Terminal count history", "/monitoring/terminals"))
	monitorCmd.AddCommand(loadCmd(app, "connections", "Connection count history", "/monitoring/connections"))
	monitorCmd.AddCommand(loadCmd(app, "network-load", "Network load history", "/monitoring/network"))
	monitorCmd.AddCommand(simpleGetCmd(app, "downstream", "Downstream traffic", "/monitoring/downstream"))

	for _, c := range []*cobra.Command{
		pagedCmd(app, "clients-online", "Online IPv4 clients", "/monitoring/clients-online"),
		pagedCmd(app, "clients-offline", "Offline IPv4 clients", "/monitoring/clients-offline"),
		pagedCmd(app, "clients-ip6-online", "Online IPv6 clients", "/monitoring/clients-ip6-online"),
		pagedCmd(app, "clients-ip6-offline", "Offline IPv6 clients", "/monitoring/clients-ip6-offline"),
		pagedCmd(app, "traffic-summary", "Client traffic summary", "/monitoring/clients-traffic-summary"),
		monitorTrafficLoadCmd(app),
		monitorClientProtocolsCmd(app),
		monitorClientProtocolsHistoryCmd(app),
		monitorClientAppProtocolsCmd(app),
		monitorAppTrafficSummaryCmd(app),
		simpleGetCmd(app, "protocols", "Protocol distribution", "/monitoring/protocols"),
		simpleGetCmd(app, "protocols-history", "Protocol history", "/monitoring/protocols/history-load"),
		monitorAppProtocolsLoadCmd(app),
		monitorAppProtocolsHistoryCmd(app),
		monitorAppProtocolsTerminalsCmd(app),
		simpleGetCmd(app, "interfaces", "WAN interface status", "/monitoring/interfaces-status"),
		simpleGetCmd(app, "interfaces-traffic", "Per-interface traffic", "/monitoring/interfaces-traffic"),
		simpleGetCmd(app, "interfaces-config", "Interface config", "/monitoring/interfaces-config"),
		simpleGetCmd(app, "interfaces-physical", "Physical NIC info", "/monitoring/interfaces-physical"),
		simpleGetCmd(app, "interfaces-traffic-v6", "IPv6 interface traffic", "/monitoring/interfaces-traffic-v6"),
		simpleGetCmd(app, "wireless-stats", "Wireless statistics", "/monitoring/wireless-statistics"),
		simpleGetCmd(app, "wireless-score", "Wireless quality score", "/monitoring/wireless-score"),
		simpleGetCmd(app, "wireless-traffic", "Per-AP traffic", "/monitoring/wireless-traffic"),
		simpleGetCmd(app, "ssid-clients", "SSID client distribution", "/monitoring/ssid-clients"),
		simpleGetCmd(app, "channel-clients", "Channel client distribution", "/monitoring/channel-clients"),
		simpleGetCmd(app, "cameras", "IP camera list", "/monitoring/cameras"),
		simpleGetCmd(app, "flow-shunting", "Traffic shunting", "/monitoring/flow-shunting"),
		simpleGetCmd(app, "switch", "Switch port monitoring", "/monitoring/switch"),
	} {
		monitorCmd.AddCommand(c)
	}

	return monitorCmd
}

// loadCmd creates a monitor command for load-type endpoints that support
// --time-range and --aggregate query parameters.
// All flags are optional; the router uses server-side defaults when omitted.
func loadCmd(app *cliapp.Runtime, use, short, apiPath string) *cobra.Command {
	c := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			p := map[string]string{}
			if timeRange, _ := cmd.Flags().GetString("time-range"); timeRange != "" {
				switch timeRange {
				case "hour", "day", "week", "month":
					p["datetype"] = timeRange // API param name unchanged
				default:
					return fmt.Errorf("--time-range must be one of: hour, day, week, month")
				}
			}
			if agg, _ := cmd.Flags().GetString("aggregate"); agg != "" {
				if agg != "avg" && agg != "max" {
					return fmt.Errorf("--aggregate must be one of: avg, max")
				}
				p["math"] = agg // API param name unchanged
			}
			if len(p) == 0 {
				p = nil
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+apiPath, p)
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	c.Flags().String("time-range", "", "Time range: hour, day, week, month")
	c.Flags().String("aggregate", "", "Aggregation: avg, max")
	return c
}

func simpleGetCmd(app *cliapp.Runtime, use, short, apiPath string) *cobra.Command {
	return simpleGetWithParamsCmd(app, use, short, apiPath, nil)
}

func simpleGetWithParamsCmd(app *cliapp.Runtime, use, short, apiPath string, params map[string]string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
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

func pagedCmd(app *cliapp.Runtime, use, short, apiPath string) *cobra.Command {
	c := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.RequireAuth(); err != nil {
				return err
			}
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+apiPath,
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(c)
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
			ip, _ := cmd.Flags().GetString("ip")
			mac, _ := cmd.Flags().GetString("mac")
			if ip == "" || mac == "" {
				return fmt.Errorf("both --ip and --mac are required")
			}
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
	_ = c.MarkFlagRequired("ip")
	_ = c.MarkFlagRequired("mac")
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
			ip, _ := cmd.Flags().GetString("ip")
			mac, _ := cmd.Flags().GetString("mac")
			if ip == "" || mac == "" {
				return fmt.Errorf("both --ip and --mac are required")
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/clients/protocols", map[string]string{
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
	_ = c.MarkFlagRequired("ip")
	_ = c.MarkFlagRequired("mac")
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
			ip, _ := cmd.Flags().GetString("ip")
			mac, _ := cmd.Flags().GetString("mac")
			if ip == "" || mac == "" {
				return fmt.Errorf("both --ip and --mac are required")
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/clients/protocols/history-load", map[string]string{
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
	_ = c.MarkFlagRequired("ip")
	_ = c.MarkFlagRequired("mac")
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
			ip, _ := cmd.Flags().GetString("ip")
			mac, _ := cmd.Flags().GetString("mac")
			if ip == "" || mac == "" {
				return fmt.Errorf("both --ip and --mac are required")
			}
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/clients/app-protocols/load", map[string]string{
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
	_ = c.MarkFlagRequired("ip")
	_ = c.MarkFlagRequired("mac")
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
			page, pageSize, filter, order, orderBy := cliapp.GetListParams(cmd)
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/app-traffic-summary",
				cliapp.ListParams(page, pageSize, filter, order, orderBy))
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	cliapp.AddListFlags(c)
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
			p := map[string]string{"page": intStr(page), "page_size": intStr(pageSize)}
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
	cliapp.AddListFlags(c)
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
			appids, _ := cmd.Flags().GetString("appids")
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/app-protocols/history-load", map[string]string{
				"appids": appids,
			})
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	c.Flags().String("appids", "", "App protocol IDs, comma-separated (required, e.g. 2580003,2580004)")
	_ = c.MarkFlagRequired("appids")
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
			appid, _ := cmd.Flags().GetString("appid")
			raw, err := app.APIClient.Get(cliapp.APIBase+"/monitoring/app-protocols/terminal-load", map[string]string{
				"appid": appid,
			})
			if err != nil {
				return err
			}
			app.PrintRaw(raw)
			return nil
		},
	}
	c.Flags().String("appid", "", "App protocol ID (required, e.g. 2580003)")
	_ = c.MarkFlagRequired("appid")
	return c
}

func intStr(n int) string {
	return fmt.Sprint(n)
}
