package newservice

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/internal/actions/newservice"
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/forms"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var cmdNewService = &cobra.Command{
	Use:   "new:service",
	Short: "Create a new service",
	Long:  "Create a new service",
	Run:   run,
}

func init() {
	cmdNewService.Flags().BoolVarP(&data.Flags.NoSentry, "no-sentry", "", data.Flags.NoSentry, "Disable sentry project creation")
	cmdNewService.Flags().Bool("use-dig", false, "Use uber dig for dependency injection")
	cmdNewService.Flags().StringP("legacy-path", "l", "", "Path to legacy service")
}

func run(cmd *cobra.Command, _ []string) {
	common.ChDir(cmd)

	data.LoadMeta()

	input := dtos.ServiceInput{}

	useDig, err := cmd.Flags().GetBool("use-dig")
	if err != nil {
		log.Exitf(1, "failed to obtain use-dig flag: %s", err.Error())
	}
	input.UseDig = useDig

	legacyPath, err := cmd.Flags().GetString("legacy-path")
	if err != nil {
		log.Exitf(1, "failed to obtain legacy-path flag: %s", err.Error())
	}

	if legacyPath != "" {
		input.IsLegacy = true
		input.ServicePath = legacyPath
	}

	if errForm := forms.NewService(&input); errForm != nil {
		log.Exitf(1, "failed to collect new service info: %s", errForm.Error())
	}

	if errExec := newservice.New(input).Exec(); errExec != nil {
		log.Exitf(1, "failed to create service: %s", errExec.Error())
	}

	log.Success("Service created successfully")
}

func GetCmd() *cobra.Command {
	return cmdNewService
}
