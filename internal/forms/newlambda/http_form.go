package newlambda

import (
	"errors"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
)

func httpForm(input *dtos.ServiceInput) error {
	err := inputs.Select[string]("HTTP Method:",
		inputs.WithDescription[string]("Select the HTTP method to be used for the endpoint"),
		inputs.WithValue(&input.HTTPMethod),
		inputs.WithOptions(map[string]string{
			"GET":    "GET",
			"POST":   "POST",
			"PUT":    "PUT",
			"DELETE": "DELETE",
		}),
	)

	if err != nil {
		return err
	}

	err = inputs.Text("HTTP Path:",
		inputs.WithDescription[string]("Enter the path to the service"),
		inputs.WithValue(&input.HTTPPath),
		inputs.WithValidation(func(s string) error {
			if s == "" {
				return errors.New("path cannot be empty")
			}

			return nil
		}),
	)

	return err
}
