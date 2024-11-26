package newlambda

import (
	"errors"

	"github.com/charmbracelet/huh"

	"github.com/Drafteame/draft/internal/actions/dtos"
)

func queueForm(input *dtos.Input) error {
	queueArn := huh.NewInput().
		Title("Set Queue ARN:").
		Description("Enter the ARN of the queue or some replacer that works on lambda-config.yml").
		Value(&input.QueueARN).
		Validate(func(s string) error {
			if s == "" {
				return errors.New("queue ARN cannot be empty")
			}

			return nil
		})

	return huh.NewForm(huh.NewGroup(queueArn)).WithTheme(huh.ThemeCharm()).Run()
}
