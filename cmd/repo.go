/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Drafteame/draft/tpl"
	"github.com/Drafteame/draft/utils"
	"github.com/spf13/cobra"
)

const (
	RepoTypeMongo = "mongo"
)

type repositoryArgs struct {
	RawName        string
	CollectionName string
	PackageName    string
	CammelCaseName string
	Type           string
}

// repositoryCmd represents the repository command
var repositoryCmd = &cobra.Command{
	Use:   "repo",
	Short: "Add new repository port",
	Long: `Create a new repository port to interact with database with MongoDB, QLDB
or any other supported database.`,
	Run: func(cmd *cobra.Command, args []string) {
		a := getRepositoryArgs(cmd)

		utils.PathExists(currentDir + "/.engine.yml")

		a.validateType()
		a.validatePaths()
		a.createDirs()
		a.createFiles()

		fmt.Printf("engine: Repository `%s` created successfully\n", a.PackageName)
	},
}

func init() {
	rootCmd.AddCommand(repositoryCmd)

	repositoryCmd.Flags().StringP("name", "n", "", "Name of the repository")
	repositoryCmd.Flags().StringP("entity", "e", "", "Name of the entity (table, collection, etc) that should be managed by the repository")
	repositoryCmd.Flags().StringP("type", "t", "mongo", "Repository type to be crated")
}

func getRepositoryArgs(cmd *cobra.Command) *repositoryArgs {
	name := cmd.Flag("name").Value.String()
	repoType := cmd.Flag("type").Value.String()
	entity := cmd.Flag("entity").Value.String()

	if name == "" {
		fmt.Println("engine: --name is required")
		os.Exit(1)
	}

	if repoType == "" {
		fmt.Println("engine: --type is required")
		os.Exit(1)
	}

	if entity == "" {
		entity = utils.ToSnakeCase(name)
	}

	a := &repositoryArgs{}

	a.RawName = name
	a.Type = repoType
	a.CollectionName = entity
	a.PackageName = utils.ToPackageName(a.RawName)
	a.CammelCaseName = utils.ToCammelCase(a.RawName)

	return a
}

func (ra *repositoryArgs) validateType() {
	switch ra.Type {
	case RepoTypeMongo:
		return
	default:
		fmt.Println("engine: Unsupported repository type")
		os.Exit(1)
		return
	}
}

func (ra *repositoryArgs) validatePaths() {
	file := fmt.Sprintf("%s/internal/ports/repositories/%s", currentDir, ra.PackageName)

	utils.PathNotExists(file)
}

func (ra *repositoryArgs) createDirs() {
	path := fmt.Sprintf("%s/internal/ports/repositories/%s", currentDir, ra.PackageName)

	utils.CreateFolder(path)
}

func (ra *repositoryArgs) createFiles() {
	file := fmt.Sprintf("%s/internal/ports/repositories/%s/%s.go", currentDir, ra.PackageName, ra.PackageName)

	content, err := utils.RenderTemplate(tpl.MongoRepository, ra)
	if err != nil {
		fmt.Println("engine: Error rendering repository")
		os.Exit(1)
	}

	utils.CreateFile(file, content)
}
