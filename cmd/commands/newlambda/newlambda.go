package newlambda

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/actions/dtos"
	"github.com/Drafteame/draft/internal/actions/newlambda"
	"github.com/Drafteame/draft/internal/flags"
	"github.com/Drafteame/draft/internal/forms"
)

var newLambdaCmd = &cobra.Command{
	Use:   "new:lambda",
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

	if err := forms.NewLambda(&input); err != nil {
		panic(err)
	}

	if err := newlambda.GetAction().Exec(input); err != nil {
		panic(err)
	}
}

func GetCmd() *cobra.Command {
	return newLambdaCmd
}
