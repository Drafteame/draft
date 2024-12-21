package newdomain

import (
	"errors"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/inputs"
)

func baseForm(input *dtos.DomainInput) error {
	err := inputs.Text("Domain Path:",
		inputs.WithValue(&input.DomainPath),
		inputs.WithDescription[string]("Enter the path to the domain excluding 'domains' folder."),
		inputs.WithValidation(func(val string) error {
			if val == "" {
				return errors.New("domain path cannot be empty")
			}

			return nil
		}),
	)

	if err != nil {
		return err
	}

	input.DomainPath = strings.ToLower(input.DomainPath)

	domainName, _ := lo.Last(strings.Split(input.DomainPath, "/"))

	input.DomainName = normalizeDomainName(domainName)
	input.DomainNameLower = strings.ToLower(input.DomainName)
	input.DomainNamePascal = cases.Title(language.English).String(input.DomainName)

	err = inputs.Select[string]("Select DB Type:",
		inputs.WithDescription[string]("Select the type of database you want to use"),
		inputs.WithValue(&input.DBType),
		inputs.WithOptions(map[string]string{
			"Postgres": "postgres",
		}),
	)

	return err
}
