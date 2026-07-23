package newservice

import (
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/project"
)

func GetForm(input *dtos.ServiceInput) error {
	input.ServiceFramework = "sls"
	input.LambdaName = "helloworld"
	input.LambdaType = "plain"
	input.ReservedConcurrency = "medium.http"
	input.PackageName = data.Meta.PackageName
	input.NextImportTag = data.NextImportTag
	input.NextLambdaImportTag = data.NextLambdaImportTag

	if err := baseForm(input); err != nil {
		return err
	}

	input.RoleName = project.CapitalizeServiceName(input.NormalizedServiceName)

	return frameDetails(input)
}
