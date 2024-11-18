package forms

import (
	"github.com/Drafteame/draft/internal/actions/dtos"
	"github.com/Drafteame/draft/internal/forms/newlambda"
	"github.com/Drafteame/draft/internal/forms/newservice"
)

func NewService(input *dtos.Input) error {
	return newservice.GetForm(input)
}

func NewLambda(input *dtos.Input) error {
	return newlambda.GetForm(input)
}
