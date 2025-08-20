package config

import (
	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/config"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get configuration",
	Long:  "Get current configuration file and print it",
	Run:   run,
}

func run(_ *cobra.Command, _ []string) {
	cfg, err := config.Get()
	if err != nil {
		log.Exitf(1, "failed to obtain config: %s", err.Error())
	}

	tb, err := toml.Marshal(cfg)
	if err != nil {
		log.Exitf(1, "failed to obtain config: %s", err.Error())
	}

	log.PrintCode(string(tb), "toml")
}

func GetCmd() *cobra.Command {
	return configCmd
}
