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
	typePath := input.LambdaType
	if input.LambdaType == "custom" {
		typePath = input.CustomTypePath
	}

	lambdaPath := input.ServicePath + "/cmd/" + typePath + "/" + input.LambdaName
	if input.IsLegacy {
		lambdaPath = input.ServicePath + "/" + typePath + "/" + input.LambdaName
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

	if err := nl.preCreate(); err != nil {
		return err
	}

	if err := nl.exec(); err != nil {
		return err
	}

	return nl.postCreate()
}

// preCreate performs validation before lambda creation
func (nl *NewLambda) preCreate() error {
	if !files.Exists(nl.input.ServicePath) {
		return fmt.Errorf("service %s not found", nl.input.ServicePath)
	}

	return nil
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
