package newlambda

import (
	"errors"
	"strings"

	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
)

func baseForm(input *dtos.LambdaInput) error {
	err := inputs.Text("Service Path:",
		inputs.WithDescription[string]("Enter the path to the service folder."),
		inputs.WithValue(&input.ServicePath),
		inputs.WithValidation(func(val string) error {
			if val == "" {
				return errors.New("service path cannot be empty")
			}

			return nil
		}),
	)

	if err != nil {
		return err
	}

	err = inputs.Text("Lambda Name:",
		inputs.WithDescription[string]("Enter the name of the new lambda."),
		inputs.WithValue(&input.LambdaName),
		inputs.WithValidation(func(s string) error {
			if s == "" {
				return errors.New("lambda name cannot be empty")
			}

			return nil
		}),
	)

	if err != nil {
		return err
	}

	err = inputs.Select[string]("Lambda Type:",
		inputs.WithDescription[string]("Select the type of the new lambda to be created."),
		inputs.WithValue(&input.LambdaType),
		inputs.WithOptions(map[string]string{
			"Plain":   "plain",
			"SQS":     "sqs",
			"SNS+SQS": "snssqs",
			"HTTP":    "http",
			"Cron":    "cron",
		}),
	)

	input.PackageName = data.Meta.PackageName

	if !strings.HasPrefix(input.ServicePath, "services/") {
		input.ServicePath = "services/" + input.ServicePath
	}

	return err
}
