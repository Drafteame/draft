package newlambda

import (
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/dirs"
)

func (nl *NewLambda) createSnsSqs() error {
	if err := dirs.Create(nl.lambdaPath + "/handler"); err != nil {
		return err
	}

	if err := dirs.Create(nl.lambdaPath + "/handler/worker"); err != nil {
		return err
	}

	if err := dirs.Create(nl.lambdaPath + "/handler/embed"); err != nil {
		return err
	}

	if err := dirs.Create(nl.lambdaPath + "/handler/dtos"); err != nil {
		return err
	}

	filesEntries := []dtos.FileEntry{
		{Path: "/main.go", Data: nl.tmpl.SnsSqs.MainGo},
		{Path: "/lambda-config.yml", Data: nl.tmpl.SnsSqs.LambdaConfigYAML},
		{Path: "/handler/bootstrap.go", Data: nl.tmpl.SnsSqs.Handler.BootstrapGo},
		{Path: "/handler/worker/worker.go", Data: nl.tmpl.SnsSqs.Handler.WorkerGo},
		{Path: "/handler/worker/resources.go", Data: nl.tmpl.SnsSqs.Handler.ResourcesGo},
		{Path: "/handler/embed/_.yaml", Data: nl.tmpl.SnsSqs.Handler.EmbedYML},
		{Path: "/handler/dtos/dto.go", Data: nl.tmpl.SnsSqs.Handler.DtosGo},
	}

	return nl.createFiles(filesEntries...)
}
