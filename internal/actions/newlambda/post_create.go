package newlambda

import (
	"fmt"
	"strings"

	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/pkg/files"
)

func (nl *NewLambda) postCreate() error {
	if nl.input.IsLegacy {
		println("Command executed in legacy mode. No deps and serverless.yml changes created. Please add manually")

		return nil
	}

	actions := []func() error{
		nl.addToDepsGo,
		nl.addToServerlessYAML,
		nl.format,
	}

	for _, action := range actions {
		if err := action(); err != nil {
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

	line := "- ${file(cmd/%s/%s/lambda-config.yml):function}\n  " + data.NextLambdaImportTag
	line = fmt.Sprintf(line, nl.input.LambdaType, nl.input.LambdaName)

	newContent := strings.ReplaceAll(string(content), data.NextLambdaImportTag, line)

	return files.Create(path, []byte(newContent))
}

func (nl *NewLambda) addToDepsGo() error {
	path := nl.input.ServicePath + "/deps.go"
	content, err := files.Read(path)
	if err != nil {
		return err
	}

	line := "_ \"%s/%s/cmd/%s/%s/handler\"\n\t" + data.NextImportTag
	line = fmt.Sprintf(line, nl.input.PackageName, nl.input.ServicePath, nl.input.LambdaType, nl.input.LambdaName)

	newContent := strings.ReplaceAll(string(content), data.NextImportTag, line)

	return files.Create(path, []byte(newContent))
}
