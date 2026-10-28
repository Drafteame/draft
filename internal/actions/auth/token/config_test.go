package token

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigurationPrecedenceAndEnvironmentIsolation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(`[sentry]
token = "unrelated-secret"
[auth.workos.dev]
client_id = "file-dev"
port = 3002
provider = "authkit"
[auth.workos.prod]
client_id = "file-prod"
`), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := loadConfig("dev", path, func(key string) string {
		if key == "DRAFT_WORKOS_CLIENT_ID" {
			return "env-dev"
		}
		return ""
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ClientID != "env-dev" || cfg.Port != 3002 || cfg.Provider != "authkit" {
		t.Fatalf("incorrect merged config: %+v", cfg)
	}
	prod, err := loadConfig("prod", path, func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	if prod.ClientID != "file-prod" || prod.Port != 3001 || prod.Provider != "GoogleOAuth" {
		t.Fatalf("environment settings leaked: %+v", prod)
	}
}

func TestConfigDefaultsAndInvalidInputs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.toml")
	for _, environment := range []string{"dev", "prod"} {
		cfg, err := loadConfig(environment, path, func(string) string { return "" })
		if err != nil {
			t.Fatal(err)
		}
		if err := cfg.Validate(); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := loadConfig("staging", path, func(string) string { return "" }); err == nil {
		t.Fatal("unsupported environment accepted")
	}
	if _, err := loadConfig("dev", path, func(k string) string {
		if k == "DRAFT_WORKOS_PORT" {
			return "wrong"
		}
		return ""
	}); err == nil {
		t.Fatal("invalid port environment value accepted")
	}
	cfg, _ := loadConfig("dev", path, func(string) string { return "" })
	for _, port := range []int{0, -1, 65536} {
		c := cfg
		c.Port = port
		if c.Validate() == nil {
			t.Errorf("accepted port %d", port)
		}
	}
	for _, endpoint := range []string{"http://api.workos.com/token", "/relative", "https://user:password@example.com/token", "https://example.com/token#fragment"} {
		c := cfg
		c.TokenURL = endpoint
		if c.Validate() == nil {
			t.Errorf("accepted invalid endpoint")
		}
	}
}

func TestConfigReadAndParseErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.toml")
	if err := os.WriteFile(path, []byte("[invalid"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadConfig("dev", path, func(string) string { return "" }); err == nil {
		t.Fatal("invalid TOML accepted")
	}
	if _, err := loadConfig("dev", filepath.Dir(path), func(string) string { return "" }); err == nil {
		t.Fatal("unreadable configuration accepted")
	}
}
