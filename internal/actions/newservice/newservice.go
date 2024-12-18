package newservice

import (
	"errors"
	"os"

	"github.com/Drafteame/draft/internal/data"
	dtos2 "github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/dirs"
	"github.com/Drafteame/draft/internal/project"
	"github.com/Drafteame/draft/internal/templates"
)

type NewService struct {
	tmpl  templates.SLS
	input dtos2.ServiceInput
}

func GetAction(input dtos2.ServiceInput) *NewService {
	input.LambdaName = "helloworld"
	input.LambdaType = "plain"
	input.PackageName = data.Meta.PackageName
	input.NormalizedServiceName = project.NormalizeServiceName(input.ServiceName)

	if input.ServicePath == "" {
		input.ServicePath = input.ServiceName
	}

	input.ServicePath = "services/" + input.ServicePath

	return &NewService{input: input}
}

func (css *NewService) Exec() error {
	if err := css.preCreate(); err != nil {
		return err
	}

	if err := css.exec(); err != nil {
		return err
	}

	return css.postCreate()
}

func (css *NewService) exec() error {
	css.tmpl = templates.NewSLS(css.input)

	switch css.input.ServiceFramework {
	case "sls":
		return css.createServerless()
	default:
		return errors.New("unsupported service framework")
	}
}

func (css *NewService) createServerless() error {
	if err := css.createAllDirs(); err != nil {
		return err
	}

	if err := css.createLambdaFolders(); err != nil {
		return err
	}

	return css.createFiles()
}

func (css *NewService) createAllDirs() error {
	folders := []string{
		css.input.ServicePath + "/cmd/plain",
		css.input.ServicePath + "/config/app",
		css.input.ServicePath + "/config/sls",
	}

	return dirs.Create(folders...)
}

func (css *NewService) createLambdaFolders() error {
	folders := []string{
		css.input.ServicePath + "/cmd/" + css.input.LambdaType + "/" + css.input.LambdaName + "/handler",
	}

	return dirs.Create(folders...)
}

func (css *NewService) createFiles() error {
	files := css.getFileList()

	for _, file := range files {
		if err := os.WriteFile(css.input.ServicePath+file.Path, file.Data, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (css *NewService) getFileList() []dtos2.FileEntry {
	entries := []dtos2.FileEntry{
		{Path: "/serverless.yml", Data: css.tmpl.ServerlessYAML},
		{Path: "/package.json", Data: css.tmpl.PackageJSON},
		{Path: "/deps.go", Data: css.tmpl.DepsGo},
		{Path: "/config/app/app.pkl", Data: css.tmpl.Config.App.AppPkl},
		{Path: "/config/app/modules.pkl", Data: css.tmpl.Config.App.ModulesPkl},
		{Path: "/config/sls/environment.yml", Data: css.tmpl.Config.Sls.EnvironmentYAML},
		{Path: "/config/sls/iam.yml", Data: css.tmpl.Config.Sls.IamYAML},
	}

	if css.input.FrameVersion == "v2" {
		entries = append(entries, css.getFrameV2Entries()...)
	}

	return entries
}

func (css *NewService) getFrameV2Entries() []dtos2.FileEntry {
	return []dtos2.FileEntry{
		{Path: "/cmd/plain/" + css.input.LambdaName + "/main.go", Data: css.tmpl.FrameV2.Plain.MainGo},
		{Path: "/cmd/plain/" + css.input.LambdaName + "/lambda-config.yml", Data: css.tmpl.FrameV2.Plain.LambdaConfigYAML},
		{Path: "/cmd/plain/" + css.input.LambdaName + "/handler/handler.go", Data: css.tmpl.FrameV2.Plain.Handler.HandlerGo},
		{Path: "/cmd/plain/" + css.input.LambdaName + "/handler/bootstrap.go", Data: css.tmpl.FrameV2.Plain.Handler.BootstrapGo},
	}
}
