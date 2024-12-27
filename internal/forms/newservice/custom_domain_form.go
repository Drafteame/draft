package newservice

import (
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
)

func customDomain(input *dtos.ServiceInput) error {
	err := inputs.Text("Specify domain base path:",
		inputs.WithDescription[string]("Set the base path that should respond to service endpoints (default should be service name)."),
		inputs.WithPlaceholder[string](input.NormalizedServiceName),
		inputs.WithValue(&input.DomainPath),
	)

	if err != nil {
		return err
	}

	if input.DomainPath == "" {
		input.DomainPath = input.NormalizedServiceName
	}

	return nil
}
