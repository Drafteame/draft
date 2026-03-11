package dev

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/internal/actions/deploy"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var cmd = &cobra.Command{
	Use:   "deploy:func:dev <service|path> <function-name>",
	Short: "Deploy a single Lambda function to dev",
	Long: `Package and deploy a single Lambda function to dev.

The service argument can be a service name (from serverless.yml) or a path.

Examples:
  draft deploy:func:dev gamestats storegamestats
  draft deploy:func:dev services/gamestats storegamestats`,
	Run:  run,
	Args: cobra.ExactArgs(2),
}

func GetCmd() *cobra.Command { return cmd }

func run(c *cobra.Command, args []string) {
	common.ChDir(c)

	if err := deploy.DeployFunction(deploy.DevEnv, args[0], args[1]); err != nil {
		log.Exitf(1, "deploy:func:dev failed: %v", err)
	}
}
