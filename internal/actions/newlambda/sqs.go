package newlambda

import (
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/dirs"
)

func (nl *NewLambda) createSqs() error {
	if err := dirs.Create(nl.lambdaPath + "/handler"); err != nil {
		return err
	}

	filesEntries := []dtos.FileEntry{
		{Path: "/main.go", Data: nl.tmpl.FrameV2.Sqs.MainGo},
		{Path: "/lambda-config.yml", Data: nl.tmpl.FrameV2.Sqs.LambdaConfigYAML},
		{Path: "/handler/handler.go", Data: nl.tmpl.FrameV2.Sqs.Handler.HandlerGo},
		{Path: "/handler/worker.go", Data: nl.tmpl.FrameV2.Sqs.Handler.WorkerGo},
		{Path: "/handler/bootstrap.go", Data: nl.tmpl.FrameV2.Sqs.Handler.BootstrapGo},
	}

	return nl.createFiles(filesEntries...)
}
