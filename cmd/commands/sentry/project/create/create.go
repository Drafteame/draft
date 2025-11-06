package create

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/internal/actions/sentry/createproject"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/forms"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var sentryCreateProjectCmd = &cobra.Command{
	Use:   "sentry:project:create",
	Short: "Create a new Sentry project",
	Long:  "Create a new Sentry project and configure initial stages (dev and prod)",
	Run:   run,
}

func run(cmd *cobra.Command, _ []string) {
	common.ChDir(cmd)

	input := dtos.CreateProjectInput{}

	if err := forms.CreateProject(&input); err != nil {
		log.Exitf(1, "Failed to collect project info: %v", err)
	}

	if err := createproject.New(input).Exec(); err != nil {
		log.Exitf(1, "Failed to create sentry project: %v", err)
	}
}

func GetCmd() *cobra.Command {
	return sentryCreateProjectCmd
}
