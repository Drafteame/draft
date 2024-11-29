package config

import (
	"os"

	"github.com/BurntSushi/toml"

	"github.com/Drafteame/draft/internal/files"
)

const configPath = "$HOME/.draftea/draft/config.toml"

type Config struct {
	Sentry Sentry `toml:"sentry"`
}

type Sentry struct {
	Token        string `toml:"token"`
	Organization string `toml:"organization"`
	Team         string `toml:"team"`
}

func Get() Config {
	cfg, err := load()
	if err != nil {
		panic(err)
	}

	loadEnvs(&cfg)

	return cfg
}

func load() (Config, error) {
	if !files.Exists(configPath) {
		return Config{}, nil
	}

	content, err := files.Read(configPath)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	if errUnm := toml.Unmarshal(content, &cfg); errUnm != nil {
		return Config{}, errUnm
	}

	return cfg, nil
}

func loadEnvs(cfg *Config) {
	loadSentryEnvs(cfg)
}

func loadSentryEnvs(cfg *Config) {
	token := os.Getenv("DRAFT_SENTRY_TOKEN")
	if token != "" {
		cfg.Sentry.Token = token
	}

	team := os.Getenv("DRAFT_SENTRY_TEAM")
	if team != "" {
		cfg.Sentry.Team = team
	}

	org := os.Getenv("DRAFT_SENTRY_ORGANIZATION")
	if org != "" {
		cfg.Sentry.Organization = org
	}
}
