package newdomain

import (
	"errors"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/charmbracelet/huh"
)

func postgresForm(input *dtos.DomainInput) error {
	err := huh.NewInput().
		Title("Table Name:").
		Description("Enter the name of the table to use on this domain (you can specify schema to by using 'schema.table' notation).").
		Value(&input.TableName).
		Validate(func(s string) error {
			if s == "" {
				return errors.New("table name cannot be empty")
			}

			return nil
		}).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		return err
	}

	err = huh.NewInput().
		Title("Set ID prefix:").
		Description("Enter the prefix to use for values on the ID field (length should be 3 chars)").
		Value(&input.DBPrefix).
		Validate(func(s string) error {
			if len(s) != 3 {
				return errors.New("prefix should be 3 characters long")
			}

			return nil
		}).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		return err
	}

	err = huh.NewSelect[string]().
		Title("Select an available database to connect:").
		Description("Select the database that should be connected on the domain").
		Value(&input.DBName).
		Options(
			huh.NewOption("General", "General"),
			huh.NewOption("Turbo", "turbo"),
			huh.NewOption("Kyc", "Kyc"),
			huh.NewOption("Fraud", "Fraud"),
			huh.NewOption("Audiences", "Audiences"),
			huh.NewOption("Notification Engine", "NotificationEngine"),
			huh.NewOption("User Audiences", "UserAudiences"),
		).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		return err
	}

	input.DBProviderFuncName = input.DBName

	return nil
}
