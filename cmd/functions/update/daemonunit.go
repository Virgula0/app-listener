package update

import (
	"github.com/spf13/cobra"

	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/systemd"
)

// DaemonUnitCmd prints the daemon unit this build bundles: `update` and `install binary`
// run it on the binary they deploy, so the unit diffed is the new release's, not the running one's.
var DaemonUnitCmd = &cobra.Command{
	Use:   systemd.DaemonUnitCommand,
	Short: "Print the bundled daemon systemd unit (update runs it on the binary it deploys)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		unit, err := inst.SampleContent(systemd.DaemonServiceName + ".service")
		if err != nil {
			return err
		}
		_, err = cmd.OutOrStdout().Write(unit)
		return err
	},
}
