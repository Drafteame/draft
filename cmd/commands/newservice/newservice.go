package newservice

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/actions/dtos"
	"github.com/Drafteame/draft/internal/actions/newservice"
	"github.com/Drafteame/draft/internal/flags"
	"github.com/Drafteame/draft/internal/forms"
)

var newServiceCmd = &cobra.Command{
	Use:   "new:service",
	Short: "Create a new service",
	Long:  "Create a new service",
	Run:   run,
}

func run(_ *cobra.Command, _ []string) {
	if flags.Flags.WorkingDir != "" {
		if err := os.Chdir(flags.Flags.WorkingDir); err != nil {
			panic(err)
		}
	}

	input := dtos.Input{}

	if err := forms.NewService(&input); err != nil {
		panic(err)
	}

	action := newservice.GetAction(input)

	if err := action.Exec(); err != nil {
		panic(err)
	}

	if err := action.PostCreate(); err != nil {
		panic(err)
	}

	println("Service created successfully")
}

func GetCmd() *cobra.Command {
	return newServiceCmd
}
