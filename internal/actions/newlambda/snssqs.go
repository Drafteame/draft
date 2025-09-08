package newlambda

import (
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/dirs"
)

func (nl *NewLambda) createSnsSqs() error {
	if !nl.input.WithFrame {
		return nl.createNativeSnsSqs()
	}
	return nl.createSnsSqsWithFrame()
}

func (nl *NewLambda) createSnsSqsWithFrame() error {
	if err := dirs.Create(nl.lambdaPath + "/handler"); err != nil {
		return err
	}

	filesEntries := []dtos.FileEntry{
		{Path: "/main.go", Data: nl.tmpl.SnsSqs.MainGo},
		{Path: "/lambda-config.yml", Data: nl.tmpl.SnsSqs.LambdaConfigYAML},
		{Path: "/handler/handler.go", Data: nl.tmpl.SnsSqs.Handler.HandlerGo},
		{Path: "/handler/worker.go", Data: nl.tmpl.SnsSqs.Handler.WorkerGo},
		{Path: "/handler/bootstrap.go", Data: nl.tmpl.SnsSqs.Handler.BootstrapGo},
		{Path: "/handler/provider.go", Data: nl.tmpl.SnsSqs.Handler.ProviderGo},
	}

	return nl.createFiles(filesEntries...)
}

func (nl *NewLambda) createNativeSnsSqs() error {
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
		{Path: "/main.go", Data: nl.tmpl.SnsSqs.MainGo},
		{Path: "/lambda-config.yml", Data: nl.tmpl.SnsSqs.LambdaConfigYAML},
		{Path: "/handler/bootstrap.go", Data: nl.tmpl.SnsSqs.Handler.BootstrapGo},
		{Path: "/handler/worker/worker.go", Data: nl.tmpl.SnsSqs.Handler.WorkerGo},
		{Path: "/handler/worker/resources.go", Data: nl.tmpl.SnsSqs.Handler.ResourcesGo},
		{Path: "/handler/embed/_.yaml", Data: nl.tmpl.SnsSqs.Handler.EmbedYML},
	}

	return nl.createFiles(filesEntries...)
}
