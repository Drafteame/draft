package main

import (
	"github.com/Drafteame/draft/cmd/commands"
	"github.com/Drafteame/draft/cmd/commands/config"
	"github.com/Drafteame/draft/cmd/commands/local/invoke"
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

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
