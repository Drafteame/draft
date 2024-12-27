package newservice

import (
	"errors"
	"os"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/dirs"
	"github.com/Drafteame/draft/internal/templates"
)

type NewService struct {
	tmpl  templates.SLS
	input dtos.ServiceInput
}

func GetAction(input dtos.ServiceInput) *NewService {
	return &NewService{input: input}
}

func (ns *NewService) Exec() error {
	if err := ns.preCreate(); err != nil {
		return err
	}

	if err := ns.exec(); err != nil {
		return err
	}

	return ns.postCreate()
}

func (ns *NewService) exec() error {
	ns.tmpl = templates.NewSLS(ns.input)

	switch ns.input.ServiceFramework {
	case "sls":
		return ns.createServerless()
	default:
		return errors.New("unsupported service framework")
	}
}

func (ns *NewService) createServerless() error {
	if err := ns.createAllDirs(); err != nil {
		return err
	}

	if err := ns.createLambdaFolders(); err != nil {
		return err
	}

	return ns.createFiles()
}

func (ns *NewService) createAllDirs() error {
	folders := []string{
		ns.input.ServicePath + "/cmd/plain",
		ns.input.ServicePath + "/config/app",
		ns.input.ServicePath + "/config/sls",
	}

	return dirs.Create(folders...)
}

func (ns *NewService) createLambdaFolders() error {
	folders := []string{
		ns.input.ServicePath + "/cmd/" + ns.input.LambdaType + "/" + ns.input.LambdaName + "/handler",
	}

	return dirs.Create(folders...)
}

func (ns *NewService) createFiles() error {
	files := ns.getFileList()

	for _, file := range files {
		if err := os.WriteFile(ns.input.ServicePath+file.Path, file.Data, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (ns *NewService) getFileList() []dtos.FileEntry {
	entries := []dtos.FileEntry{
		{Path: "/serverless.yml", Data: ns.tmpl.ServerlessYAML},
		{Path: "/package.json", Data: ns.tmpl.PackageJSON},
		{Path: "/deps.go", Data: ns.tmpl.DepsGo},
		{Path: "/config/app/app.pkl", Data: ns.tmpl.Config.App.AppPkl},
		{Path: "/config/app/modules.pkl", Data: ns.tmpl.Config.App.ModulesPkl},
		{Path: "/config/sls/environment.yml", Data: ns.tmpl.Config.Sls.EnvironmentYAML},
		{Path: "/config/sls/iam.yml", Data: ns.tmpl.Config.Sls.IamYAML},
	}

	if ns.input.FrameVersion == "v2" {
		entries = append(entries, ns.getFrameV2Entries()...)
	}

	return entries
}

func (ns *NewService) getFrameV2Entries() []dtos.FileEntry {
	return []dtos.FileEntry{
		{Path: "/cmd/plain/" + ns.input.LambdaName + "/main.go", Data: ns.tmpl.FrameV2.Plain.MainGo},
		{Path: "/cmd/plain/" + ns.input.LambdaName + "/lambda-config.yml", Data: ns.tmpl.FrameV2.Plain.LambdaConfigYAML},
		{Path: "/cmd/plain/" + ns.input.LambdaName + "/handler/handler.go", Data: ns.tmpl.FrameV2.Plain.Handler.HandlerGo},
		{Path: "/cmd/plain/" + ns.input.LambdaName + "/handler/bootstrap.go", Data: ns.tmpl.FrameV2.Plain.Handler.BootstrapGo},
	}
}
