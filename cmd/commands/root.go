package commands

import (
	"github.com/spf13/cobra"
)

const VERSION = "unversioned"

var rootCmd = &cobra.Command{
	Use:     "draft <command>",
	Example: "draft new:service",
	Version: VERSION,
	Run:     run,
}

func run(cmd *cobra.Command, _ []string) {
	if err := cmd.Help(); err != nil {
		panic(err)
	}
}

func GetCmd() *cobra.Command {
	return rootCmd
}
