package newlambda

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/actions/newlambda"
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/forms"
)

var newLambdaCmd = &cobra.Command{
	Use:   "new:lambda",
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

	input := dtos.ServiceInput{}

	if err := forms.NewLambda(&input); err != nil {
		panic(err)
	}

	action := newlambda.GetAction(input)

	if err := action.Exec(); err != nil {
		panic(err)
	}

	println("Lambda created successfully")
}

func GetCmd() *cobra.Command {
	return newLambdaCmd
}
