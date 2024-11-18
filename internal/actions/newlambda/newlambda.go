package newlambda

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Drafteame/draft/internal/actions/dtos"
	"github.com/Drafteame/draft/internal/files"
	"github.com/Drafteame/draft/internal/templates"
)

type NewLambda struct {
	tmpl       templates.SLS
	input      dtos.Input
	lambdaPath string
}

func GetAction() *NewLambda {
	return &NewLambda{}
}

func (nl *NewLambda) Exec(input dtos.Input) error {
	nl.tmpl = templates.NewSLS(input)
	nl.input = input
	nl.lambdaPath = input.ServicePath + "/cmd/" + input.LambdaType + "/" + input.LambdaName

	var err error

	switch input.LambdaType {
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

func (nl *NewLambda) createFiles(files ...dtos.FileEntry) error {
	for _, file := range files {
		if err := os.WriteFile(nl.input.ServicePath+file.Path, file.Data, 0755); err != nil {
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

	data := "- ${file(cmd/%s/%s/lambda-config.yml):function}\n  #:next"
	data = fmt.Sprintf(data, nl.input.LambdaType, nl.input.LambdaName)

	newContent := strings.ReplaceAll(string(content), "#:next", data)

	return files.Create(path, []byte(newContent))
}

func (nl *NewLambda) addToDepsGo() error {
	path := nl.input.ServicePath + "/deps.go"
	content, err := files.Read(path)
	if err != nil {
		return err
	}

	data := "_ \"github.com/Drafteame/api-draftea/services/%s/cmd/%s/%s/handler\"\n\t//:next"
	data = fmt.Sprintf(data, nl.input.ServicePath, nl.input.LambdaType, nl.input.LambdaName)

	newContent := strings.ReplaceAll(string(content), "//:next", data)

	return files.Create(path, []byte(newContent))
}
