package newservice

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/actions/dtos"
	"github.com/Drafteame/draft/internal/actions/newservice"
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/forms"
)

var newServiceCmd = &cobra.Command{
	Use:   "new:service",
	Short: "Create a new service",
	Long:  "Create a new service",
	Run:   run,
}

func run(_ *cobra.Command, _ []string) {
	if data.Flags.WorkingDir != "" {
		if err := os.Chdir(data.Flags.WorkingDir); err != nil {
			panic(err)
		}
	}

	data.LoadMeta()

	input := dtos.Input{}

	if err := forms.NewService(&input); err != nil {
		panic(err)
	}

	action := newservice.GetAction(input)

	if err := action.Exec(); err != nil {
		panic(err)
	}

	println("Service created successfully")
}

func GetCmd() *cobra.Command {
	return newServiceCmd
}
