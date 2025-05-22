package command

import (
	"fmt"
	"os"
)

func (a *Action) Exec() error {
	if a.Input.WorkingDir != "" {
		if err := os.Chdir(a.Input.WorkingDir); err != nil {
			panic(err)
		}
	}

	config, err := a.loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load local migrations config: %w", err)
	}

	if a.Input.All {
		if err := a.migrateAll(config); err != nil {
			return fmt.Errorf("failed to execute '%s' migration on all databases: %w", a.Input.Command, err)
		}

		_, _ = fmt.Printf("Migration '%s' executed successfully\n", a.Input.Command)
		return nil
	}

	dbName := a.Input.Database

	if dbName == "" {
		name, err := a.promptSelectDB(config)
		if err != nil {
			return fmt.Errorf("failed to select database: %w", err)
		}

		dbName = name
	}

	if err := a.migrateOne(config, dbName); err != nil {
		return fmt.Errorf("failed to execute '%s' migration on database %s: %w", a.Input.Command, dbName, err)
	}

	_, _ = fmt.Printf("Migration '%s' executed successfully\n", a.Input.Command)
	return nil
}
