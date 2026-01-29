package newservice

import (
	"errors"
	"strings"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
	"github.com/Drafteame/draft/internal/project"
)

func baseForm(input *dtos.ServiceInput) error {
	if !input.IsLegacy {
		if err := promptServiceName(input); err != nil {
			return err
		}

		if err := promptRoleName(input); err != nil {
			return err
		}

		if err := promptServicePath(input); err != nil {
			return err
		}

		input.ServicePackage = project.NormalizeServicePackage(input.ServiceName)

		setDefaultServicePath(input)

		return nil
	}

	normalizeLegacyService(input)

	return nil
}

func promptServiceName(input *dtos.ServiceInput) error {
	return inputs.Text("Service Name:",
		inputs.WithDescription[string]("Set the name of the service to used for configuration"),
		inputs.WithValue(&input.ServiceName),
		inputs.WithValidation(func(s string) error {
			if s == "" {
				return errors.New("service name cannot be empty")
			}
			return nil
		}),
	)
}

func promptRoleName(input *dtos.ServiceInput) error {
	return inputs.Text("Role Name:",
		inputs.WithDescription[string]("Set the name of the role to be used by the service (must be in PascalCase, e.g., GameStats, UserTracking, GameEngine)"),
		inputs.WithValue(&input.RoleName),
		inputs.WithValidation(func(s string) error {
			if s == "" {
				return errors.New("role name cannot be empty")
			}
			return nil
		}),
	)
}

func promptServicePath(input *dtos.ServiceInput) error {
	input.NormalizedServiceName = project.NormalizeServiceName(input.ServiceName)
	return inputs.Text("Service Path:",
		inputs.WithDescription[string]("Set the path of the service to use for configuration (default: services/)"),
		inputs.WithValue(&input.ServicePath),
		inputs.WithPlaceholder[string](input.NormalizedServiceName),
	)
}

func setDefaultServicePath(input *dtos.ServiceInput) {
	if input.ServicePath == "" {
		input.ServicePath = "services/" + input.NormalizedServiceName
	}
}

func normalizeLegacyService(input *dtos.ServiceInput) {
	splitServicePath := strings.Split(input.ServicePath, "/")
	sizeSplitServicePath := len(splitServicePath)
	serviceName := splitServicePath[sizeSplitServicePath-1]
	input.NormalizedServiceName = project.NormalizeServiceName(serviceName)
	input.ServicePath = strings.Join(splitServicePath[:sizeSplitServicePath-1], "/") + "/" + input.NormalizedServiceName
	input.ServicePackage = project.NormalizeServicePackage(serviceName)
}
