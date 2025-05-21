package force

import (
	"fmt"
	"os"

	"github.com/Drafteame/draft/internal/pkg/gomigrate"
)

func (a *Action) migrateOne(config Config, dbName string) error {
	db, ok := config.Migrations.Databases[dbName]
	if !ok {
		return fmt.Errorf("database not found: %s", dbName)
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

	version := fmt.Sprintf("%d", a.Input.Version)

	return gomigrate.Exec(gomigrate.ActionForce, gomigrate.Config{
		Source:   source,
		Database: database,
		Args:     []string{version},
	})
}
