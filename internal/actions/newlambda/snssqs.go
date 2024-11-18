package newlambda

import (
	"github.com/Drafteame/draft/internal/actions/dtos"
	"github.com/Drafteame/draft/internal/dirs"
)

func (nl *NewLambda) createSnsSqs() error {
	if err := dirs.Create(nl.lambdaPath + "/handler"); err != nil {
		return err
	}

	filesEntries := []dtos.FileEntry{
		{Path: "/cmd/snssqs/" + nl.input.LambdaName + "/main.go", Data: nl.tmpl.FrameV2.SnsSqs.MainGo},
		{Path: "/cmd/snssqs/" + nl.input.LambdaName + "/lambda-config.yml", Data: nl.tmpl.FrameV2.SnsSqs.LambdaConfigYAML},
		{Path: "/cmd/snssqs/" + nl.input.LambdaName + "/handler/handler.go", Data: nl.tmpl.FrameV2.SnsSqs.Handler.HandlerGo},
		{Path: "/cmd/snssqs/" + nl.input.LambdaName + "/handler/worker.go", Data: nl.tmpl.FrameV2.SnsSqs.Handler.WorkerGo},
		{Path: "/cmd/snssqs/" + nl.input.LambdaName + "/handler/bootstrap.go", Data: nl.tmpl.FrameV2.SnsSqs.Handler.BootstrapGo},
	}

	return nl.createFiles(filesEntries...)
}
