package connect

import (
	"github.com/spf13/cobra"

	dblist "github.com/Drafteame/draft/cmd/commands/db/connect/list"
	dbstart "github.com/Drafteame/draft/cmd/commands/db/connect/start"
	dbstatus "github.com/Drafteame/draft/cmd/commands/db/connect/status"
	dbstop "github.com/Drafteame/draft/cmd/commands/db/connect/stop"
)

var cmd = &cobra.Command{
	Use:   "db:connect",
	Short: "Manage SSM port-forwarding tunnels to remote databases",
}

func GetCmd() *cobra.Command {
	cmd.AddCommand(dblist.GetCmd())
	cmd.AddCommand(dbstatus.GetCmd())
	cmd.AddCommand(dbstart.GetCmd())
	cmd.AddCommand(dbstop.GetCmd())

	return cmd
}
