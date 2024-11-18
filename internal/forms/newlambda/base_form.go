package newlambda

import (
	"github.com/charmbracelet/huh"

	"github.com/Drafteame/draft/internal/actions/dtos"
)

func baseForm(input *dtos.Input) error {
	servicePath := huh.NewInput().
		Title("Service Path:").
		Description("Enter the path to the service").
		Value(&input.ServicePath)

	lambdaName := huh.NewInput().
		Title("Lambda Name:").
		Description("Enter the name of the new lambda").
		Value(&input.LambdaName)

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
