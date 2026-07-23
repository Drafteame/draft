package list

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/actions/db/connect"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var cmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured database connections",
	Run:   run,
}

func GetCmd() *cobra.Command {
	return cmd
}

func run(_ *cobra.Command, _ []string) {
	if err := connect.List(); err != nil {
		log.Exitf(1, "failed to list connections: %s", err.Error())
	}
}
