package down

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/cmd/commands/local/migrate/internal/flags"
	comm "github.com/Drafteame/draft/internal/actions/local/migrate/command"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var migrateDownCmd = &cobra.Command{
	Use:   "local:migrate:down [flags] [number]",
	Short: "Execute migrate down command to rollback the last migrations",
	Long:  "Rollback one or more database migrations using go migrate on a specific database",
	Run:   run,
}

func GetCmd() *cobra.Command {
	return migrateDownCmd
}

func init() {
	flags.Register(migrateDownCmd)
}

func run(cmd *cobra.Command, args []string) {
	common.ChDir(cmd)

	var (
		numMigrations int64
		err           error
	)

	if len(args) > 0 {
		numMigrations, err = strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Exitf(1, "invalid number of migrations: %s", err.Error())
		}
	}

	f := flags.GetFlags(cmd)

	command := "down"

	input := comm.Input{
		Command:            command,
		Database:           f.Database,
		LocalMigrateConfig: f.Config,
		Group:              f.Group,
		All:                f.All,
		NumberMigrations:   numMigrations,
	}

	errExec := comm.New(input).Exec()
	if errExec != nil {
		log.Exitf(1, "failed to execute command: %s", errExec.Error())
	}

	log.Success("Migration down completed")
}
