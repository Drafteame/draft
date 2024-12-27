package newlambda

import (
	"errors"
	"strings"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
)

func httpForm(input *dtos.LambdaInput) error {
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
		inputs.WithDescription[string]("Enter the path to the service. If your path has parameters, use <param> syntax."),
		inputs.WithValue(&input.HTTPPath),
		inputs.WithValidation(func(s string) error {
			if s == "" {
				return errors.New("path cannot be empty")
			}

			return nil
		}),
	)

	apiGatewayPath, echoPath := transformPath(input.HTTPPath)

	input.HTTPPathAPIGateway = apiGatewayPath
	input.HTTPPathEcho = echoPath

	return err
}

func transformPath(path string) (string, string) {
	// Transform for API Gateway
	apiGatewayPath := strings.ReplaceAll(path, "<", "{")
	apiGatewayPath = strings.ReplaceAll(apiGatewayPath, ">", "}")

	// Transform for Echo
	echoPath := strings.ReplaceAll(path, "<", ":")
	echoPath = strings.ReplaceAll(echoPath, ">", "")

	return apiGatewayPath, echoPath
}
