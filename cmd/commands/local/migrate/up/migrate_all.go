package up

import "fmt"

func migrateAll(config Config) error {
	for db := range config.Migrations.Databases {
		_, _ = fmt.Printf("Executing migrations for '%s'\n", db)
		if err := migrateOne(config, db); err != nil {
			return err
		}
	}

	return nil
}
