/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/Drafteame/draft/tpl"
	"github.com/Drafteame/draft/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// migrate:newCmd represents the migrate:new command
var migrateNewCmd = &cobra.Command{
	Use:   "migrate:new",
	Short: "Generate new migration file",
	Long: `Generate a new migration file for the specified database.

If the database migrations folder is empty ir creates all needed structure
to generate the new migration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PathExists(currentDir + "/.engine.yml")

		a := getMigrateNewArgs(cmd)

		a.validatePaths()
		a.addMigration()

		fmt.Printf("engine:migrate: Migration file created\n\n")
		fmt.Printf("\tRun `go mod tidy` to update dependencies.\n\n")
	},
}

type migrateNewArgs struct {
	Name          string
	Namespace     string
	SnakeCaseName string
	Version       int64
	Mongo         bool
	Qldb          bool
}

func init() {
	rootCmd.AddCommand(migrateNewCmd)

	migrateNewCmd.Flags().StringP("name", "n", "", "Name of the migration")
	migrateNewCmd.Flags().Bool("mongo", false, "Specify MongoDB migration")
	migrateNewCmd.Flags().Bool("qldb", false, "Specify QLDB migration")
}

func getMigrateNewArgs(cmd *cobra.Command) *migrateNewArgs {
	name := cmd.Flag("name").Value.String()

	if name == "" {
		fmt.Println("engine:migrate: Migration name is required")
		os.Exit(1)
	}

	a := &migrateNewArgs{
		Name:          name,
		Namespace:     viper.GetString("namespace"),
		SnakeCaseName: utils.ToSnakeCase(name),
		Version:       time.Now().Unix(),
		Qldb:          cmd.Flag("qldb").Value.String() == "true",
		Mongo:         cmd.Flag("mongo").Value.String() == "true",
	}

	return a
}

func (a *migrateNewArgs) addMigration() {
	fileName := fmt.Sprintf("%d_%s.go", a.Version, a.SnakeCaseName)

	if a.Mongo {
		a.addMongoMigration(fileName)
		return
	}

	if a.Qldb {
		a.addQldbMigration(fileName)
		return
	}

	fmt.Println("engine:migrate: no database specified")
}

func (a *migrateNewArgs) addMongoMigration(fileName string) {
	path := fmt.Sprintf("%s/migrations/mongo/migrate/%s", currentDir, fileName)

	render, err := utils.RenderTemplate(tpl.MigrateMongoGo, a)
	if err != nil {
		fmt.Printf("engine:migrate: %v\n", err)
		os.Exit(1)
	}

	if err = utils.CreateFile(path, render); err != nil {
		fmt.Printf("engine:migrate: %v\n", err)
		os.Exit(1)
	}
}

func (a *migrateNewArgs) addQldbMigration(_ string) {
	fmt.Println("engine:migrate: QLDB migration not implemented yet")
	os.Exit(0)
}

func (a *migrateNewArgs) validatePaths() {
	if a.Mongo {
		a.validateMongoPaths()
		return
	}

	if a.Qldb {
		a.validateQldbPaths()
		return
	}

	fmt.Println("engine:migrate: no database specified")
}

func (a *migrateNewArgs) validateMongoPaths() {
	if utils.PathNotExists(currentDir + "/migrations/mongo/migrate") {
		if err := utils.CreateFolder(currentDir + "/migrations/mongo/migrate"); err != nil {
			fmt.Printf("engine:migrate: %v\n", err)
			os.Exit(1)
		}
	}

	if utils.PathNotExists(currentDir + "/migrations/mongo/main.go") {
		render, err := utils.RenderTemplate(tpl.MigrateMongoMainGo, a)
		if err != nil {
			fmt.Printf("engine:migrate: %v\n", err)
			os.Exit(1)
		}

		if err = utils.CreateFile(currentDir+"/migrations/mongo/main.go", render); err != nil {
			fmt.Printf("engine:migrate: %v\n", err)
			os.Exit(1)
		}
	}
}

func (a *migrateNewArgs) validateQldbPaths() {
	if utils.PathNotExists(currentDir + "/migrations/qldb/migrate") {
		if err := utils.CreateFolder(currentDir + "/migrations/mongo"); err != nil {
			fmt.Printf("engine:migrate: %v\n", err)
			os.Exit(1)
		}
	}

	if utils.PathNotExists(currentDir + "/migrations/qldb/main.go") {
		render, err := utils.RenderTemplate(tpl.MigrateQLDBMainGo, a)
		if err != nil {
			fmt.Printf("engine:migrate: %v\n", err)
			os.Exit(1)
		}

		if err = utils.CreateFile(currentDir+"/migrations/qldb/main.go", render); err != nil {
			fmt.Printf("engine:migrate: %v\n", err)
			os.Exit(1)
		}
	}
}
