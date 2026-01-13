package newdomain

import (
	"errors"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
)

func dynamoForm(input *dtos.DomainInput) error {
	// Prompt for table name only if not provided via flag
	if input.TableName == "" {
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
	}

	// Validate table name if provided via flag
	if input.TableName == "" {
		return errors.New("table name cannot be empty")
	}

	return nil
}
