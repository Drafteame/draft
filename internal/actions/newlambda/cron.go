package newlambda

import (
	"github.com/Drafteame/draft/internal/actions/dtos"
	"github.com/Drafteame/draft/internal/dirs"
)

func (nl *NewLambda) createCron() error {
	if err := dirs.Create(nl.lambdaPath + "/handler"); err != nil {
		return err
	}

	filesEntries := []dtos.FileEntry{
		{Path: "/cmd/cron/" + nl.input.LambdaName + "/main.go", Data: nl.tmpl.FrameV2.Cron.MainGo},
		{Path: "/cmd/cron/" + nl.input.LambdaName + "/lambda-config.yml", Data: nl.tmpl.FrameV2.Cron.LambdaConfigYAML},
		{Path: "/cmd/cron/" + nl.input.LambdaName + "/handler/handler.go", Data: nl.tmpl.FrameV2.Cron.Handler.HandlerGo},
		{Path: "/cmd/cron/" + nl.input.LambdaName + "/handler/bootstrap.go", Data: nl.tmpl.FrameV2.Cron.Handler.BootstrapGo},
	}

	return nl.createFiles(filesEntries...)
}
