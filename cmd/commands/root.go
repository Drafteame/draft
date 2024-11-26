package commands

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/flags"
)

var version = "unversioned"

var rootCmd = &cobra.Command{
	Use:     "draft <command>",
	Example: "draft new:service",
	Version: version,
	Run:     run,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flags.Flags.WorkingDir, "working-dir", "w", "", "Working directory")
}

func run(cmd *cobra.Command, _ []string) {
	if err := cmd.Help(); err != nil {
		panic(err)
	}
}

func GetCmd() *cobra.Command {
	return rootCmd
}
