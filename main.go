package main

import (
	"github.com/Drafteame/draft/cmd/commands"
	"github.com/Drafteame/draft/cmd/commands/config"
	"github.com/Drafteame/draft/cmd/commands/local/invoke"
	migratedown "github.com/Drafteame/draft/cmd/commands/local/migrate/down"
	migrateforce "github.com/Drafteame/draft/cmd/commands/local/migrate/force"
	migrateup "github.com/Drafteame/draft/cmd/commands/local/migrate/up"
	testsetup "github.com/Drafteame/draft/cmd/commands/local/setup"
	"github.com/Drafteame/draft/cmd/commands/mock/generate"
	"github.com/Drafteame/draft/cmd/commands/newdomain"
	"github.com/Drafteame/draft/cmd/commands/newlambda"
	"github.com/Drafteame/draft/cmd/commands/newservice"
	"github.com/Drafteame/draft/cmd/commands/sentry/project/create"
	"github.com/Drafteame/draft/cmd/commands/sentry/project/delete"
	"github.com/Drafteame/draft/internal/pkg/log"
)

func main() {
	cmd := commands.GetCmd()

	cmd.AddCommand(config.GetCmd())
	cmd.AddCommand(newservice.GetCmd())
	cmd.AddCommand(newlambda.GetCmd())
	cmd.AddCommand(newdomain.GetCmd())
	cmd.AddCommand(create.GetCmd())
	cmd.AddCommand(delete.GetCmd())
	cmd.AddCommand(invoke.GetCmd())
	cmd.AddCommand(migrateup.GetCmd())
	cmd.AddCommand(migrateforce.GetCmd())
	cmd.AddCommand(migratedown.GetCmd())
	cmd.AddCommand(testsetup.GetCmd())
	cmd.AddCommand(generate.GetCmd())

	if err := cmd.Execute(); err != nil {
		log.Exit(1, err.Error())
	}
}
