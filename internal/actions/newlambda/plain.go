package newlambda

import (
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/dirs"
)

func (nl *NewLambda) createPlain() error {
	if err := dirs.Create(nl.lambdaPath + "/handler"); err != nil {
		return err
	}

	filesEntries := []dtos.FileEntry{
		{Path: "/main.go", Data: nl.tmpl.FrameV2.Plain.MainGo},
		{Path: "/lambda-config.yml", Data: nl.tmpl.FrameV2.Plain.LambdaConfigYAML},
		{Path: "/handler/handler.go", Data: nl.tmpl.FrameV2.Plain.Handler.HandlerGo},
		{Path: "/handler/bootstrap.go", Data: nl.tmpl.FrameV2.Plain.Handler.BootstrapGo},
	}

	return nl.createFiles(filesEntries...)
}
