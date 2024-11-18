package newlambda

import (
	"github.com/charmbracelet/huh"

	"github.com/Drafteame/draft/internal/actions/dtos"
)

func httpForm(input *dtos.Input) error {
	httpMethod := huh.NewSelect[string]().
		Title("Select HTTP Method:").
		Options(
			huh.NewOption("GET", "GET"),
			huh.NewOption("POST", "POST"),
			huh.NewOption("PUT", "PUT"),
			huh.NewOption("DELETE", "DELETE"),
		).Value(&input.HTTPMethod)

	httpPath := huh.NewInput().
		Title("Set HTTP Path:").
		Description("Enter the path to the service").
		Value(&input.HTTPPath)

	group := huh.NewGroup(
		httpMethod,
		httpPath,
	).Title("HTTP Details")

	return huh.NewForm(group).WithTheme(huh.ThemeCharm()).Run()
}
