package force

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/cmd/commands/local/migrate/internal/flags"
	comm "github.com/Drafteame/draft/internal/actions/local/migrate/command"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var migrateForceCmd = &cobra.Command{
	Use:   "local:migrate:force [flags] [version]",
	Short: "Force a specific migration version to resolve dirty state",
	Long:  "This command forces a specific version to avoid dirty state and continue with migration corrections.",
	Args:  cobra.ExactArgs(1),
	Run:   run,
}

func GetCmd() *cobra.Command {
	return migrateForceCmd
}

func init() {
	flags.Register(migrateForceCmd)
}

func run(cmd *cobra.Command, args []string) {
	common.ChDir(cmd)

	version, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		log.Exitf(1, "invalid version number: %s", err.Error())
	}

	f := flags.GetFlags(cmd)

	command := "force"

	input := comm.Input{
		Command:            command,
		Database:           f.Database,
		LocalMigrateConfig: f.Config,
		Group:              f.Group,
		All:                false,
		Version:            version,
	}

	errExec := comm.New(input).Exec()
	if errExec != nil {
		log.Exitf(1, "failed to execute command: %s", errExec.Error())
	}

	log.Success("Migration force completed")
}
