package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/cmd/functions/monitor"
)

var rootCmd = &cobra.Command{
	Use:   "app-listener",
	Short: "Monitor file system events using eBPF",
	Long:  `app-listener is a TUI application that monitors file system operations using eBPF.`,
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
	rootCmd.AddCommand(monitor.MonitorCmd)
}
