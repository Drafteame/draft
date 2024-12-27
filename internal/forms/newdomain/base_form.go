package newdomain

import (
	"errors"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/Drafteame/draft/internal/dtos"
	inputs2 "github.com/Drafteame/draft/internal/pkg/inputs"
)

func baseForm(input *dtos.DomainInput) error {
	err := inputs2.Text("Domain Path:",
		inputs2.WithValue(&input.DomainPath),
		inputs2.WithDescription[string]("Enter the path to the domain excluding 'domains' folder."),
		inputs2.WithValidation(func(val string) error {
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

	err = inputs2.Select[string]("Select DB Type:",
		inputs2.WithDescription[string]("Select the type of database you want to use"),
		inputs2.WithValue(&input.DBType),
		inputs2.WithOptions(map[string]string{
			"Postgres": "postgres",
		}),
	)

	return err
}
