package migrateconfig

import (
	"fmt"
	"os"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"

	"github.com/Drafteame/draft/internal/pkg/files"
)

const (
	migrateConfigFile = ".local-migrate-config.yml"
)

type MigrateConfig struct {
	Migrations Migrations `yaml:"migrations"`
}

type Migrations struct {
	BasePath  string              `yaml:"base_path"`
	Databases map[string]DBConfig `yaml:"databases"`
}

type DBConfig struct {
	Group      string         `yaml:"group"`
	Folder     string         `yaml:"folder"`
	Connection map[string]any `yaml:"connection"`
}

// GetAvailableDatabases reads the .local-migrate-config.yml file and returns
// a list of available databases that are not in the test group
func GetAvailableDatabases(workingDir string) (map[string]string, error) {
	configPath := fmt.Sprintf("%s/%s", workingDir, migrateConfigFile)

	if !files.Exists(configPath) {
		return nil, fmt.Errorf("migrate config file not found: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrate config: %w", err)
	}

	var config MigrateConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse migrate config: %w", err)
	}

	databases := make(map[string]string)

	for dbName, dbConfig := range config.Migrations.Databases {
		// Skip test databases
		if dbConfig.Group == "test" {
			continue
		}

		// Use the database name as display label
		displayName := formatDisplayName(dbName)
		databases[displayName] = dbName
	}

	if len(databases) == 0 {
		return nil, fmt.Errorf("no non-test databases found in migrate config")
	}

	return databases, nil
}

// ToPascalCase converts a snake_case string to PascalCase
// Examples:
//   - general -> General
//   - user_preferences -> UserPreferences
//   - games_core -> GamesCore
func ToPascalCase(s string) string {
	if s == "" {
		return ""
	}

	caser := cases.Title(language.English)
	parts := strings.Split(s, "_")

	pascalParts := lo.Map(parts, func(part string, _ int) string {
		return caser.String(part)
	})

	return strings.Join(pascalParts, "")
}

// formatDisplayName formats a database name for display in the selection list
// Examples:
//   - general -> General
//   - user_preferences -> User Preferences
//   - games_core -> Games Core
func formatDisplayName(dbName string) string {
	caser := cases.Title(language.English)
	parts := strings.Split(dbName, "_")

	titleParts := lo.Map(parts, func(part string, _ int) string {
		return caser.String(part)
	})

	return strings.Join(titleParts, " ")
}
