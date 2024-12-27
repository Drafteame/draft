package newlambda

import (
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/dirs"
)

func (nl *NewLambda) createPlain() error {
	if err := dirs.Create(nl.lambdaPath + "/handler"); err != nil {
		return err
	}

	filesEntries := []dtos.FileEntry{
		{Path: "/main.go", Data: nl.tmpl.Plain.MainGo},
		{Path: "/lambda-config.yml", Data: nl.tmpl.Plain.LambdaConfigYAML},
		{Path: "/handler/handler.go", Data: nl.tmpl.Plain.Handler.HandlerGo},
		{Path: "/handler/bootstrap.go", Data: nl.tmpl.Plain.Handler.BootstrapGo},
	}

	return nl.createFiles(filesEntries...)
}
