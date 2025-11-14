package newdomain

import (
	"fmt"

	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
)

func GetForm(input *dtos.DomainInput) error {
	input.PackageName = data.Meta.PackageName

	if err := baseForm(input); err != nil {
		return err
	}

	switch input.DBType {
	case data.DBTypePostgres:
		if err := postgresForm(input); err != nil {
			return err
		}
	case data.DBTypeDynamo:
		if err := dynamoForm(input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("database type %s not supported", input.DBType)
	}

	return nil
}
