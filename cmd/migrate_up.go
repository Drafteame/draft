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
var migrateUpCmd = &cobra.Command{
	Use:   "migrate:up",
	Short: "Execute migrations UP command",
	Long:  `Makes a system call to execute up migrations accorgind configurations.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PathExists(currentDir + "/.engine.yml")

		a := getMigrateUpArgs(cmd)

		a.executeMigrationsUp()

		fmt.Println("engine:migrate: migrations UP and ready...")
	},
}

type migrateUpArgs struct {
	Mongo bool
	Qldb  bool
}

func init() {
	rootCmd.AddCommand(migrateUpCmd)

	migrateUpCmd.Flags().Bool("mongo", false, "Specify MongoDB migration")
	migrateUpCmd.Flags().Bool("qldb", false, "Specify QLDB migration")
}

func getMigrateUpArgs(cmd *cobra.Command) *migrateUpArgs {
	return &migrateUpArgs{
		Mongo: cmd.Flag("mongo").Value.String() == "true",
		Qldb:  cmd.Flag("qldb").Value.String() == "true",
	}
}

func (a *migrateUpArgs) executeMigrationsUp() {
	if a.Mongo {
		executeMongoMigrationsUp()
		return
	}

	if a.Qldb {
		executeQldbMigrationsUp()
		return
	}

	fmt.Println("engine:migrate: No database migration specified")
	os.Exit(1)
}

func executeMongoMigrationsUp() {
	err := utils.Run("go", "run", fmt.Sprintf("%s/migrations/mongo/main.go", currentDir), "up")
	if err != nil {
		fmt.Println("engine:migrate: Error executing MongoDB migration")
		fmt.Printf("engine:migrate: %v\n", err)
		os.Exit(1)
	}
}

func executeQldbMigrationsUp() {
	err := utils.Run("go", "run", fmt.Sprintf("%s/migrations/qldb/main.go", currentDir), "up")
	if err != nil {
		fmt.Println("engine:migrate: Error executing QLDB migration")
		fmt.Printf("engine:migrate: %v\n", err)
		os.Exit(1)
	}
}
