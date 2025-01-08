package deleteproject

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/actions/sentry/deleteproject"
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/forms"
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

	input := dtos.DeleteProjectInput{}

	if err := forms.DeleteProject(&input); err != nil {
		if err.Error() == "operation cancelled by the user" {
			fmt.Println("The operation was cancelled by the user.")
			return
		}
		panic(err)
	}

	action, err := deleteproject.GetAction(input)
	if err != nil {
		panic(err)
	}

	if err = action.Exec(); err != nil {
		panic(err)
	}
}

func GetCmd() *cobra.Command {
	return cmd
}
