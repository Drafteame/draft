package forms

import (
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/forms/newdomain"
	"github.com/Drafteame/draft/internal/forms/newlambda"
	"github.com/Drafteame/draft/internal/forms/newservice"
)

func NewService(input *dtos.ServiceInput) error {
	return newservice.GetForm(input)
}

func NewLambda(input *dtos.ServiceInput) error {
	return newlambda.GetForm(input)
}

func NewDomain(input *dtos.DomainInput) error {
	return newdomain.GetForm(input)
}
