package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	localsetup "github.com/Drafteame/draft/internal/actions/local/setup"
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
	defer func() {
		if r := recover(); r != nil {
			_, _ = fmt.Fprintln(os.Stderr, r)
			os.Exit(1)
		}
	}()

	workDir, err := cmd.Parent().Flags().GetString("working-dir")
	if err != nil {
		panic(err)
	}

	prune, err := cmd.Flags().GetBool("prune")
	if err != nil {
		panic(err)
	}

	bypassDocker, err := cmd.Flags().GetBool("bypass-docker")
	if err != nil {
		panic(err)
	}

	input := localsetup.Input{
		WorkingDir:   workDir,
		Prune:        prune,
		BypassDocker: bypassDocker,
	}

	if err := localsetup.New(input).Exec(); err != nil {
		panic(err)
	}

	println("Setup test environment completed")
	os.Exit(0)
}
