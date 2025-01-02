package deleteproject

import (
	"github.com/Drafteame/draft/internal/actions/sentry/deleteproject"
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/forms"
	"github.com/spf13/cobra"
	"os"
)

var cmd = &cobra.Command{
	Use:   "sentry:delete-project",
	Short: "Delete a Sentry project",
	Long:  "Delete a Sentry project",
	Run:   run,
}

func run(_ *cobra.Command, _ []string) {
	if data.Flags.WorkingDir != "" {
		if err := os.Chdir(data.Flags.WorkingDir); err != nil {
			panic(err)
		}
	}

	data.LoadMeta() // carga el nombre del paquete

	input := dtos.DeleteProjectInput{}

	//form
	if err := forms.DeleteProject(&input); err != nil {
		panic(err)
	}

	//Action
	action, err := deleteproject.GetAction(input)
	if err != nil {
		panic(err)
	}

	if err := action.Exec(); err != nil {
		panic(err)
	}

	println("Sentry project deleted successfully")

}

func GetCmd() *cobra.Command {
	return cmd
}
