package start

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/actions/db/connect"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var cmd = &cobra.Command{
	Use:   "start <type> <name>",
	Short: "Start an SSM tunnel to a database (type: postgres, redis, mongo)",
	Args:  cobra.ExactArgs(2),
	Run:   run,
}

func GetCmd() *cobra.Command {
	cmd.Flags().IntP("port", "p", 0, "Override local port (default: use config value)")
	return cmd
}

func run(c *cobra.Command, args []string) {
	localPort, err := c.Flags().GetInt("port")
	if err != nil {
		log.Exitf(1, "invalid port flag: %s", err.Error())
	}

	input := connect.StartInput{
		DBType:    args[0],
		Name:      args[1],
		LocalPort: localPort,
	}

	if err := connect.Start(input); err != nil {
		log.Exitf(1, "failed to start tunnel: %s", err.Error())
	}
}
