package newdomain

import (
	"errors"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
)

func postgresForm(input *dtos.DomainInput) error {
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

	err = inputs.Text("Set ID Prefix:",
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

	err = inputs.Select[string]("Select an available database to connect:",
		inputs.WithDescription[string]("Select the database that should be connected on the domain"),
		inputs.WithValue(&input.DBName),
		inputs.WithOptions(map[string]string{
			"Audiences":           "Audiences",
			"Data Products":       "DataProducts",
			"Fraud":               "fraud",
			"Games Core":          "GamesCore",
			"General":             "General",
			"Kyc":                 "Kyc",
			"Notification Engine": "NotificationEngine",
			"Scores":              "Scores",
			"Stats":               "Stats",
			"Turbo":               "Turbo",
			"User Preferences":    "UserPreferences",
		}),
	)

	if err != nil {
		return err
	}

	input.DBProviderFuncName = input.DBName

	return nil
}
