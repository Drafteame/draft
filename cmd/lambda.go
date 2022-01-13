/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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

// lambdaCmd represents the lambda command
var lambdaCmd = &cobra.Command{
	Use:   "lambda",
	Short: "Add a new Lambda function with an event handler",
	Long: `Generate all needed files to register a new Lambda function
with an event handler like sqs, sns, plain lambda etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		a := getLambdaArgs(cmd)

		utils.PathExists(currentDir + "/.engine.yml")

		fmt.Println("engine: Creating Lambda function", a.Name)

		a.createDirs()
		a.createFiles()
		a.registerEvent()

		fmt.Println("engine: Lambda function created")
		fmt.Println("\n\t Run `go mod tidy` to update dependencies.")
	},
}

func init() {
	rootCmd.AddCommand(lambdaCmd)

	lambdaCmd.Flags().StringP("name", "n", "", "Name of the Lambda function")
	lambdaCmd.Flags().StringP("type", "t", "", "Type of the Lambda function (sqs, sns, plain)")
}

type lambdaArgs struct {
	Name           string
	Type           string
	Namespace      string
	SnakeCaseName  string
	CammelCaseName string
	PackageName    string
}

func getLambdaArgs(cmd *cobra.Command) *lambdaArgs {
	name := cmd.Flag("name").Value.String()
	typeLambda := cmd.Flag("type").Value.String()

	if name == "" {
		fmt.Println("engine: --name is required")
		os.Exit(1)
	}

	if typeLambda == "" {
		typeLambda = "plain"
	}

	return &lambdaArgs{
		Name:           name,
		Type:           typeLambda,
		Namespace:      viper.GetString("namespace"),
		SnakeCaseName:  utils.ToSnakeCase(name),
		CammelCaseName: utils.ToCammelCase(name),
		PackageName:    utils.ToPackageName(name),
	}
}

func (la *lambdaArgs) createDirs() {
	dirs := []string{
		fmt.Sprintf("%s/internal/handlers/%s", currentDir, la.PackageName),
		fmt.Sprintf("%s/cmd/%s", currentDir, la.PackageName),
	}

	for _, dir := range dirs {
		utils.CreateFolder(dir)
	}
}

func (la *lambdaArgs) createFiles() {
	files := []struct {
		name string
		path string
	}{
		{
			name: "handler",
			path: fmt.Sprintf("%s/internal/handlers/%s/%s.go", currentDir, la.PackageName, la.PackageName),
		},
		{
			name: "main",
			path: fmt.Sprintf("%s/cmd/%s/main.go", currentDir, la.PackageName),
		},
	}

	for _, file := range files {
		template, err := la.getLambdaTemplate(file.name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		template, err = utils.RenderTemplate(template, la)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		utils.CreateFile(file.path, template)
	}
}

func (la *lambdaArgs) registerEvent() {
	template, err := la.getEventTemplate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	template, err = utils.RenderTemplate(template, la)
	if err != nil {
		fmt.Printf("engine: %v\n", err)
		os.Exit(1)
	}

	path := fmt.Sprintf("%s/serverless.yml", currentDir)

	content := strings.ReplaceAll(utils.ReadFile(path), "# engine:serverless:functions", template)

	utils.ReplaceFileContent(path, content)
}

func (la *lambdaArgs) getLambdaTemplate(name string) (string, error) {
	switch name {
	case "handler":
		return la.getHandlerTemplate()
	case "main":
		return la.getMainTemplate()
	case "serverless":
		return la.getEventTemplate()
	default:
		return "", fmt.Errorf("engine: template '%s' is not supported", name)
	}
}

func (la *lambdaArgs) getHandlerTemplate() (string, error) {
	switch la.Type {
	case "plain":
		return tpl.PlainHandlerGo, nil
	case "sqs":
		return tpl.SqsHandlerGo, nil
	default:
		return "", fmt.Errorf("engine: lambda type '%s' is not supported", la.Type)
	}
}

func (la *lambdaArgs) getMainTemplate() (string, error) {
	switch la.Type {
	case "plain":
		return tpl.PlainMainGo, nil
	case "sqs":
		return tpl.SqsMainGo, nil
	default:
		return "", fmt.Errorf("engine: lambda type '%s' is not supported", la.Type)
	}
}

func (la *lambdaArgs) getEventTemplate() (string, error) {
	switch la.Type {
	case "plain":
		return tpl.ServerlessPlainEvent, nil
	case "sqs":
		return tpl.ServerlessSQSEvent, nil
	default:
		return "", fmt.Errorf("engine: lambda type '%s' is not supported", la.Type)
	}
}
