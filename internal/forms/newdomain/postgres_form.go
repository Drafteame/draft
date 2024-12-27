package newdomain

import (
	"errors"

	"github.com/Drafteame/draft/internal/dtos"
	inputs2 "github.com/Drafteame/draft/internal/pkg/inputs"
)

func postgresForm(input *dtos.DomainInput) error {
	err := inputs2.Text("Table Name:",
		inputs2.WithDescription[string]("Enter the name of the table to use on this domain (you can specify schema to by using 'schema.table' notation)."),
		inputs2.WithValue(&input.TableName),
		inputs2.WithValidation(func(s string) error {
			if s == "" {
				return errors.New("table name cannot be empty")
			}

			return nil
		}),
	)

	if err != nil {
		return err
	}

	err = inputs2.Text("Set ID Prefix:",
		inputs2.WithDescription[string]("Enter the prefix to use for values on the ID field (length should be 3 chars)"),
		inputs2.WithValue(&input.DBPrefix),
		inputs2.WithValidation(func(s string) error {
			if len(s) != 3 {
				return errors.New("prefix should be 3 characters long")
			}

			return nil
		}),
	)

	if err != nil {
		return err
	}

	err = inputs2.Select[string]("Select an available database to connect:",
		inputs2.WithDescription[string]("Select the database that should be connected on the domain"),
		inputs2.WithValue(&input.DBName),
		inputs2.WithOptions(map[string]string{
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
