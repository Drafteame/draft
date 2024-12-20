package commands

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/data"
)

var version = "unversioned"

var rootCmd = &cobra.Command{
	Use:     "draft <command>",
	Example: "draft new:service",
	Version: version,
	Run:     run,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&data.Flags.WorkingDir, "working-dir", "w", "", "Working directory")
	rootCmd.PersistentFlags().BoolVarP(&data.Flags.Debug, "debug", "d", false, "Debug mode")
	rootCmd.PersistentFlags().BoolVarP(&data.Flags.TTY, "tty", "t", true, "TTY mode")
}

func run(cmd *cobra.Command, _ []string) {
	if err := cmd.Help(); err != nil {
		panic(err)
	}
}

func GetCmd() *cobra.Command {
	return rootCmd
}
