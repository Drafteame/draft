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

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new project using engine structure",
	Long: `This command will create a new project using the engine structure and will
setup all basic files and folder structure to start developing.

This is the project structure to be created by the command:

<name>
├── .github
│   └── workflows
|       ├── dev_deploy.yml
|       ├── manual_redeploy.yml
|       ├── pull_request.yml
│       └── release.yml
├── cmd
│   └── api
│       └── main.go
├── config
│   ├── concurrency.yml
│   ├── cors.yml
│   ├── domains.yml
│   ├── environment.yml
│   ├── iam.yml
│   └── vpc.yml
├── internal
│   ├── adapters
│   │   └── .gitkeep
│   ├── handlers
│   │   └── .gitkeep
│   ├── models
│   │   └── .gitkeep
│   ├── ports
│   │   └── repositories
│   │       └── .gitkeep
│   ├── routes
│   │   └── routes.go
│   ├── schemas
│   │   └── .gitkeep
├── migrations
│   └── .gitkeep
├── .cz.toml
├── .engine.yml
├── .gitignore
├── .golangci.yml
├── .goreleaser.yml
├── build.sh
├── go.mod
├── Makefile
├── package.json
├── README.md
└── serverless.yml
- `,
	Run: func(cmd *cobra.Command, args []string) {
		a := getNewArgs(cmd)

		fmt.Printf("engine: Generating %s project on %s\n", a.Name, currentDir+"/"+a.Name)

		a.createDirs()
		a.createFiles()

		fmt.Printf("engine: Project %s created successfully. \n\n\tRun `go mod tidy` to download dependencies.\n", a.Name)
	},
}

type newArgs struct {
	Name               string
	Namespace          string
	ReleaseTagReplacer string
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringP("name", "n", "", "Name of the project")
	newCmd.Flags().StringP("namespace", "s", "", "Namespace that should be used by go packages. If is not present will be automatically set according to the project name")
}

func getNewArgs(cmd *cobra.Command) *newArgs {
	name := cmd.Flag("name").Value.String()
	namespace := cmd.Flag("namespace").Value.String()

	if name == "" {
		fmt.Println("engine: --name is required")
		os.Exit(1)
	}

	if namespace == "" {
		namespace = utils.ToPackageName(name)
	}

	return &newArgs{
		Name:               name,
		Namespace:          namespace,
		ReleaseTagReplacer: "{{ .Tag }}",
	}
}

func (args *newArgs) createDirs() {
	args.createRootFolder()

	path := currentDir + "/" + args.Name

	dirs := []string{
		".github/workflows",
		"cmd/api",
		"config",
		"migrations",
		"internal/adapters",
		"internal/handlers",
		"internal/models",
		"internal/routes",
		"internal/schemas",
		"internal/ports/repositories",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(path+"/"+dir, 0755); err != nil {
			fmt.Printf("engine: error creating directory %s: \n", path+"/"+dir)
			fmt.Printf("engine: %v\n", err)
			os.Exit(1)
		}
	}
}

func (args *newArgs) createRootFolder() {
	path := currentDir + "/" + args.Name

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		fmt.Printf("engine: project %s already exists\n", path)
		os.Exit(1)
	}

	err := os.Mkdir(path, 0755)
	if err != nil {
		fmt.Printf("engine: error creating directory: %s\n", path)
		os.Exit(1)
	}
}

func (args *newArgs) createFiles() {
	path := currentDir + "/" + args.Name

	files := [][]string{
		{".github/workflows/dev_deploy.yml", tpl.GithubActionDevDeploy},
		{".github/workflows/manual_redeploy.yml", tpl.GithubActionManualRedeploy},
		{".github/workflows/pull_request.yml", tpl.GithubActionPullRequest},
		{".github/workflows/release.yml", tpl.GithubActionRelease},
		{"cmd/api/main.go", tpl.HTTPMainGo},
		{"config/concurrency.yml", tpl.ConcurrencyYaml},
		{"config/cors.yml", tpl.CorsYaml},
		{"config/domains.yml", tpl.DomainsYaml},
		{"config/environment.yml", tpl.EnvironmentYaml},
		{"config/iam.yml", tpl.IamYaml},
		{"config/vpc.yml", tpl.VpcYaml},
		{"internal/adapters/.gitkeep", tpl.GitKeep},
		{"internal/handlers/.gitkeep", tpl.GitKeep},
		{"internal/models/.gitkeep", tpl.GitKeep},
		{"internal/ports/repositories/.gitkeep", tpl.GitKeep},
		{"internal/routes/routes.go", tpl.RoutesGo},
		{"internal/schemas/.gitkeep", tpl.GitKeep},
		{"migrations/.gitkeep", tpl.GitKeep},
		{".cz.toml", tpl.CzToml},
		{".engine.yml", tpl.EngineYaml},
		{".gitignore", tpl.GitIgnore},
		{".golangci.yml", tpl.GolangCIYaml},
		{".goreleaser.yml", tpl.GoreleaserYaml},
		{"build.sh", tpl.BuildSh},
		{"go.mod", tpl.GoMod},
		{"Makefile", tpl.Makefile},
		{"package.json", tpl.PackageJSON},
		{"README.md", tpl.ReadmeMd},
		{"serverless.yml", tpl.ServerlessYaml},
	}

	for _, file := range files {
		tmpl, err := utils.RenderTemplate(file[1], args)
		if err != nil {
			fmt.Printf("engine: error rendering template...\n")
			fmt.Printf("engine: %v\n", err)
			os.Exit(1)
		}

		fileName := path + "/" + file[0]

		if err = utils.CreateFile(fileName, tmpl); err != nil {
			fmt.Printf("engine: error writing file %s\n", fileName)
			fmt.Printf("engine: %v\n", err)
			os.Exit(1)
		}
	}
}
