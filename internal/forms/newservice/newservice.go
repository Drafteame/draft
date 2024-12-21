package newservice

import (
	"errors"

	"github.com/charmbracelet/huh"

	"github.com/Drafteame/draft/internal/dtos"
)

func GetForm(input *dtos.ServiceInput) error {
	if err := newServiceDetails(input); err != nil {
		return err
	}

	return newServiceFrameDetails(input)
}

func newServiceDetails(input *dtos.ServiceInput) error {
	input.ServiceFramework = "sls"

	serviceName := huh.NewInput().
		Title("Service Name:").
		Value(&input.ServiceName).
		Validate(func(s string) error {
			if s == "" {
				return errors.New("service name cannot be empty")
			}

			return nil
		})

	servicePath := huh.NewInput().
		Title("Service Folder:").
		Placeholder(input.ServiceName).
		Description("Enter the folder name where should be placed the service content. If not defined will be the same as service name").
		Value(&input.ServicePath)

	input.FrameVersion = "v2"

	group1 := huh.NewGroup(
		serviceName,
		servicePath,
	).Title("Service Details")

	return huh.NewForm(group1).WithTheme(huh.ThemeCharm()).Run()
}

func newServiceFrameDetails(input *dtos.ServiceInput) error {
	inputs := make([]huh.Field, 0)

	warmupEnabled := huh.NewConfirm().
		Title("Enable warmup?").
		Value(&input.WarmupEnabled)

	customDomain := huh.NewConfirm().
		Title("Configure custom domain?").
		Value(&input.CustomDomain)

	inputs = append(inputs, warmupEnabled, customDomain)

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
