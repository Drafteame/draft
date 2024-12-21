package newlambda

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Drafteame/draft/internal/data"
	dtos2 "github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/templates"
)

type NewLambda struct {
	tmpl       templates.SLS
	input      dtos2.ServiceInput
	lambdaPath string
}

func GetAction(input dtos2.ServiceInput) *NewLambda {
	input.PackageName = data.Meta.PackageName
	input.ServicePath = "services/" + input.ServicePath

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

	var err error

	switch nl.input.LambdaType {
	case "plain":
		err = nl.createPlain()
	case "sqs":
		err = nl.createSqs()
	case "http":
		err = nl.createHttp()
	case "snssqs":
		err = nl.createSnsSqs()
	case "cron":
		err = nl.createCron()
	default:
		err = errors.New("unsupported lambda type")
	}

	if err != nil {
		return err
	}

	if err := nl.addToServerlessYAML(); err != nil {
		return err
	}

	return nl.addToDepsGo()
}

func (nl *NewLambda) createFiles(files ...dtos2.FileEntry) error {
	for _, file := range files {
		path := nl.lambdaPath + file.Path

		if err := os.WriteFile(path, file.Data, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (nl *NewLambda) addToServerlessYAML() error {
	path := nl.input.ServicePath + "/serverless.yml"
	content, err := files.Read(path)
	if err != nil {
		return err
	}

	line := "- ${file(cmd/%s/%s/lambda-config.yml):function}\n  #:next"
	line = fmt.Sprintf(line, nl.input.LambdaType, nl.input.LambdaName)

	newContent := strings.ReplaceAll(string(content), "#:next", line)

	return files.Create(path, []byte(newContent))
}

func (nl *NewLambda) addToDepsGo() error {
	path := nl.input.ServicePath + "/deps.go"
	content, err := files.Read(path)
	if err != nil {
		return err
	}

	line := "_ \"%s/%s/cmd/%s/%s/handler\"\n\t//:next"
	line = fmt.Sprintf(line, nl.input.PackageName, nl.input.ServicePath, nl.input.LambdaType, nl.input.LambdaName)

	newContent := strings.ReplaceAll(string(content), "//:next", line)

	return files.Create(path, []byte(newContent))
}
