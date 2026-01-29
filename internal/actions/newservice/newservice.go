package newservice

import (
	"errors"
	"os"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/dirs"
	"github.com/Drafteame/draft/internal/templates"
)

type NewService struct {
	tmpl  *templates.ServiceTemplates
	input dtos.ServiceInput
}

func New(input dtos.ServiceInput) *NewService {
	return &NewService{input: input}
}

func (ns *NewService) Exec() error {
	if errPre := ns.preCreate(); errPre != nil {
		return errPre
	}

	tmpl, err := templates.NewServiceTemplates(ns.input)
	if err != nil {
		return err
	}

	ns.tmpl = tmpl

	if errExec := ns.exec(); errExec != nil {
		return errExec
	}

	return ns.postCreate()
}

func (ns *NewService) exec() error {
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
		ns.input.ServicePath + "/config/app",
		ns.input.ServicePath + "/config/sls",
	}

	return dirs.Create(folders...)
}

func (ns *NewService) createLambdaFolders() error {
	lambdaPath := ns.input.ServicePath + "/cmd/" + ns.input.LambdaType + "/" + ns.input.LambdaName

	if ns.input.IsLegacy {
		lambdaPath = ns.input.ServicePath + "/" + ns.input.LambdaType + "/" + ns.input.LambdaName
	}

	folders := []string{
		lambdaPath + "/handler",
		lambdaPath + "/handler/worker",
		lambdaPath + "/handler/embed",
		lambdaPath + "/handler/dtos",
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
		{Path: "/config/sls/resources.yml", Data: ns.tmpl.Config.Sls.ResourcesYAML},
	}

	entries = append(entries, ns.getEntries()...)

	return entries
}

func (ns *NewService) getEntries() []dtos.FileEntry {
	var entries []dtos.FileEntry

	lambdaPath := "/cmd/plain/"
	if ns.input.IsLegacy {
		lambdaPath = "/plain/"
	}

	entries = []dtos.FileEntry{
		{Path: lambdaPath + ns.input.LambdaName + "/main.go", Data: ns.tmpl.Lambda.Plain.MainGo},
		{Path: lambdaPath + ns.input.LambdaName + "/lambda-config.yml", Data: ns.tmpl.Lambda.Plain.LambdaConfigYAML},
		{Path: lambdaPath + ns.input.LambdaName + "/handler/bootstrap.go", Data: ns.tmpl.Lambda.Plain.Handler.BootstrapGo},
		{Path: lambdaPath + ns.input.LambdaName + "/handler/worker/worker.go", Data: ns.tmpl.Lambda.Plain.Handler.WorkerGo},
		{Path: lambdaPath + ns.input.LambdaName + "/handler/worker/resources.go", Data: ns.tmpl.Lambda.Plain.Handler.ResourcesGo},
		{Path: lambdaPath + ns.input.LambdaName + "/handler/dtos/dto.go", Data: ns.tmpl.Lambda.Plain.Handler.DtosGo},
		{Path: lambdaPath + ns.input.LambdaName + "/handler/embed/_.yaml", Data: ns.tmpl.Lambda.Plain.Handler.EmbedYML},
	}
	return entries
}
