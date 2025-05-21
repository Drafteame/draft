package up

import (
	"fmt"
	comm "github.com/Drafteame/draft/internal/actions/local/migrate/command"
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
	migrateUpCmd.Flags().String("group", "", "DB migrations group name")
}

func run(cmd *cobra.Command, _ []string) {
	defer func() {
		if r := recover(); r != nil {
			_, _ = fmt.Printf("%v\n", r)
			os.Exit(1)
		}
	}()

	workingDir, err := cmd.Parent().Flags().GetString("working-dir")
	if err != nil {
		panic(err)
	}

	database, err := cmd.Flags().GetString("database")
	if err != nil {
		panic(err)
	}

	localMigrateConfig, err := cmd.Flags().GetString("local-migrate-config")
	if err != nil {
		panic(err)
	}

	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		panic(err)
	}

	group, err := cmd.Flags().GetString("group")
	if err != nil {
		panic(err)
	}

	command := "up"

	input := comm.Input{
		Command:            command,
		WorkingDir:         workingDir,
		Database:           database,
		LocalMigrateConfig: localMigrateConfig,
		Group:              group,
		All:                all,
	}

	errExec := comm.New(input).Exec()
	if errExec != nil {
		panic(errExec)
	}
}
