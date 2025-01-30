package forms

import (
	"github.com/Masterminds/semver"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/forms/newdomain"
	"github.com/Drafteame/draft/internal/forms/newlambda"
	"github.com/Drafteame/draft/internal/forms/newservice"
	"github.com/Drafteame/draft/internal/forms/nixversion"
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

func DeleteProject(input *dtos.DeleteProjectInput) error {
	return deleteproject.GetForm(input)
}

func UpdateNixModules(input *dtos.UpdateNixModules, currentVersion, latestVersion *semver.Version) error {
	return nixversion.GetForm(input, currentVersion, latestVersion)
}
