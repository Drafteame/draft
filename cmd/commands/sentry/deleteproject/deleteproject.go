package deleteproject

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/internal/actions/sentry/deleteproject"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/forms"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var sentryDeleteProjectCmd = &cobra.Command{
	Use:   "sentry:delete-project",
	Short: "Delete a Sentry project",
	Long:  "Delete a Sentry project",
	Run:   run,
}

func run(cmd *cobra.Command, _ []string) {
	common.ChDir(cmd)

	input := dtos.DeleteProjectInput{}

	if err := forms.DeleteProject(&input); err != nil {
		if err.Error() == "operation cancelled by the user" {
			log.Warn("The operation was cancelled by the user.")
			return
		}

		log.Exitf(1, "Failed to collect project info: %v", err)
	}

	if err := deleteproject.New(input).Exec(); err != nil {
		log.Exitf(1, "Failed to delete sentry project: %v", err)
	}
}

func GetCmd() *cobra.Command {
	return sentryDeleteProjectCmd
}
