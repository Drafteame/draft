package newdomain

import (
	"errors"
	"fmt"

	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
	"github.com/Drafteame/draft/internal/pkg/migrateconfig"
)

func postgresForm(input *dtos.DomainInput) error {
	// Prompt for table name only if not provided via flag
	if input.TableName == "" {
		err := inputs.Text("Table Name:",
			inputs.WithDescription[string]("Enter the name of the table to use on this domain (you can specify schema to by using 'schema.table' notation)."),
			inputs.WithValue(&input.TableName),
			inputs.WithValidation(func(s string) error {
				if s == "" {
					return errors.New("table name cannot be empty")
				}

				return nil
			}),
		)

		if err != nil {
			return err
		}
	}

	// Validate table name if provided via flag
	if input.TableName == "" {
		return errors.New("table name cannot be empty")
	}

	// Prompt for DB prefix only if not provided via flag
	if input.DBPrefix == "" {
		err := inputs.Text("Set ID Prefix:",
			inputs.WithDescription[string]("Enter the prefix to use for values on the ID field (length should be 3 chars)"),
			inputs.WithValue(&input.DBPrefix),
			inputs.WithValidation(func(s string) error {
				if len(s) != 3 {
					return errors.New("prefix should be 3 characters long")
				}

				return nil
			}),
		)

		if err != nil {
			return err
		}
	}

	// Validate DB prefix if provided via flag
	if len(input.DBPrefix) != 3 {
		return errors.New("db-prefix should be 3 characters long")
	}

	// Load available databases from .local-migrate-config.yml
	workingDir := data.Flags.WorkingDir
	if workingDir == "" {
		workingDir = "."
	}

	databases, err := migrateconfig.GetAvailableDatabases(workingDir)
	if err != nil {
		return fmt.Errorf("failed to load available databases: %w", err)
	}

	// Prompt for DB name only if not provided via flag
	if input.DBName == "" {
		err := inputs.Select[string]("Select an available database to connect:",
			inputs.WithDescription[string]("Select the database that should be connected on the domain"),
			inputs.WithValue(&input.DBName),
			inputs.WithOptions(databases),
		)

		if err != nil {
			return err
		}
	}

	// Validate DB name if provided via flag
	if input.DBName != "" {
		validDBName := false
		for _, dbName := range databases {
			if dbName == input.DBName {
				validDBName = true
				break
			}
		}

		if !validDBName {
			availableDBs := make([]string, 0, len(databases))
			for _, dbName := range databases {
				availableDBs = append(availableDBs, dbName)
			}
			return fmt.Errorf("invalid db-name '%s': must be one of %v", input.DBName, availableDBs)
		}
	}

	// Convert database name to PascalCase for provider function name
	input.DBProviderFuncName = migrateconfig.ToPascalCase(input.DBName)

	return nil
}
