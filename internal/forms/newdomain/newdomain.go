package newdomain

import (
	"fmt"
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"regexp"
	"strings"
)

func GetForm(input *dtos.DomainInput) error {
	input.PackageName = data.Meta.PackageName

	if err := baseForm(input); err != nil {
		return err
	}

	switch input.DBType {
	case "postgres":
		if err := postgresForm(input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("database type %s not supported", input.DBType)
	}

	return nil
}

func normalizeDomainName(name string) string {
	name = strings.ToLower(name)

	re := regexp.MustCompile(`[^a-z0-9]`)
	name = re.ReplaceAllString(name, "")

	return name
}
