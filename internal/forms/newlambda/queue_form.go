package newlambda

import (
	"github.com/charmbracelet/huh"

	"github.com/Drafteame/draft/internal/actions/dtos"
)

func queueForm(input *dtos.Input) error {
	queueArn := huh.NewInput().
		Title("Set Queue ARN:").
		Description("Enter the ARN of the queue or some replacer that works on lambda-config.yml").
		Value(&input.QueueARN)

	return huh.NewForm(huh.NewGroup(queueArn)).WithTheme(huh.ThemeCharm()).Run()
}
