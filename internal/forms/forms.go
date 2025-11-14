package forms

import (
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/forms/newdomain"
	"github.com/Drafteame/draft/internal/forms/newlambda"
	"github.com/Drafteame/draft/internal/forms/newservice"
	"github.com/Drafteame/draft/internal/forms/sentry/createproject"
	"github.com/Drafteame/draft/internal/forms/sentry/deleteproject"
)

func NewService(input *dtos.ServiceInput) error {
	return newservice.GetForm(input)
}

func NewLambda(input *dtos.LambdaInput) error {
	return newlambda.GetForm(input)
}

func NewDomain(input *dtos.DomainInput) error {
	return newdomain.GetForm(input)
}

func CreateProject(input *dtos.CreateProjectInput) error {
	return createproject.GetForm(input)
}

func DeleteProject(input *dtos.DeleteProjectInput) error {
	return deleteproject.GetForm(input)
}
