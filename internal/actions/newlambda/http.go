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
		{Path: "/main.go", Data: nl.tmpl.FrameV2.HTTP.MainGo},
		{Path: "/lambda-config.yml", Data: nl.tmpl.FrameV2.HTTP.LambdaConfigYAML},
		{Path: "/handler/handler.go", Data: nl.tmpl.FrameV2.HTTP.Handler.HandlerGo},
		{Path: "/handler/bootstrap.go", Data: nl.tmpl.FrameV2.HTTP.Handler.BootstrapGo},
	}

	return nl.createFiles(filesEntries...)
}
