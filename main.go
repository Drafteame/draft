package main

import (
	"github.com/Drafteame/draft/cmd/commands"
	"github.com/Drafteame/draft/cmd/commands/newlambda"
	"github.com/Drafteame/draft/cmd/commands/newservice"
)

func main() {
	cmd := commands.GetCmd()

	cmd.AddCommand(newservice.GetCmd())
	cmd.AddCommand(newlambda.GetCmd())

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
