package up

import (
	"errors"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/pkg/inputs"
)

type Config struct {
	Migrations Migrations `yaml:"migrations"`
}

type Migrations struct {
	BasePath  string              `yaml:"base_path"`
	Databases map[string]Database `yaml:"databases"`
}

type Database struct {
	Group      string     `yaml:"group"`
	Folder     string     `yaml:"folder"`
	Connection Connection `yaml:"connection"`
}

type Connection struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func (a *Action) loadConfig() (Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}

	configPath := cwd + "/" + a.Input.LocalMigrateConfig

	config := Config{}

	if err := files.LoadYAML(configPath, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (a *Action) promptSelectDB(config Config) (string, error) {
	dbs := make(map[string]string)

	for key := range config.Migrations.Databases {
		if config.Migrations.Databases[key].Group != a.Input.Group {
			continue
		}

		title := strings.Join(strings.Split(key, "_"), " ")
		dbs[cases.Title(language.English).String(title)] = key
	}

	if len(dbs) == 0 {
		return "", errors.New("no databases found")
	}

	var db string

	errSelect := inputs.Select[string]("Select DB to migrate:",
		inputs.WithDescription[string]("Select the database you want to migrate defined on the config file"),
		inputs.WithOptions(dbs),
		inputs.WithValue(&db),
	)

	return db, errSelect
}
