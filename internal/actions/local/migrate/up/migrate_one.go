package up

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

	return gomigrate.Exec(gomigrate.ActionUp, gomigrate.Config{
		Source:   source,
		Database: database,
	})
}
