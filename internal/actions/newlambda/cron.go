package newlambda

import (
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/dirs"
)

func (nl *NewLambda) createCron() error {
	if err := dirs.Create(nl.lambdaPath + "/handler"); err != nil {
		return err
	}

	if err := dirs.Create(nl.lambdaPath + "/handler/worker"); err != nil {
		return err
	}

	if err := dirs.Create(nl.lambdaPath + "/handler/embed"); err != nil {
		return err
	}

	filesEntries := []dtos.FileEntry{
		{Path: "/main.go", Data: nl.tmpl.Cron.MainGo},
		{Path: "/lambda-config.yml", Data: nl.tmpl.Cron.LambdaConfigYAML},
		{Path: "/handler/bootstrap.go", Data: nl.tmpl.Cron.Handler.BootstrapGo},
		{Path: "/handler/worker/worker.go", Data: nl.tmpl.Cron.Handler.WorkerGo},
		{Path: "/handler/worker/resources.go", Data: nl.tmpl.Cron.Handler.ResourcesGo},
		{Path: "/handler/embed/_.yaml", Data: nl.tmpl.Cron.Handler.EmbedYML},
	}

	return nl.createFiles(filesEntries...)
}
