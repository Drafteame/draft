package newlambda

import (
	"errors"
	"github.com/Drafteame/draft/internal/dtos"

	"github.com/charmbracelet/huh"
)

func cronForm(input *dtos.ServiceInput) error {
	cronExpression := huh.NewInput().
		Title("Set Cron Expression:").
		Description("Enter the cron expression").
		Value(&input.CronExpression).
		Validate(func(s string) error {
			if s == "" {
				return errors.New("cron expression cannot be empty")
			}

			return nil
		})

	return huh.NewForm(huh.NewGroup(cronExpression)).WithTheme(huh.ThemeCharm()).Run()
}
