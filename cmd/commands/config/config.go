package config

import (
	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get configuration",
	Long:  "Get current configuration file and print it",
	Run:   run,
}

func run(_ *cobra.Command, _ []string) {
	cfg := config.Get()

	tb, err := toml.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	println(string(tb))
}

func GetCmd() *cobra.Command {
	return configCmd
}
