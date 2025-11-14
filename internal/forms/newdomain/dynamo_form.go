package newdomain

import (
	"errors"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
)

func dynamoForm(input *dtos.DomainInput) error {
	err := inputs.Text("Table Name:",
		inputs.WithDescription[string]("Enter the name of the DynamoDB table to use on this domain."),
		inputs.WithValue(&input.TableName),
		inputs.WithValidation(func(s string) error {
			if s == "" {
				return errors.New("table name cannot be empty")
			}

			return nil
		}),
	)

	if err != nil {
		return err
	}

	return nil
}
