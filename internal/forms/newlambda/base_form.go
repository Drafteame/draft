package newlambda

import (
	"errors"
	"strings"

	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
	"github.com/Drafteame/draft/internal/project"
)

func baseForm(input *dtos.LambdaInput) error {
	services, err := project.GetServices()
	if err != nil {
		return err
	}

	err = inputs.Select[string]("Service:",
		inputs.WithDescription[string]("Select the service where the new lambda will be created."),
		inputs.WithValue(&input.ServicePath),
		inputs.WithOptions(services),
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
