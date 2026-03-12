package deploy

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	deployaction "github.com/Drafteame/draft/internal/actions/deploy"
	"github.com/Drafteame/draft/internal/pkg/log"
)

// NewFuncCmd returns a deploy function command for the given environment.
func NewFuncCmd(env deployaction.EnvConfig) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("deploy:func:%s <service|path> <function-name>", env.Stage),
		Short: fmt.Sprintf("Deploy a single Lambda function to %s", env.Stage),
		Long: fmt.Sprintf(`Package and deploy a single Lambda function to %s.

The service argument can be a service name (from serverless.yml) or a path.

Examples:
  draft deploy:func:%s gamestats storegamestats
  draft deploy:func:%s services/gamestats storegamestats`, env.Stage, env.Stage, env.Stage),
		Args: cobra.ExactArgs(2),
		Run: func(c *cobra.Command, args []string) {
			common.ChDir(c)

			if err := deployaction.DeployFunction(env, args[0], args[1]); err != nil {
				log.Exitf(1, "deploy:func:%s failed: %v", env.Stage, err)
			}
		},
	}
}

// NewServiceCmd returns a deploy service command for the given environment.
func NewServiceCmd(env deployaction.EnvConfig) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("deploy:%s <service|path> [service|path...]", env.Stage),
		Short: fmt.Sprintf("Deploy one or more services to %s", env.Stage),
		Long: fmt.Sprintf(`Deploy Serverless services to %s without changing directories.

Accepts service names (from serverless.yml service: field) or paths.
Searches the git repository root to resolve service names.

Examples:
  draft deploy:%s gamestats
  draft deploy:%s gamestats notification-engine-v2 usertracking
  draft deploy:%s services/gamestats
  draft deploy:%s gamestats services/other-service`, env.Stage, env.Stage, env.Stage, env.Stage, env.Stage),
		Args: cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			common.ChDir(c)

			results := deployaction.DeployService(env, args)

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
		},
	}
}
