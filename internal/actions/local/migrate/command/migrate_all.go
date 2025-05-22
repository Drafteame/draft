package command

import "fmt"

func (a *Action) migrateAll(config Config) error {
	for db := range config.Migrations.Databases {
		if config.Migrations.Databases[db].Group != a.Input.Group {
			continue
		}

		_, _ = fmt.Printf("Executing '%s' migration for '%s'\n", a.Input.Command, db)
		if err := a.migrateOne(config, db); err != nil {
			return err
		}
	}

	return nil
}
