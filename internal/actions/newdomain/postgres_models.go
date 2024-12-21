package newdomain

import (
	"errors"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/samber/lo"

	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/pkg/files"
)

func (nd *NewDomain) postgresModels() error {
	var err error

	spin := spinner.New().Title("Adding domain doas to provider test migrations")

	action := func() {
		err = nd.addPostgresTestModelsOnProvider()
	}

	spinErr := spin.Action(action).Run()

	return errors.Join(spinErr, err)
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

	content, err := files.ReadString(filePath)
	if err != nil {
		return err
	}

	for key, value := range replacers {
		content = strings.ReplaceAll(content, key, value)
	}

	return files.Create(filePath, []byte(content))
}
