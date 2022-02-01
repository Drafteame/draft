/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Drafteame/draft/tpl"
	"github.com/Drafteame/draft/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// adapterCmd represents the adapter command
var adapterCmd = &cobra.Command{
	Use:   "adapter",
	Short: "Adds a simple adapter interface and startup configuration",
	Long: `Adds the adapter configuration needed to be used by all handlers.

/internal/adapters/<adapter_name>
├── <adapter_name>.go
└── adapter.go
`,
	Run: func(cmd *cobra.Command, args []string) {
		a := new(adapterArgs)

		if err := a.validate(cmd); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := a.createFolder(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := a.createFiles(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("engine: Adapter %s created successfully. \n\n", a.Name)
		fmt.Printf("\tRun `go mod tidy` to download dependencies.\n")
	},
}

func init() {
	rootCmd.AddCommand(adapterCmd)

	adapterCmd.Flags().StringP("name", "n", "", "Name of the adapter")
	adapterCmd.Flags().StringP("model", "m", "", "Related model name of the adapter")
	adapterCmd.Flags().StringP("repo", "r", "", "Related repository package")
}

type adapterArgs struct {
	Name        string
	Namespace   string
	PackageName string
	RepoName    string
	ModelName   string
}

func (a *adapterArgs) validate(cmd *cobra.Command) error {
	a.Namespace = viper.GetString("namespace")
	a.Name = cmd.Flag("name").Value.String()
	a.PackageName = utils.ToPackageName(a.Name)
	a.RepoName = cmd.Flag("repo").Value.String()
	a.ModelName = cmd.Flag("model").Value.String()

	if a.Name == "" {
		return fmt.Errorf("engine: --name is required")
	}

	if a.RepoName == "" {
		return fmt.Errorf("engine: --repo is required")
	}

	if a.ModelName == "" {
		return fmt.Errorf("engine: --model is required")
	}

	if utils.PathNotExists(fmt.Sprintf("%s/internal/ports/repositories/%s", currentDir, a.RepoName)) {
		return fmt.Errorf("engine: %s/internal/ports/repositories/%s does not exist", currentDir, a.RepoName)
	}

	return nil
}

func (a *adapterArgs) createFolder() error {
	path := fmt.Sprintf("%s/internal/adapters/%s", currentDir, a.PackageName)

	if utils.PathExists(path) {
		return fmt.Errorf("engine: %s/internal/adapters/%s already exists", currentDir, a.PackageName)
	}

	if err := utils.CreateFolder(path); err != nil {
		return fmt.Errorf("engine: %s", err)
	}

	return nil
}

func (a *adapterArgs) createFiles() error {
	files := [][]string{
		{
			fmt.Sprintf("%s/internal/adapters/%s/%s.go", currentDir, a.PackageName, a.PackageName),
			tpl.AdapterInterface,
		},
		{
			fmt.Sprintf("%s/internal/adapters/%s/adapter.go", currentDir, a.PackageName),
			tpl.AdapterMongoImplementation,
		},
	}

	for _, file := range files {
		render, err := utils.RenderTemplate(file[1], a)
		if err != nil {
			return fmt.Errorf("engine: %v", err)
		}

		if err = utils.CreateFile(file[0], render); err != nil {
			return fmt.Errorf("engine: %s", err)
		}
	}

	return nil
}
