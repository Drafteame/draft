package main

import (
	"github.com/Drafteame/draft/cmd/commands"
	"github.com/Drafteame/draft/cmd/commands/config"
	"github.com/Drafteame/draft/cmd/commands/local/invoke"
	migrateforce "github.com/Drafteame/draft/cmd/commands/local/migrate/force"
	migrateup "github.com/Drafteame/draft/cmd/commands/local/migrate/up"
	testsetup "github.com/Drafteame/draft/cmd/commands/local/setup"
	"github.com/Drafteame/draft/cmd/commands/newdomain"
	"github.com/Drafteame/draft/cmd/commands/newlambda"
	"github.com/Drafteame/draft/cmd/commands/newservice"
	"github.com/Drafteame/draft/cmd/commands/sentry/deleteproject"
)

func main() {
	cmd := commands.GetCmd()

	cmd.AddCommand(config.GetCmd())
	cmd.AddCommand(newservice.GetCmd())
	cmd.AddCommand(newlambda.GetCmd())
	cmd.AddCommand(newdomain.GetCmd())
	cmd.AddCommand(deleteproject.GetCmd())
	cmd.AddCommand(invoke.GetCmd())
	cmd.AddCommand(migrateup.GetCmd())
	cmd.AddCommand(migrateforce.GetCmd())
	cmd.AddCommand(testsetup.GetCmd())

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
