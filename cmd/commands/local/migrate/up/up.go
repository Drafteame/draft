package up

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/cmd/commands/local/migrate/internal/flags"
	comm "github.com/Drafteame/draft/internal/actions/local/migrate/command"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var migrateUpCmd = &cobra.Command{
	Use:   "local:migrate:up [flags] ",
	Short: "Execute migrate up command to apply all pending migrations",
	Long:  "Execute migrations using go migrate over all databases or a specific one",
	Run:   run,
}

func GetCmd() *cobra.Command {
	return migrateUpCmd
}

func init() {
	flags.Register(migrateUpCmd)
}

func run(cmd *cobra.Command, _ []string) {
	common.ChDir(cmd)

	f := flags.GetFlags(cmd)

	command := "up"

	input := comm.Input{
		Command:            command,
		Database:           f.Database,
		LocalMigrateConfig: f.Config,
		Group:              f.Group,
		All:                f.All,
	}

	errExec := comm.New(input).Exec()
	if errExec != nil {
		log.Exitf(1, "failed to execute command: %s", errExec.Error())
	}

	log.Success("Migration up completed")
}
