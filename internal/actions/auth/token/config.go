package token

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/BurntSushi/toml"
)

// Config contains public OAuth client settings, never credentials or tokens.
type Config struct {
	ClientID string `toml:"client_id"`
	Provider string `toml:"provider"`
	Port     int    `toml:"port"`
	AuthURL  string `toml:"auth_url"`
	TokenURL string `toml:"token_url"`
}

type settings struct {
	ClientID string `toml:"client_id"`
	Provider string `toml:"provider"`
	Port     *int   `toml:"port"`
	AuthURL  string `toml:"auth_url"`
	TokenURL string `toml:"token_url"`
}

// LoadConfig loads only auth settings; it does not initialize Sentry or AWS.
func LoadConfig(environment string) (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("resolve home directory: %w", err)
	}
	return loadConfig(environment, filepath.Join(home, ".draftea", "draft", "config.toml"), os.Getenv)
}

func loadConfig(environment, path string, getenv func(string) string) (Config, error) {
	ids := map[string]string{
		"dev":  "client_01KY2KVCJDRJ7T69XMQDZ1G7PA",
		"prod": "client_01KY2KVDCTFB1RCAWEHT4G67X7",
	}
	id, ok := ids[environment]
	if !ok {
		return Config{}, fmt.Errorf("environment must be dev or prod")
	}
	cfg := Config{ClientID: id, Provider: "GoogleOAuth", Port: 3001,
		AuthURL: "https://api.workos.com/user_management/authorize", TokenURL: "https://api.workos.com/user_management/authenticate"}
	var file struct {
		Auth struct {
			WorkOS map[string]settings `toml:"workos"`
		} `toml:"auth"`
	}
	content, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("read auth configuration: %w", err)
	}
	if err == nil {
		if _, err := toml.Decode(string(content), &file); err != nil {
			return Config{}, fmt.Errorf("parse auth configuration: %w", err)
		}
		s := file.Auth.WorkOS[environment]
		if s.ClientID != "" {
			cfg.ClientID = s.ClientID
		}
		if s.Provider != "" {
			cfg.Provider = s.Provider
		}
		if s.Port != nil {
			cfg.Port = *s.Port
		}
		if s.AuthURL != "" {
			cfg.AuthURL = s.AuthURL
		}
		if s.TokenURL != "" {
			cfg.TokenURL = s.TokenURL
		}
	}
	for key, dest := range map[string]*string{
		"DRAFT_WORKOS_CLIENT_ID": &cfg.ClientID, "DRAFT_WORKOS_PROVIDER": &cfg.Provider,
		"DRAFT_WORKOS_AUTH_URL": &cfg.AuthURL, "DRAFT_WORKOS_TOKEN_URL": &cfg.TokenURL,
	} {
		if value := getenv(key); value != "" {
			*dest = value
		}
	}
	if value := getenv("DRAFT_WORKOS_PORT"); value != "" {
		port, err := strconv.Atoi(value)
		if err != nil {
			return Config{}, fmt.Errorf("DRAFT_WORKOS_PORT must be an integer")
		}
		cfg.Port = port
	}
	return cfg, nil
}

// Validate rejects invalid settings before opening a browser or a listener.
func (c Config) Validate() error {
	if c.ClientID == "" {
		return fmt.Errorf("WorkOS client ID is required")
	}
	if c.Provider == "" {
		return fmt.Errorf("WorkOS provider is required")
	}
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("callback port must be between 1 and 65535")
	}
	for _, endpoint := range []string{c.AuthURL, c.TokenURL} {
		u, err := url.Parse(endpoint)
		if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.Fragment != "" {
			return fmt.Errorf("WorkOS endpoints must be absolute HTTPS URLs without user info or fragments")
		}
	}
	return nil
}
