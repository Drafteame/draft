package command

import (
	"fmt"
)

func (a *Action) Exec() error {
	config, err := a.loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load local migrations config: %w", err)
	}

	if a.Input.All {
		if errAll := a.migrateAll(config); errAll != nil {
			return fmt.Errorf("failed to execute '%s' migration on all databases: %w", a.Input.Command, errAll)
		}

		return nil
	}

	dbName := a.Input.Database

	if dbName == "" {
		name, errSelect := a.promptSelectDB(config)
		if errSelect != nil {
			return fmt.Errorf("failed to select database: %w", errSelect)
		}

		dbName = name
	}

	if errMigrate := a.migrateOne(config, dbName); errMigrate != nil {
		return fmt.Errorf("failed to execute '%s' migration on database %s: %w", a.Input.Command, dbName, errMigrate)
	}

	return nil
}
