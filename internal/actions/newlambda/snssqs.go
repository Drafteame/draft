package newlambda

import (
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/dirs"
)

func (nl *NewLambda) createSnsSqs() error {
	if err := dirs.Create(nl.lambdaPath + "/handler"); err != nil {
		return err
	}

	filesEntries := []dtos.FileEntry{
		{Path: "/main.go", Data: nl.tmpl.SnsSqs.MainGo},
		{Path: "/lambda-config.yml", Data: nl.tmpl.SnsSqs.LambdaConfigYAML},
		{Path: "/handler/handler.go", Data: nl.tmpl.SnsSqs.Handler.HandlerGo},
		{Path: "/handler/worker.go", Data: nl.tmpl.SnsSqs.Handler.WorkerGo},
		{Path: "/handler/bootstrap.go", Data: nl.tmpl.SnsSqs.Handler.BootstrapGo},
	}

	return nl.createFiles(filesEntries...)
}
