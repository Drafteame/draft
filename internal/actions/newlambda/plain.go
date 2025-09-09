package newlambda

import (
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/dirs"
)

func (nl *NewLambda) createPlain() error {
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
		{Path: "/main.go", Data: nl.tmpl.Plain.MainGo},
		{Path: "/lambda-config.yml", Data: nl.tmpl.Plain.LambdaConfigYAML},
		{Path: "/handler/bootstrap.go", Data: nl.tmpl.Plain.Handler.BootstrapGo},
		{Path: "/handler/worker/worker.go", Data: nl.tmpl.Plain.Handler.WorkerGo},
		{Path: "/handler/worker/resources.go", Data: nl.tmpl.Plain.Handler.ResourcesGo},
		{Path: "/handler/embed/_.yaml", Data: nl.tmpl.Plain.Handler.EmbedYML},
		{Path: "/handler/dtos/dto.go", Data: nl.tmpl.Plain.Handler.DtosGo},
	}

	return nl.createFiles(filesEntries...)
}
