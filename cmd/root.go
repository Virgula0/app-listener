package cmd

import (
	"os"

	"github.com/Virgula0/app-listener/cmd/entity"
	"github.com/Virgula0/app-listener/cmd/functions/guard"
	"github.com/Virgula0/app-listener/cmd/functions/monitor"
	"github.com/Virgula0/app-listener/cmd/functions/networkguard"
	"github.com/Virgula0/app-listener/cmd/functions/networkmonitor"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app-listener",
	Short: "Monitor and guard file system events using eBPF",
	Long: `app-listener is a TUI/GUI application that monitors and guards
file system operations using eBPF.

Use --gui to launch the graphical interface instead of the terminal TUI
(only available for the monitor subcommand).`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		initLogger()
		return nil
	},
}

var logLevel string

func Execute() error {
	return rootCmd.Execute()
}

func initLogger() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:            true,
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
	log.SetOutput(os.Stderr)

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Warnf("invalid log level %q, using info", logLevel)
		level = log.InfoLevel
	}
	log.SetLevel(level)
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Set log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolVarP(&entity.GUI, "gui", "", false, "Launch GUI instead of TUI (default: false)")
	rootCmd.AddCommand(monitor.MonitorCmd)
	rootCmd.AddCommand(guard.GuardCmd)
	rootCmd.AddCommand(networkmonitor.NetworkMonitorCmd)
	rootCmd.AddCommand(networkguard.NetworkGuardCmd)
}
