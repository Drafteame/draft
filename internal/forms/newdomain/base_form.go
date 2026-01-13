package newdomain

import (
	"errors"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/inputs"
)

func baseForm(input *dtos.DomainInput) error {
	// Prompt for domain path only if not provided via flag
	if input.DomainPath == "" {
		err := inputs.Text("Domain Path:",
			inputs.WithValue(&input.DomainPath),
			inputs.WithDescription[string]("Enter the path to the domain folder."),
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
	}

	// Validate domain path if provided via flag
	if input.DomainPath == "" {
		return errors.New("domain path cannot be empty")
	}

	input.DomainPath = strings.ToLower(input.DomainPath)

	domainName, _ := lo.Last(strings.Split(input.DomainPath, "/"))

	input.DomainName = normalizeDomainName(domainName)
	input.DomainNameLower = strings.ToLower(input.DomainName)
	input.DomainNamePascal = cases.Title(language.English).String(input.DomainName)

	if !strings.HasPrefix(input.DomainPath, "domains/") {
		input.DomainPath = "domains/" + input.DomainPath
	}

	// Prompt for DB type only if not provided via flag
	if input.DBType == "" {
		err := inputs.Select[string]("Select DB Type:",
			inputs.WithDescription[string]("Select the type of database you want to use"),
			inputs.WithValue(&input.DBType),
			inputs.WithOptions(map[string]string{
				"Postgres": data.DBTypePostgres,
				"DynamoDB": data.DBTypeDynamo,
			}),
		)

		if err != nil {
			return err
		}
	}

	// Validate DB type if provided via flag
	if input.DBType != data.DBTypePostgres && input.DBType != data.DBTypeDynamo {
		return errors.New("invalid db-type: must be 'postgres' or 'dynamo'")
	}

	return nil
}

func normalizeDomainName(name string) string {
	name = strings.ToLower(name)

	re := regexp.MustCompile(`[^a-z0-9]`)
	name = re.ReplaceAllString(name, "")

	return name
}
