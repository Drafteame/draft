package commands

import (
	"github.com/spf13/cobra"
)

const VERSION = "dev"

var rootCmd = &cobra.Command{
	Use:     "back <command>",
	Example: "back some-command",
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
