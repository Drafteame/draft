package newlambda

import (
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/dirs"
)

func (nl *NewLambda) createCustom() error {
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
		{Path: "/main.go", Data: nl.tmpl.Custom.MainGo},
		{Path: "/lambda-config.yml", Data: nl.tmpl.Custom.LambdaConfigYAML},
		{Path: "/handler/bootstrap.go", Data: nl.tmpl.Custom.Handler.BootstrapGo},
		{Path: "/handler/worker/worker.go", Data: nl.tmpl.Custom.Handler.WorkerGo},
		{Path: "/handler/worker/resources.go", Data: nl.tmpl.Custom.Handler.ResourcesGo},
		{Path: "/handler/worker/worker_setup_test.go", Data: nl.tmpl.Custom.Handler.WorkerSetupTestGo},
		{Path: "/handler/worker/worker_test.go", Data: nl.tmpl.Custom.Handler.WorkerTestGo},
		{Path: "/handler/embed/_.yaml", Data: nl.tmpl.Custom.Handler.EmbedYML},
	}

	// Add idempotency files if requested
	if nl.input.UseIdempotency {
		filesEntries = append(filesEntries, dtos.FileEntry{
			Path: "/handler/worker/idempotency.go",
			Data: nl.tmpl.Custom.Handler.IdempotencyGo,
		})
		filesEntries = append(filesEntries, dtos.FileEntry{
			Path: "/handler/worker/interfaces.go",
			Data: nl.tmpl.Custom.Handler.InterfacesGo,
		})
	}

	return nl.createFiles(filesEntries...)
}
