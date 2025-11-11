package newdomain

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh/spinner"
	"gopkg.in/yaml.v3"

	"github.com/Drafteame/draft/internal/pkg/exec"
	"github.com/Drafteame/draft/internal/pkg/files"
)

func (nd *NewDomain) mockery() error {
	var err error

	spin := spinner.New().Title("Adding mockery configs")

	action := func() {
		spin.Update("Creating mockery package")
		err = nd.addMockeryPackages()
		if err != nil {
			return
		}

		spin.Update("Running mockery to create files")
		err = nd.createMockeryFiles()
		if err != nil {
			return
		}
	}

	spinErr := spin.Action(action).Run()

	return errors.Join(spinErr, err)
}

func (nd *NewDomain) addMockeryPackages() error {
	paths := map[string]string{
		nd.input.PackageName + "/" + nd.input.DomainPath + "/service":    nd.input.DomainPath + "/service/mocks",
		nd.input.PackageName + "/" + nd.input.DomainPath + "/repository": nd.input.DomainPath + "/repository/mocks",
	}

	mockeryConfig, err := files.Read(".mockery.yml")
	if err != nil {
		return err
	}

	config := map[string]any{}

	if err := yaml.Unmarshal(mockeryConfig, &config); err != nil {
		return err
	}

	packages, ok := config["packages"].(map[string]any)
	if !ok {
		return fmt.Errorf("packages key not found or invalid in .mockery.yml")
	}

	for pkgPath, mocksDir := range paths {
		packages[pkgPath] = map[string]any{
			"config": map[string]any{
				"all":       true,
				"dir":       mocksDir,
				"filename":  "mock_{{.InterfaceName}}.go",
				"inpackage": false,
			},
		}
	}

	newConfig, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return files.Create(".mockery.yml", newConfig)
}

func (nd *NewDomain) createMockeryFiles() error {
	_, err := exec.Command("mockery")
	if err != nil {
		return fmt.Errorf("command 'mockery' failed: %w", err)
	}

	return nil
}
