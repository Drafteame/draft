package newdomain

import (
	"strings"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"

	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/pkg/files"
)

func (nd *NewDomain) postCreate() error {
	actions := []func() error{
		nd.addMockeryPackages,
		nd.addPostgresTestModelsOnProvider,
	}

	for _, action := range actions {
		if err := action(); err != nil {
			return err
		}
	}

	return nil
}

func (nd *NewDomain) addMockeryPackages() error {
	paths := []string{
		nd.input.PackageName + "/domains/" + nd.input.DomainPath + "/service",
		nd.input.PackageName + "/domains/" + nd.input.DomainPath + "/repository",
	}

	mockeryConfig, err := files.Read(".mockery.yml")
	if err != nil {
		return err
	}

	config := map[string]any{}

	if err := yaml.Unmarshal(mockeryConfig, &config); err != nil {
		return err
	}

	for _, path := range paths {
		config["packages"].(map[string]any)[path] = map[string]any{}
	}

	newConfig, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return files.Create(".mockery.yml", newConfig)
}

func (nd *NewDomain) addPostgresTestModelsOnProvider() error {
	fileName := lo.SnakeCase(nd.input.DBName) + ".go"
	filePath := "pkg/providers/postgres/" + fileName

	daoPackage := nd.input.PackageName + "/domains/" + nd.input.DomainPath + "/repository/daos"
	alias := "dao" + nd.input.DomainNameLower
	fullImport := alias + ` "` + daoPackage + `"`
	fullModel := alias + "." + nd.input.DomainNamePascal + "{},"

	replacers := map[string]string{
		data.NextImportTag:  fullImport + "\n\t" + data.NextImportTag,
		data.NextDbModelTag: fullModel + "\n\t" + data.NextDbModelTag,
	}

	return files.ScanAndWrite(filePath, func(line string) (string, error) {
		for key, val := range replacers {
			if strings.Contains(line, key) {
				line = strings.ReplaceAll(line, key, val)
			}
		}

		return line, nil
	})
}
