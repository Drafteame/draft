package stop

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/actions/db/connect"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var cmd = &cobra.Command{
	Use:   "stop [type] [name]",
	Short: "Stop active SSM tunnel(s)",
	Long: `Stop one or more active SSM tunnels.

Examples:
  draft dbconnect stop --all              Stop all active connections
  draft dbconnect stop postgres           Stop all active postgres connections
  draft dbconnect stop redis              Stop all active redis connections
  draft dbconnect stop postgres users-dev Stop a specific connection`,
	Args: cobra.RangeArgs(0, 2),
	Run:  run,
}

func GetCmd() *cobra.Command {
	cmd.Flags().Bool("all", false, "Stop all active connections")
	return cmd
}

func run(c *cobra.Command, args []string) {
	all, err := c.Flags().GetBool("all")
	if err != nil {
		log.Exitf(1, "invalid flag: %s", err.Error())
	}

	if !all && len(args) == 0 {
		log.Exit(1, "specify a type, a type and name, or use --all\n\n"+
			"  draft dbconnect stop --all\n"+
			"  draft dbconnect stop postgres\n"+
			"  draft dbconnect stop postgres users-dev")
	}

	input := connect.StopInput{}

	if !all {
		input.DBType = args[0]

		if len(args) == 2 {
			input.Name = args[1]
		}
	}

	if err := connect.Stop(input); err != nil {
		log.Exitf(1, "failed to stop: %s", err.Error())
	}
}
