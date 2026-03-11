package dev

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/internal/actions/deploy"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var cmd = &cobra.Command{
	Use:   "deploy:dev <service|path> [service|path...]",
	Short: "Deploy one or more services to dev",
	Long: `Deploy Serverless services to dev without changing directories.

Accepts service names (from serverless.yml service: field) or paths.
Searches the git repository root to resolve service names.

Examples:
  draft deploy:dev gamestats
  draft deploy:dev gamestats notification-engine-v2 usertracking
  draft deploy:dev services/gamestats
  draft deploy:dev gamestats services/other-service`,
	Run:  run,
	Args: cobra.MinimumNArgs(1),
}

func GetCmd() *cobra.Command { return cmd }

func run(c *cobra.Command, args []string) {
	common.ChDir(c)

	results := deploy.DeployService(deploy.DevEnv, args)

	if len(results) > 1 {
		log.Info("\n─── Deploy Summary ───")
		for _, r := range results {
			if r.Err != nil {
				log.Errorf("✗ %s: %v", r.Name, r.Err)
			} else {
				log.Successf("✓ %s", r.Name)
			}
		}
	}

	for _, r := range results {
		if r.Err != nil {
			log.Exitf(1, "one or more deploys failed")
		}
	}
}
