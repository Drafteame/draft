package start

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/actions/db/connect"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var cmd = &cobra.Command{
	Use:   "start <type> <name>",
	Short: "Start an SSM tunnel to a database (type: postgres, redis, mongo)",
	Long: `Start an SSM port-forwarding tunnel to a database.

Postgres connections go through the RDS Proxy with IAM authentication. Once the
tunnel is up an RDS IAM authentication token is generated and written to the
PostgreSQL password file ($PGPASSFILE or ~/.pgpass), so psql and any libpq
client connect without a password. The token is only valid for new logins
during 15 minutes; running the same command again while the tunnel is up just
refreshes the token without restarting anything.

The PostgreSQL user is your Draftea email, taken from DRAFTEA_EMAIL (set with
hms-set-email).

Examples:
  draft db:connect start postgres turbo-dev
  draft db:connect start postgres turbo-dev --port 15432
  draft db:connect start redis api-cache-dev`,
	Args: cobra.ExactArgs(2),
	Run:  run,
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

	if localPort < 0 || localPort > 65535 {
		log.Exitf(1, "invalid port %d: must be between 0 and 65535", localPort)
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
