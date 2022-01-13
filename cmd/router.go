/*
Copyright © 2021 Drafteame eduardo.aguilar@draftea.com

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
	"strings"

	"github.com/Drafteame/draft/tpl"
	"github.com/Drafteame/draft/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type routerArgs struct {
	RawName        string
	PackageName    string
	CammelCaseName string
	SnakeCaseName  string
	Name           string
	Namespace      string
}

// routerCmd represents the router command
var routerCmd = &cobra.Command{
	Use:   "router",
	Short: "Add New router to the project",
	Long: `Generates all needed files to add a new router and a sample path to the
project.`,
	Run: func(cmd *cobra.Command, args []string) {
		a := getRouterArgs(cmd)

		utils.PathExists(currentDir + "/.engine.yml")

		a.validatePaths()
		a.createFiles()
		a.registerRoutes()

		fmt.Printf("engine: router `%s` created\n", a.PackageName)
	},
}

func init() {
	rootCmd.AddCommand(routerCmd)

	routerCmd.Flags().StringP("name", "n", "", "Name of the router")
}

func getRouterArgs(cmd *cobra.Command) *routerArgs {
	name := cmd.Flag("name").Value.String()

	if name == "" {
		fmt.Println("engine: --name is required")
		os.Exit(1)
	}

	return &routerArgs{
		RawName:        name,
		PackageName:    utils.ToPackageName(name),
		CammelCaseName: utils.ToCammelCase(name),
		SnakeCaseName:  utils.ToSnakeCase(name),
		Name:           viper.GetString("name"),
		Namespace:      viper.GetString("namespace"),
	}
}

func (args *routerArgs) validatePaths() {
	pathsExists := []string{
		handlersPath,
		routesPath,
		schemasPath,
	}

	for _, path := range pathsExists {
		if utils.PathNotExists(path) {
			fmt.Printf("engine: %s not found.\n", path)
			os.Exit(1)
		}
	}

	pathsNotExists := []string{
		fmt.Sprintf("%s/%s", handlersPath, args.PackageName),
		fmt.Sprintf("%s/%s/handler.go", handlersPath, args.PackageName),
		fmt.Sprintf("%s/%s/%s.go", handlersPath, args.PackageName, args.PackageName),
		fmt.Sprintf("%s/%s.go", routesPath, args.PackageName),
		fmt.Sprintf("%s/%s.go", schemasPath, args.PackageName),
	}

	for _, path := range pathsNotExists {
		if utils.PathNotExists(path) {
			fmt.Printf("engine: %s alredy exists.\n", path)
			os.Exit(1)
		}
	}
}

func (args *routerArgs) createFiles() {
	if err := utils.CreateFolder(fmt.Sprintf("%s/%s", handlersPath, args.PackageName)); err != nil {
		fmt.Printf("engine: error creating handler folder at %s/%s\n", handlersPath, args.PackageName)
		os.Exit(1)
	}

	files := [][]string{
		{fmt.Sprintf("%s/%s.go", schemasPath, args.PackageName), tpl.JSONSchemas},
		{fmt.Sprintf("%s/%s/%s.go", handlersPath, args.PackageName, args.PackageName), tpl.HandlerInterface},
		{fmt.Sprintf("%s/%s/handler.go", handlersPath, args.PackageName), tpl.HandlerStruct},
		{fmt.Sprintf("%s/%s.go", routesPath, args.PackageName), tpl.Router},
	}

	for _, file := range files {
		tmpl, err := utils.RenderTemplate(file[1], args)
		if err != nil {
			fmt.Printf("engine: error rendering template: %s\n", err.Error())
			os.Exit(1)
		}

		if err = utils.CreateFile(file[0], tmpl); err != nil {
			fmt.Printf("engine: error writing file %s: %s\n", file[0], err.Error())
			os.Exit(1)
		}
	}
}

func (args *routerArgs) registerRoutes() {
	mainRouterPath := fmt.Sprintf("%s/routes.go", routesPath)
	router := utils.ReadFile(mainRouterPath)

	tmpl, err := utils.RenderTemplate(tpl.RegisterRouter, args)
	if err != nil {
		fmt.Printf("engine: error rendering template: %s\n", err.Error())
		os.Exit(1)
	}

	router = strings.ReplaceAll(router, "// router:register", tmpl)

	utils.ReplaceFileContent(mainRouterPath, router)
}
