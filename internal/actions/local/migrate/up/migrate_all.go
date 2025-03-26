package up

import "fmt"

func (a *Action) migrateAll(config Config) error {
	for db := range config.Migrations.Databases {
		_, _ = fmt.Printf("Executing migrations for '%s'\n", db)
		if err := a.migrateOne(config, db); err != nil {
			return err
		}
	}

	return nil
}
