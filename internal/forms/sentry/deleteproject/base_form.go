package deleteproject

import (
	"errors"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
	"github.com/Drafteame/draft/internal/pkg/sentry"
)

func baseForm(input *dtos.DeleteProjectInput) error {
	projects, err := sentry.ListProjects()

	if err != nil {
		return err
	}

	err = inputs.Select[string]("Sentry Project Name:",
		inputs.WithDescription[string]("Select the sentry project you want to delete."),
		inputs.WithValue(&input.ProjectName),
		inputs.WithSaveKey[string](),
		inputs.WithOptions(projects),
	)

	if err != nil {
		return err
	}

	err = inputs.Confirm("Confirmation to delete project: "+input.ProjectName,
		inputs.WithDescription[bool]("Are you sure you want to delete this project?"),
		inputs.WithValue(&input.Confirmation),
	)

	if err != nil {
		return err
	}

	if !input.Confirmation {
		return nil
	}

	err = inputs.Text("Sentry Project ID:",
		inputs.WithDescription[string]("Type project id to delete."),
		inputs.WithValue(&input.ProjectID),
		inputs.WithValidation(func(s string) error {
			if s == "" {
				return errors.New("project id cannot be empty")
			}
			if projects[input.ProjectName] != input.ProjectID {
				return errors.New("project id does not match with project: ")
			}
			return nil
		}),
	)
	if err != nil {
		return err
	}

	return nil
}
