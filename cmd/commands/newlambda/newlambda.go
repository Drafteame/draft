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

func init() {
	// TODO: add flags for lambda creation
}

func run(cmd *cobra.Command, _ []string) {
	if data.Flags.WorkingDir != "" {
		if err := os.Chdir(data.Flags.WorkingDir); err != nil {
			panic(err)
		}
	}

	data.LoadMeta()

	useDig, err := cmd.Parent().Flags().GetBool("use-dig")
	if err != nil {
		panic(err)
	}
	input := dtos.LambdaInput{
		UseDig: useDig,
	}

	legacyPath, err := cmd.Parent().Flags().GetString("legacy-path")
	if err != nil {
		panic(err)
	}

	if legacyPath != "" {
		input.IsLegacy = true
		input.ServicePath = legacyPath
	}

	if err := forms.NewLambda(&input); err != nil {
		panic(err)
	}

	action, err := newlambda.GetAction(input)
	if err != nil {
		panic(err)
	}

	if err := action.Exec(); err != nil {
		panic(err)
	}

	println("Lambda created successfully")
}

func GetCmd() *cobra.Command {
	return newLambdaCmd
}
