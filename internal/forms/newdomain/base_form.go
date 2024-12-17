package newdomain

import (
	"errors"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/charmbracelet/huh"
	"github.com/samber/lo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func baseForm(input *dtos.DomainInput) error {
	err := huh.NewInput().
		Title("Domain Path:").
		Description("Enter the path to the domain excluding 'domains' folder.").
		Value(&input.DomainPath).
		Validate(func(s string) error {
			if s == "" {
				return errors.New("domain path cannot be empty")
			}

			return nil
		}).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		return err
	}

	input.DomainPath = strings.ToLower(input.DomainPath)

	domainName, _ := lo.Last(strings.Split(input.DomainPath, "/"))

	input.DomainName = normalizeDomainName(domainName)
	input.DomainNameLower = strings.ToLower(input.DomainName)
	input.DomainNamePascal = cases.Title(language.English).String(input.DomainName)

	err = huh.NewSelect[string]().
		Title("Select DB Type:").
		Description("Select the type of database you want to use").
		Value(&input.DBType).
		Options(
			huh.NewOption("Postgres", "postgres"),
			// huh.NewOption("Mongo", "mongo"),
		).WithTheme(huh.ThemeCharm()).Run()

	return err
}
