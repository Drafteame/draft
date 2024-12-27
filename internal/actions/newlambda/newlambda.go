package newlambda

import (
	"fmt"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/templates"
)

type NewLambda struct {
	tmpl       templates.SLS
	input      dtos.ServiceInput
	lambdaPath string
}

func GetAction(input dtos.ServiceInput) *NewLambda {
	return &NewLambda{
		input:      input,
		tmpl:       templates.NewSLS(input),
		lambdaPath: input.ServicePath + "/cmd/" + input.LambdaType + "/" + input.LambdaName,
	}
}

func (nl *NewLambda) Exec() error {
	if !files.Exists(nl.input.ServicePath) {
		return fmt.Errorf("service %s not found", nl.input.ServicePath)
	}

	if err := nl.exec(); err != nil {
		return err
	}

	return nl.postCreate()
}

func (nl *NewLambda) createFiles(entries ...dtos.FileEntry) error {
	for _, file := range entries {
		path := nl.lambdaPath + file.Path

		if err := files.Create(path, file.Data); err != nil {
			return err
		}
	}

	return nil
}
