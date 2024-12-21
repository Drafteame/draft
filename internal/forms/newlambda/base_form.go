package newlambda

import (
	"errors"

	"github.com/charmbracelet/huh"

	"github.com/Drafteame/draft/internal/dtos"
)

func baseForm(input *dtos.ServiceInput) error {
	servicePath := huh.NewInput().
		Title("Service Path:").
		Description("Enter the path to the service excluding 'services' folder.").
		Value(&input.ServicePath).
		Validate(func(s string) error {
			if s == "" {
				return errors.New("service path cannot be empty")
			}

			return nil
		})

	lambdaName := huh.NewInput().
		Title("Lambda Name:").
		Description("Enter the name of the new lambda.").
		Value(&input.LambdaName).
		Validate(func(s string) error {
			if s == "" {
				return errors.New("lambda name cannot be empty")
			}

			return nil
		})

	lambdaType := huh.NewSelect[string]().
		Title("Select lambda type:").
		Options(
			huh.NewOption("Plain", "plain"),
			huh.NewOption("SQS", "sqs"),
			huh.NewOption("SNS+SQS", "snssqs"),
			huh.NewOption("HTTP", "http"),
			huh.NewOption("Cron", "cron"),
		).
		Value(&input.LambdaType)

	group := huh.NewGroup(
		servicePath,
		lambdaName,
		lambdaType,
	).Title("Lambda Details")

	return huh.NewForm(group).WithTheme(huh.ThemeCharm()).Run()
}
