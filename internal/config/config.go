package config

import (
	"os"
	"sync"

	"github.com/BurntSushi/toml"

	"github.com/Drafteame/draft/internal/pkg/aws"
	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/pkg/log"
)

const (
	configPath   = "$HOME/.draftea/draft/config.toml"
	ssmParamName = "/service/sentry/dev/SENTRY_TOKEN"
)

type Config struct {
	Sentry Sentry `toml:"sentry"`
}

type Sentry struct {
	Token        string `toml:"token"`
	Organization string `toml:"organization"`
	Team         string `toml:"team"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() (Config, error) {
	var err error

	once.Do(func() {
		cfg, err = load()
		if err != nil {
			return
		}

		loadEnvs(&cfg)
	})

	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func load() (Config, error) {
	if !files.Exists(configPath) {
		return Config{}, nil
	}

	content, err := files.Read(configPath)
	if err != nil {
		return Config{}, err
	}

	c := Config{}
	if errUnm := toml.Unmarshal(content, &c); errUnm != nil {
		return Config{}, errUnm
	}

	return c, nil
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

	if cfg.Sentry.Token == "" {
		ssmToken, err := aws.GetParameter(ssmParamName)
		if err != nil {
			log.Warn("Failed to get Sentry token from SSM:", err)
		} else if ssmToken != "" {
			cfg.Sentry.Token = ssmToken
		}
	}
}
