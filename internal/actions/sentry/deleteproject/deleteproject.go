package deleteproject

import (
	"github.com/Drafteame/draft/internal/dtos"
)

type DeleteProject struct {
	input dtos.DeleteProjectInput
}

func GetAction(input dtos.DeleteProjectInput) (*DeleteProject, error) {
	return &DeleteProject{
		input: input,
	}, nil
}

func (dl *DeleteProject) Exec() error {
	if err := dl.sentry(); err != nil {
		return err
	}
	return nil
}
