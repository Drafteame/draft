package newlambda

import (
	"errors"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
)

func queueForm(input *dtos.ServiceInput) error {
	return inputs.Text("Queue ARN:",
		inputs.WithDescription[string]("Enter the ARN of the queue or some replacer that works on lambda-config.yml"),
		inputs.WithValue(&input.QueueARN),
		inputs.WithValidation(func(s string) error {
			if s == "" {
				return errors.New("queue ARN cannot be empty")
			}

			return nil
		}),
	)
}
