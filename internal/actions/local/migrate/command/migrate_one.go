package command

import (
	"errors"
	"fmt"
	"os"

	"github.com/Drafteame/draft/internal/pkg/gomigrate"
)

func (a *Action) migrateOne(config Config, dbName string) error {
	if dbName == "all" {
		return nil
	}

	db, ok := config.Migrations.Databases[dbName]
	if !ok {
		return errors.New("database not found")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	source := fmt.Sprintf("%s/%s/%s", cwd, config.Migrations.BasePath, db.Folder)
	database := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.Connection.User,
		db.Connection.Password,
		db.Connection.Host,
		db.Connection.Port,
		db.Connection.Database,
	)

	cfg := gomigrate.Config{
		Source:   source,
		Database: database,
	}

	switch a.Input.Command {
	case "up":
		return gomigrate.Exec(gomigrate.ActionUp, cfg)
	case "force":
		cfg.Args = []string{fmt.Sprintf("%d", a.Input.Version)}
		return gomigrate.Exec(gomigrate.ActionForce, cfg)
	case "down":
		if a.Input.NumberMigrations != 0 {
			cfg.Args = []string{fmt.Sprintf("%d", a.Input.NumberMigrations)}
		} else {
			cfg.Args = []string{"-all"}
		}
		return gomigrate.Exec(gomigrate.ActionDown, cfg)
	default:
		return errors.New("invalid command")
	}
}
