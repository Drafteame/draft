/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Drafteame/draft/utils"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateDownCmd = &cobra.Command{
	Use:   "migrate:down",
	Short: "Execute migrations DOWN command",
	Long:  `Makes a system call to execute down migrations accorgind configurations.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PathExists(currentDir + "/.engine.yml")

		a := getMigrateDownArgs(cmd)

		a.executeMigrationsDown()

		fmt.Println("engine:migrate: migrations DOWN and offline...")
	},
}

type migrateDownArgs struct {
	Mongo bool
	Qldb  bool
}

func init() {
	rootCmd.AddCommand(migrateDownCmd)

	migrateDownCmd.Flags().Bool("mongo", false, "Specify MongoDB migration")
	migrateDownCmd.Flags().Bool("qldb", false, "Specify QLDB migration")
}

func getMigrateDownArgs(cmd *cobra.Command) *migrateDownArgs {
	return &migrateDownArgs{
		Mongo: cmd.Flag("mongo").Value.String() == "true",
		Qldb:  cmd.Flag("qldb").Value.String() == "true",
	}
}

func (a *migrateDownArgs) executeMigrationsDown() {
	if a.Mongo {
		executeMongoMigrationsDown()
		return
	}

	if a.Qldb {
		executeQldbMigrationsDown()
		return
	}

	fmt.Println("engine:migrate: No database migration specified")
	os.Exit(1)
}

func executeMongoMigrationsDown() {
	err := utils.Run("go", "run", fmt.Sprintf("%s/migrations/mongo/main.go", currentDir), "down")
	if err != nil {
		fmt.Println("engine:migrate: Error executing MongoDB migration")
		fmt.Printf("engine:migrate: %v\n", err)
		os.Exit(1)
	}
}

func executeQldbMigrationsDown() {
	err := utils.Run("go", "run", fmt.Sprintf("%s/migrations/qldb/main.go", currentDir), "down")
	if err != nil {
		fmt.Println("engine:migrate: Error executing QLDB migration")
		fmt.Printf("engine:migrate: %v\n", err)
		os.Exit(1)
	}
}
