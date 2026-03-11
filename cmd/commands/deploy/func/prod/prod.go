package prod

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/internal/actions/deploy"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var cmd = &cobra.Command{
	Use:   "deploy:func:prod <service|path> <function-name>",
	Short: "Deploy a single Lambda function to prod",
	Long: `Package and deploy a single Lambda function to prod.

The service argument can be a service name (from serverless.yml) or a path.

Examples:
  draft deploy:func:prod gamestats storegamestats
  draft deploy:func:prod services/gamestats storegamestats`,
	Run:  run,
	Args: cobra.ExactArgs(2),
}

func GetCmd() *cobra.Command { return cmd }

func run(c *cobra.Command, args []string) {
	common.ChDir(c)

	if err := deploy.DeployFunction(deploy.ProdEnv, args[0], args[1]); err != nil {
		log.Exitf(1, "deploy:func:prod failed: %v", err)
	}
}
