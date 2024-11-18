package newlambda

import (
	"github.com/Drafteame/draft/internal/actions/dtos"
	"github.com/Drafteame/draft/internal/dirs"
)

func (nl *NewLambda) createPlain() error {
	if err := dirs.Create(nl.lambdaPath + "/handler"); err != nil {
		return err
	}

	filesEntries := []dtos.FileEntry{
		{Path: "/cmd/plain/" + nl.input.LambdaName + "/main.go", Data: nl.tmpl.FrameV2.Plain.MainGo},
		{Path: "/cmd/plain/" + nl.input.LambdaName + "/lambda-config.yml", Data: nl.tmpl.FrameV2.Plain.LambdaConfigYAML},
		{Path: "/cmd/plain/" + nl.input.LambdaName + "/handler/handler.go", Data: nl.tmpl.FrameV2.Plain.Handler.HandlerGo},
		{Path: "/cmd/plain/" + nl.input.LambdaName + "/handler/bootstrap.go", Data: nl.tmpl.FrameV2.Plain.Handler.BootstrapGo},
	}

	return nl.createFiles(filesEntries...)
}
