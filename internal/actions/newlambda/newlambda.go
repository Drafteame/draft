package newlambda

import (
	"fmt"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/templates"
)

type NewLambda struct {
	tmpl       *templates.LambdaTemplates
	input      dtos.LambdaInput
	lambdaPath string
}

func New(input dtos.LambdaInput) *NewLambda {
	lambdaPath := input.ServicePath + "/cmd/" + input.LambdaType + "/" + input.LambdaName
	if input.IsLegacy {
		lambdaPath = input.ServicePath + "/" + input.LambdaType + "/" + input.LambdaName
	}

	return &NewLambda{
		input:      input,
		lambdaPath: lambdaPath,
	}
}

func (nl *NewLambda) Exec() error {
	tmpl, err := templates.NewLambdaTemplates(nl.input)
	if err != nil {
		return err
	}

	nl.tmpl = tmpl

	if !files.Exists(nl.input.ServicePath) {
		return fmt.Errorf("service %s not found", nl.input.ServicePath)
	}

	if errExec := nl.exec(); errExec != nil {
		return errExec
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
