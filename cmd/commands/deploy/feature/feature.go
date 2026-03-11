package feature

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/internal/actions/deploy"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var cmd = &cobra.Command{
	Use:   "deploy:feature <service|path> [service|path...]",
	Short: "Deploy one or more services to feature",
	Long: `Deploy Serverless services to feature without changing directories.

Accepts service names (from serverless.yml service: field) or paths.
Searches the git repository root to resolve service names.

Examples:
  draft deploy:feature gamestats
  draft deploy:feature gamestats notification-engine-v2 usertracking
  draft deploy:feature services/gamestats
  draft deploy:feature gamestats services/other-service`,
	Run:  run,
	Args: cobra.MinimumNArgs(1),
}

func GetCmd() *cobra.Command { return cmd }

func run(c *cobra.Command, args []string) {
	common.ChDir(c)

	results := deploy.DeployService(deploy.FeatureEnv, args)

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
