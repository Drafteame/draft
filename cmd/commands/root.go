package commands

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var Version = "development"

var rootCmd = &cobra.Command{
	Use:     "draft <command>",
	Example: "draft new:service",
	Version: Version,
	Run:     run,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&data.Flags.WorkingDir, "working-dir", "w", "", "Working directory")
	rootCmd.PersistentFlags().BoolVarP(&data.Flags.Debug, "debug", "d", false, "Debug mode")
	rootCmd.PersistentFlags().BoolVarP(&data.Flags.TTY, "tty", "t", true, "TTY mode")
}

func run(cmd *cobra.Command, _ []string) {
	if err := cmd.Help(); err != nil {
		log.Exitf(1, "failed to print help: %s", err.Error())
	}
}

func GetCmd() *cobra.Command {
	return rootCmd
}
