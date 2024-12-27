package newservice

import (
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
)

func frameDetails(input *dtos.ServiceInput) error {
	err := inputs.Confirm("Enable Warmup?:",
		inputs.WithDescription[bool]("Select if you want to enable warmup plugin for all service"),
		inputs.WithValue(&input.WarmupEnabled),
	)

	if err != nil {
		return err
	}

	err = inputs.Confirm("Configure custom domain?:",
		inputs.WithDescription[bool]("Select if you want to configure custom domain"),
		inputs.WithValue(&input.CustomDomain),
	)

	if err != nil {
		return err
	}

	if input.CustomDomain {
		return customDomain(input)
	}

	return nil
}
