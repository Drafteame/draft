package commands

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/pkg/log"
	nixversion "github.com/Drafteame/draft/internal/pkg/version/nix"
)

var Version = "development"

var rootCmd = &cobra.Command{
	Use:              "draft <command>",
	Example:          "draft new:service",
	Version:          Version,
	Run:              run,
	PersistentPreRun: persistentPreRun,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&data.Flags.WorkingDir, "working-dir", "w", "", "Working directory")
	rootCmd.PersistentFlags().BoolVarP(&data.Flags.Debug, "debug", "d", false, "Debug mode")
	rootCmd.PersistentFlags().BoolVarP(&data.Flags.TTY, "tty", "t", true, "TTY mode")
	rootCmd.Flags().Bool("use-dig", false, "Use uber dig for dependency injection")
	rootCmd.Flags().String("legacy-path", "", "Path to legacy service")
}

func run(cmd *cobra.Command, _ []string) {
	if err := cmd.Help(); err != nil {
		log.Exitf(1, "failed to print help: %s", err.Error())
	}
}

func persistentPreRun(_ *cobra.Command, _ []string) {
	nixversion.CheckNixModulesVersion()
}

func GetCmd() *cobra.Command {
	return rootCmd
}
