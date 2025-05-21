package force

import (
	"fmt"
	comm "github.com/Drafteame/draft/internal/actions/local/migrate/command"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var migrateForceCmd = cobra.Command{
	Use:   "local:migrate:force [flags] [version]",
	Short: "Force a specific migration version to resolve dirty state",
	Long:  "When testing migrations, inconsistencies can occur. This command forces a specific version to avoid dirty state and continue with migration corrections.",
	Args:  cobra.ExactArgs(1),
	Run:   run,
}

func GetCmd() *cobra.Command {
	return &migrateForceCmd
}

func init() {
	migrateForceCmd.Flags().StringP("database", "D", "", "database name")
	migrateForceCmd.Flags().StringP("local-migrate-config", "c", ".local-migrate-config.yml", "path to the migrations config file")
	migrateForceCmd.Flags().String("group", "", "DB migrations group name")
}

func run(cmd *cobra.Command, args []string) {
	defer func() {
		if r := recover(); r != nil {
			_, _ = fmt.Printf("%v\n", r)
			os.Exit(1)
		}
	}()

	version, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		panic(fmt.Errorf("invalid version number: %w", err))
	}

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

	group, err := cmd.Flags().GetString("group")
	if err != nil {
		panic(err)
	}

	command := "force"

	input := comm.Input{
		Command:            command,
		WorkingDir:         workingDir,
		Database:           database,
		LocalMigrateConfig: localMigrateConfig,
		Group:              group,
		All:                false,
		Version:            version,
	}

	errExec := comm.New(input).Exec()
	if errExec != nil {
		panic(errExec)
	}
}
