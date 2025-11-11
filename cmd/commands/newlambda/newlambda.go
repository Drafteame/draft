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
	Short: "Create a new lambda function",
	Long:  "Create a new lambda function with the specified configuration and trigger type",
	Run:   run,
}

func init() {
	newLambdaCmd.Flags().Bool("use-dig", false, "Use uber dig for dependency injection")
	newLambdaCmd.Flags().StringP("legacy-path", "l", "", "Path to legacy service")
}

func run(cmd *cobra.Command, _ []string) {
	common.ChDir(cmd)

	data.LoadMeta()

	useDig := common.GetBoolFlag(cmd, "use-dig")
	input := dtos.LambdaInput{
		UseDig: useDig,
	}

	legacyPath := common.GetStringFlag(cmd, "legacy-path")

	if legacyPath != "" {
		input.IsLegacy = true
		input.ServicePath = legacyPath
	}

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
