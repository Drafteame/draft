package newservice

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/actions/newservice"
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/forms"
)

var cmdNewService = &cobra.Command{
	Use:   "new:service",
	Short: "Create a new service",
	Long:  "Create a new service",
	Run:   run,
}

func init() {
	cmdNewService.Flags().BoolVarP(&data.Flags.NoSentry, "no-sentry", "", data.Flags.NoSentry, "Disable sentry project creation")
}

func run(cmd *cobra.Command, _ []string) {
	if data.Flags.WorkingDir != "" {
		if err := os.Chdir(data.Flags.WorkingDir); err != nil {
			panic(err)
		}
	}

	data.LoadMeta()

	input := dtos.ServiceInput{}

	useDig, err := cmd.Parent().Flags().GetBool("use-dig")
	if err != nil {
		panic(err)
	}
	input.UseDig = useDig

	legacyPath, err := cmd.Parent().Flags().GetString("legacy-path")
	if err != nil {
		panic(err)
	}

	if legacyPath != "" {
		input.IsLegacy = true
		input.ServicePath = legacyPath
	}

	if err := forms.NewService(&input); err != nil {
		panic(err)
	}

	action, err := newservice.GetAction(input)
	if err != nil {
		panic(err)
	}

	if errExec := action.Exec(); errExec != nil {
		panic(errExec)
	}

	println("Service created successfully")
}

func GetCmd() *cobra.Command {
	return cmdNewService
}
