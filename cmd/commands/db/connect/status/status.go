package status

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/actions/db/connect"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var cmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of all active database tunnels",
	Run:   run,
}

func GetCmd() *cobra.Command {
	return cmd
}

func run(_ *cobra.Command, _ []string) {
	if err := connect.Status(); err != nil {
		log.Exitf(1, "failed to get status: %s", err.Error())
	}
}
