package feature

import (
	"github.com/spf13/cobra"

	deploycmd "github.com/Drafteame/draft/cmd/commands/deploy"
	deployaction "github.com/Drafteame/draft/internal/actions/deploy"
)

func GetCmd() *cobra.Command {
	return deploycmd.NewServiceCmd(deployaction.FeatureEnv)
}
