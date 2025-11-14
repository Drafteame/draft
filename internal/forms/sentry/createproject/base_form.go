package createproject

import (
	"errors"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
)

func baseForm(input *dtos.CreateProjectInput) error {
	err := inputs.Text("Sentry Project Name:",
		inputs.WithDescription[string]("Enter the name for the new Sentry project."),
		inputs.WithValue(&input.ProjectName),
		inputs.WithValidation(func(s string) error {
			if s == "" {
				return errors.New("project name cannot be empty")
			}
			return nil
		}),
	)

	if err != nil {
		return err
	}

	return nil
}
