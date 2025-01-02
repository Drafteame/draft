package deleteproject

import (
	"github.com/Drafteame/draft/internal/dtos"
)

func GetForm(input *dtos.DeleteProjectInput) error {

	if err := baseForm(input); err != nil {
		return err
	}

	return nil
}
