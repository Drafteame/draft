package deleteproject

import (
	"github.com/Drafteame/draft/internal/dtos"
)

type DeleteProject struct {
	input dtos.DeleteProjectInput
}

func New(input dtos.DeleteProjectInput) *DeleteProject {
	return &DeleteProject{
		input: input,
	}
}

func (dl *DeleteProject) Exec() error {
	return dl.sentry()
}
