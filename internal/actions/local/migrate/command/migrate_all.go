package command

import "fmt"

func (a *Action) migrateAll(config Config) error {
	for db := range config.Migrations.Databases {
		if config.Migrations.Databases[db].Group != a.Input.Group {
			continue
		}

		_, _ = fmt.Printf("Executing migrations for '%s'\n", db)
		if err := a.migrateOne(config, db); err != nil {
			return err
		}
	}

	return nil
}
