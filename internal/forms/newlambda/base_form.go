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
	if !input.IsLegacy {
		services, errGet := project.GetServices()
		if errGet != nil {
			return errGet
		}

		errSrv := inputs.Select[string]("Service:",
			inputs.WithDescription[string]("Select the service where the new lambda will be created."),
			inputs.WithValue(&input.ServicePath),
			inputs.WithOptions(services),
		)

		if errSrv != nil {
			return errSrv
		}
	}

	errName := inputs.Text("Lambda Name:",
		inputs.WithDescription[string]("Enter the name of the new lambda."),
		inputs.WithValue(&input.LambdaName),
		inputs.WithValidation(func(s string) error {
			if s == "" {
				return errors.New("lambda name cannot be empty")
			}

			return nil
		}),
	)

	if errName != nil {
		return errName
	}

	errType := inputs.Select[string]("Lambda Type:",
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

	if errType != nil {
		return errType
	}

	if input.LambdaType != "cron" && input.LambdaType != "http" {
		errConcurrency := inputs.Select[string]("Reserved Concurrency:",
			inputs.WithDescription[string]("Select the tier of reserved concurrency for the new lambda."),
			inputs.WithValue(&input.ReservedConcurrency),
			inputs.WithOptions(map[string]string{
				"MACRO - HTTP":          "macro.http",
				"MACRO - EVENT DRIVEN":  "macro.eventDriven",
				"HIGH - HTTP":           "high.http",
				"HIGH - EVENT DRIVEN":   "high.eventDriven",
				"MEDIUM - HTTP":         "medium.http",
				"MEDIUM - EVENT DRIVEN": "medium.eventDriven",
				"LOW - HTTP":            "low.http",
				"LOW - EVENT DRIVEN":    "low.eventDriven",
				"MIN - HTTP":            "min.http",
				"MIN - EVENT DRIVEN":    "min.eventDriven",
				"MICRO - HTTP":          "micro.http",
				"MICRO - EVENT DRIVEN":  "micro.eventDriven",
			}),
		)

		if errConcurrency != nil {
			return errConcurrency
		}
	}

	input.PackageName = data.Meta.PackageName

	if !strings.HasPrefix(input.ServicePath, "services/") && !input.IsLegacy {
		input.ServicePath = "services/" + input.ServicePath
	}

	splitServicePath := strings.Split(input.ServicePath, "/")
	serviceName := splitServicePath[len(splitServicePath)-1]
	input.ServiceName = strings.ToLower(serviceName)

	return nil
}
