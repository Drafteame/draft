package newservice

import (
	"errors"

	"github.com/charmbracelet/huh"

	"github.com/Drafteame/draft/internal/actions/dtos"
)

func GetForm(input *dtos.Input) error {
	if err := newServiceDetails(input); err != nil {
		return err
	}

	return newServiceFrameDetails(input)
}

func newServiceDetails(input *dtos.Input) error {
	serviceTypes := huh.NewSelect[string]().
		Title("Select service framework:").
		Options(
			huh.NewOption("Serverless", "sls"),
			// huh.NewOption("CDK", "cdk"),
		).
		Value(&input.ServiceFramework)

	serviceName := huh.NewInput().
		Title("Service Name:").
		Value(&input.ServiceName).
		Validate(func(s string) error {
			if s == "" {
				return errors.New("service name cannot be empty")
			}

			return nil
		})

	ServicePath := huh.NewInput().
		Title("Service Folder:").
		Description("Enter the folder name where should be placed the service content").
		Value(&input.ServicePath)

	frameworkVersion := huh.NewSelect[string]().
		Title("Select Framework version:").
		Options(
			huh.NewOption("Framework V2", "v2"),
			// huh.NewOption("Experimental Engine", "exp-engine"),
		).
		Value(&input.FrameVersion)

	group1 := huh.NewGroup(
		serviceTypes,
		serviceName,
		ServicePath,
		frameworkVersion,
	).Title("Service Details")

	return huh.NewForm(group1).WithTheme(huh.ThemeCharm()).Run()
}

func newServiceFrameDetails(input *dtos.Input) error {
	inputs := make([]huh.Field, 0)

	if input.ServiceFramework == "sls" {
		warmupEnabled := huh.NewConfirm().
			Title("Enable warmup?").
			Value(&input.WarmupEnabled)

		customDomain := huh.NewConfirm().
			Title("Configure custom domain?").
			Value(&input.CustomDomain)

		inputs = append(inputs, warmupEnabled, customDomain)
	}

	group := huh.NewGroup(inputs...).Title("Service Configs")

	if err := huh.NewForm(group).WithTheme(huh.ThemeCharm()).Run(); err != nil {
		return err
	}

	if input.CustomDomain {
		domainPath := huh.NewInput().
			Title("Specify domain base path:").
			Value(&input.DomainPath).
			Validate(func(s string) error {
				if s == "" {
					return errors.New("domain path cannot be empty")
				}

				return nil
			})

		if err := huh.NewForm(huh.NewGroup(domainPath)).WithTheme(huh.ThemeCharm()).Run(); err != nil {
			return err
		}
	}

	return nil
}
