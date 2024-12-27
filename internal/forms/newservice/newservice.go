package newservice

import (
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
)

func GetForm(input *dtos.ServiceInput) error {
	input.ServiceFramework = "sls"
	input.FrameVersion = "v2"
	input.LambdaName = "helloworld"
	input.LambdaType = "plain"
	input.PackageName = data.Meta.PackageName
	input.NextImportTag = data.NextImportTag
	input.NextLambdaImportTag = data.NextLambdaImportTag

	if err := baseForm(input); err != nil {
		return err
	}

	return frameDetails(input)
}
