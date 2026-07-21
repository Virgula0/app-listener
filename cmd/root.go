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
}

func Execute() error {
	initLogger()
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
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func init() {
	rootCmd.AddCommand(monitor.MonitorCmd)
}
