package newlambda

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/internal/actions/newlambda"
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/forms"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var newLambdaCmd = &cobra.Command{
	Use:   "new:lambda",
	Short: "Create a new service",
	Long:  "Create a new service",
	Run:   run,
}

func init() {
	newLambdaCmd.Flags().Bool("use-dig", false, "Use uber dig for dependency injection")
	newLambdaCmd.Flags().StringP("legacy-path", "l", "", "Path to legacy service")
	newLambdaCmd.Flags().BoolP("with-frame", "f", false, "Use framev2 for lambda creation")
}

func run(cmd *cobra.Command, _ []string) {
	common.ChDir(cmd)

	data.LoadMeta()

	useDig, err := cmd.Flags().GetBool("use-dig")
	if err != nil {
		log.Exitf(1, "failed to obtain use-dig flag: %s", err.Error())
	}
	input := dtos.LambdaInput{
		UseDig: useDig,
	}

	legacyPath, err := cmd.Flags().GetString("legacy-path")
	if err != nil {
		log.Exitf(1, "failed to obtain legacy-path flag: %s", err.Error())
	}

	if legacyPath != "" {
		input.IsLegacy = true
		input.ServicePath = legacyPath
	}

	withFrame, err := cmd.Flags().GetBool("with-frame")
	if err != nil {
		log.Exitf(1, "failed to obtain with-frame flag: %s", err.Error())
	}
	input.WithFrame = withFrame

	if errForm := forms.NewLambda(&input); errForm != nil {
		log.Exitf(1, "failed to collect new lambda info: %s", errForm.Error())
	}

	if errExec := newlambda.New(input).Exec(); errExec != nil {
		log.Exitf(1, "failed to create lambda: %s", errExec.Error())
	}

	log.Success("Lambda created successfully")
}

func GetCmd() *cobra.Command {
	return newLambdaCmd
}
