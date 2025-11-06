package createproject

import (
	"github.com/Drafteame/draft/internal/dtos"
)

type CreateProject struct {
	input dtos.CreateProjectInput
}

func New(input dtos.CreateProjectInput) *CreateProject {
	return &CreateProject{
		input: input,
	}
}

func (cp *CreateProject) Exec() error {
	return cp.sentry()
}
