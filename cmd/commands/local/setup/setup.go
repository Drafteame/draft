package setup

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	localsetup "github.com/Drafteame/draft/internal/actions/local/setup"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var localSetupCmd = &cobra.Command{
	Use:   "local:setup [flags]",
	Short: "Setup local environment",
	Long:  "Setup local environment",
	Run:   run,
}

func GetCmd() *cobra.Command {
	return localSetupCmd
}

func init() {
	localSetupCmd.Flags().Bool("prune", false, "Prune environment before creation")
	localSetupCmd.Flags().Bool("bypass-docker", false, "Avoid running docker commands")
}

func run(cmd *cobra.Command, _ []string) {
	common.ChDir(cmd)

	prune, err := cmd.Flags().GetBool("prune")
	if err != nil {
		log.Exitf(1, "failed to obtain prune flag: %s", err.Error())
	}

	bypassDocker, err := cmd.Flags().GetBool("bypass-docker")
	if err != nil {
		log.Exitf(1, "failed to obtain bypass-docker flag: %s", err.Error())
	}

	input := localsetup.Input{
		Prune:        prune,
		BypassDocker: bypassDocker,
	}

	if errExec := localsetup.New(input).Exec(); errExec != nil {
		log.Exitf(1, "failed to setup test environment: %s", errExec.Error())
	}

	log.Success("Setup test environment completed")
}
