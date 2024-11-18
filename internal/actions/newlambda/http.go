package newlambda

import (
	"github.com/Drafteame/draft/internal/actions/dtos"
	"github.com/Drafteame/draft/internal/dirs"
)

func (nl *NewLambda) createHttp() error {
	if err := dirs.Create(nl.lambdaPath + "/handler"); err != nil {
		return err
	}

	filesEntries := []dtos.FileEntry{
		{Path: "/cmd/http/" + nl.input.LambdaName + "/main.go", Data: nl.tmpl.FrameV2.HTTP.MainGo},
		{Path: "/cmd/http/" + nl.input.LambdaName + "/lambda-config.yml", Data: nl.tmpl.FrameV2.HTTP.LambdaConfigYAML},
		{Path: "/cmd/http/" + nl.input.LambdaName + "/handler/handler.go", Data: nl.tmpl.FrameV2.HTTP.Handler.HandlerGo},
		{Path: "/cmd/http/" + nl.input.LambdaName + "/handler/bootstrap.go", Data: nl.tmpl.FrameV2.HTTP.Handler.BootstrapGo},
	}

	return nl.createFiles(filesEntries...)
}
