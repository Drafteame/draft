package newlambda

import (
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/dirs"
)

func (nl *NewLambda) createSqs() error {
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
		{Path: "/main.go", Data: nl.tmpl.Sqs.MainGo},
		{Path: "/lambda-config.yml", Data: nl.tmpl.Sqs.LambdaConfigYAML},
		{Path: "/handler/bootstrap.go", Data: nl.tmpl.Sqs.Handler.BootstrapGo},
		{Path: "/handler/worker/worker.go", Data: nl.tmpl.Sqs.Handler.WorkerGo},
		{Path: "/handler/worker/resources.go", Data: nl.tmpl.Sqs.Handler.ResourcesGo},
		{Path: "/handler/worker/idempotency.go", Data: nl.tmpl.Sqs.Handler.IdempotencyGo},
		{Path: "/handler/worker/interfaces.go", Data: nl.tmpl.Sqs.Handler.InterfacesGo},
		{Path: "/handler/embed/_.yaml", Data: nl.tmpl.Sqs.Handler.EmbedYML},
		{Path: "/handler/dtos/dto.go", Data: nl.tmpl.Sqs.Handler.DtosGo},
	}

	return nl.createFiles(filesEntries...)
}
