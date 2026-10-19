// Package network groups the network modes: `network monitor` and `network guard`.
package network

import (
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/cmd/functions/networkguard"
	"github.com/Virgula0/app-listener/cmd/functions/networkmonitor"
)

var NetworkCmd = &cobra.Command{
	Use:     "network",
	Aliases: []string{"net"},
	Short:   "Monitor or deny the network operations of binaries",
	Args:    cobra.NoArgs,
}

func init() {
	NetworkCmd.AddCommand(networkmonitor.NetworkMonitorCmd, networkguard.NetworkGuardCmd)
}
