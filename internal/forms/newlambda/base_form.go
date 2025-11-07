package newlambda

import (
	"errors"
	"strings"

	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
	"github.com/Drafteame/draft/internal/project"
)

var reservedConcurrency = map[string]map[string]string{
	"http": {
		"MACRO":  "macro.http",
		"HIGH":   "high.http",
		"MEDIUM": "medium.http",
		"LOW":    "low.http",
		"MIN":    "min.http",
		"MICRO":  "micro.http",
	},
	"eventDriven": {
		"MACRO":  "macro.eventDriven",
		"HIGH":   "high.eventDriven",
		"MEDIUM": "medium.eventDriven",
		"LOW":    "low.eventDriven",
		"MIN":    "min.eventDriven",
		"MICRO":  "micro.eventDriven",
	},
}

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

	rcType := "eventDriven"
	if input.LambdaType == "http" {
		rcType = "http"
	}

	errConcurrency := inputs.Select[string]("Reserved Concurrency:",
		inputs.WithDescription[string]("Select the tier of reserved concurrency for the new lambda."),
		inputs.WithValue(&input.ReservedConcurrency),
		inputs.WithOptions(reservedConcurrency[rcType]),
	)

	if errConcurrency != nil {
		return errConcurrency
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
