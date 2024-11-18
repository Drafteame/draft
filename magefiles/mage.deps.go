package main

import (
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"magefiles/files"
)

// Deps namespace for dependency management
type Deps mg.Namespace

type packageJSON struct {
	Name            string            `json:"name"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

type slsEntry struct {
	workingDir string
	deps       packageJSON
}

// InstallHooks install Husky hooks.
func (Deps) InstallHooks() error {
	return sh.Run("husky", "install")
}

// UpgradeSls upgrade serverless dependencies on all services.
func (d Deps) UpgradeSls() error {
	fileList, err := d.getSlsFileList()
	if err != nil {
		return err
	}

	entries, err := d.getSlsEntries(fileList)
	if err != nil {
		return err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if err := d.upgradeSls(pwd, entry); err != nil {
			return err
		}
	}

	return nil
}

func (Deps) installSls(resetPath string, entry slsEntry) error {
	defer func() {
		if err := os.Chdir(resetPath); err != nil {
			panic(err)
		}
	}()

	if err := os.Chdir(entry.workingDir); err != nil {
		panic(err)
	}

	println("==== Installing dependencies for service", entry.deps.Name)

	return sh.Run("npm", "install")
}

func (d Deps) upgradeSls(resetPath string, entry slsEntry) error {
	defer func() {
		if err := os.Chdir(resetPath); err != nil {
			panic(err)
		}
	}()

	if err := os.Chdir(entry.workingDir); err != nil {
		panic(err)
	}

	println("==== Upgrading dependencies for service", entry.deps.Name)

	if err := d.upgradeSlsDevDeps(entry); err != nil {
		return err
	}

	return d.upgradeSlsDeps(entry)
}

func (Deps) upgradeSlsDevDeps(entry slsEntry) error {
	println("---- Dev dependencies")

	for devDep := range entry.deps.DevDependencies {
		println(devDep)

		if err := sh.Run("npm", "install", devDep+"@latest", "--save-dev"); err != nil {
			return err
		}
	}

	return nil
}

func (Deps) upgradeSlsDeps(entry slsEntry) error {
	println("---- Dependencies")

	for dep := range entry.deps.Dependencies {
		println(dep)

		if err := sh.Run("npm", "install", dep+"@latest", "--save-dev"); err != nil {
			return err
		}
	}

	return nil
}

func (Deps) getSlsFileList() ([]string, error) {
	omit := []string{
		".serverless",
		".bin",
		".warmup",
		"node_modules",
		"config",
		"http",
		"sqs",
		"snssqs",
		"plain",
		"kafka",
	}

	return files.Search(".", "package.json", files.WithOmit(omit...))
}

func (Deps) getSlsEntries(fileList []string) ([]slsEntry, error) {
	entries := make([]slsEntry, 0, len(fileList))

	for _, filePath := range fileList {
		content := packageJSON{}

		if err := files.BindJSON(filePath, &content); err != nil {
			return nil, err
		}

		entries = append(entries, slsEntry{
			workingDir: filepath.Dir(filePath),
			deps:       content,
		})
	}

	return entries, nil
}
