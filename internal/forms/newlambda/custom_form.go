package newlambda

import (
	"errors"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
)

func customForm(input *dtos.LambdaInput) error {
	err := inputs.Text("Custom Type Path:",
		inputs.WithDescription[string]("Enter the custom type path (e.g., 'custom', 'events', 'triggers')."),
		inputs.WithValue(&input.CustomTypePath),
		inputs.WithPlaceholder[string]("custom"),
		inputs.WithValidation(func(s string) error {
			if s == "" {
				return errors.New("custom type path cannot be empty")
			}

			return nil
		}),
	)

	if err != nil {
		return err
	}

	// Set default value if empty
	if input.CustomTypePath == "" {
		input.CustomTypePath = "custom"
	}

	err = inputs.Confirm("Use Idempotency:",
		inputs.WithDescription[bool]("Should this lambda implement idempotency functions?"),
		inputs.WithValue(&input.UseIdempotency),
	)

	if err != nil {
		return err
	}

	return nil
}
