package newlambda

import (
	"github.com/Drafteame/draft/internal/actions/dtos"
	"github.com/Drafteame/draft/internal/dirs"
)

func (nl *NewLambda) createSqs() error {
	if err := dirs.Create(nl.lambdaPath + "/handler"); err != nil {
		return err
	}

	filesEntries := []dtos.FileEntry{
		{Path: "/cmd/sqs/" + nl.input.LambdaName + "/main.go", Data: nl.tmpl.FrameV2.Sqs.MainGo},
		{Path: "/cmd/sqs/" + nl.input.LambdaName + "/lambda-config.yml", Data: nl.tmpl.FrameV2.Sqs.LambdaConfigYAML},
		{Path: "/cmd/sqs/" + nl.input.LambdaName + "/handler/handler.go", Data: nl.tmpl.FrameV2.Sqs.Handler.HandlerGo},
		{Path: "/cmd/sqs/" + nl.input.LambdaName + "/handler/worker.go", Data: nl.tmpl.FrameV2.Sqs.Handler.WorkerGo},
		{Path: "/cmd/sqs/" + nl.input.LambdaName + "/handler/bootstrap.go", Data: nl.tmpl.FrameV2.Sqs.Handler.BootstrapGo},
	}

	return nl.createFiles(filesEntries...)
}
