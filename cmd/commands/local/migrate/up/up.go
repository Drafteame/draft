package up

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var migrateUpCmd = cobra.Command{
	Use:   "local:migrate:up [flags] ",
	Short: "Execute migrate up command to apply all pending migrations",
	Long:  "Execute migrations using go migrate over all databases or a specific one",
	Run:   run,
}

func GetCmd() *cobra.Command {
	return &migrateUpCmd
}

func init() {
	migrateUpCmd.Flags().StringP("database", "D", "", "database name")
	migrateUpCmd.Flags().StringP("local-migrate-config", "c", ".local-migrate-config.yml", "path to the migrations config file")
	migrateUpCmd.Flags().Bool("all", false, "migrate all databases")
}

func run(cmd *cobra.Command, _ []string) {
	defer func() {
		if r := recover(); r != nil {
			_, _ = fmt.Printf("%v\n", r)
			os.Exit(1)
		}
	}()

	if workDir := cmd.Parent().Flag("working-dir").Value.String(); workDir != "" {
		if err := os.Chdir(workDir); err != nil {
			panic(err)
		}
	}

	config, err := loadConfig(cmd)
	if err != nil {
		panic(fmt.Sprintf("Failed to load local migrations config: %s\n", err.Error()))
	}

	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		panic(fmt.Sprintf("Failed to get all flag: %s\n", err.Error()))
	}

	if all {
		if err := migrateAll(config); err != nil {
			panic(fmt.Sprintf("Failed to migrate all databases: %s\n", err.Error()))
		}

		_, _ = fmt.Println("Migrations executed successfully")
		return
	}

	dbName := cmd.Flag("database").Value.String()

	if dbName == "" {
		name, err := promptSelectDD(config)
		if err != nil {
			panic(fmt.Sprintf("Failed to select database: %s\n", err.Error()))
		}

		dbName = name
	}

	if err := migrateOne(config, dbName); err != nil {
		panic(fmt.Sprintf("Failed to migrate database %s: %s\n", dbName, err.Error()))
	}

	_, _ = fmt.Println("Migrations executed successfully")
}
