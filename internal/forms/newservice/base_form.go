package newservice

import (
	"errors"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
	"github.com/Drafteame/draft/internal/project"
)

func baseForm(input *dtos.ServiceInput) error {
	err := inputs.Text("Service Name:",
		inputs.WithDescription[string]("Set the name of the service to used for configuration."),
		inputs.WithValue(&input.ServiceName),
		inputs.WithValidation(func(s string) error {
			if s == "" {
				return errors.New("service name cannot be empty")
			}

			return nil
		}),
	)

	if err != nil {
		return err
	}

	input.NormalizedServiceName = project.NormalizeServiceName(input.ServiceName)

	err = inputs.Text("Service Path:",
		inputs.WithDescription[string]("Set the path of the service to use for configuration."),
		inputs.WithValue(&input.ServicePath),
		inputs.WithPlaceholder[string](input.NormalizedServiceName),
	)

	if err != nil {
		return err
	}

	if input.ServicePath == "" {
		input.ServicePath = input.NormalizedServiceName
	}

	input.ServicePath = "services/" + input.ServicePath

	return nil
}
