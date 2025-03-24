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

var legacyPath string

func init() {
	newLambdaCmd.Flags().StringVarP(&legacyPath, "legacy-path", "l", "", "Path to legacy service")
	newLambdaCmd.Flags().Bool("use-dig", false, "Use uber dig for dependency injection")
}

func run(cmd *cobra.Command, _ []string) {
	if data.Flags.WorkingDir != "" {
		if err := os.Chdir(data.Flags.WorkingDir); err != nil {
			panic(err)
		}
	}

	data.LoadMeta()

	useDig, err := cmd.Flags().GetBool("use-dig")
	if err != nil {
		panic(err)
	}
	input := dtos.LambdaInput{
		UseDig: useDig,
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
