package generate

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/internal/actions/mock/generate"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var mockGenerateCmd = &cobra.Command{
	Use:   "mock:generate",
	Short: "Generate mocks using mockery with concatenated configuration",
	Long: `Generates mocks by combining .mockery.base.yml with all .mockery-packages.yml files found in the project.

This command:
  1. Reads the base mockery configuration from .mockery.base.yml
  2. Searches for all .mockery-packages.yml files in the project (using fd)
  3. Concatenates them into a temporary .mockery.yml file
  4. Runs mockery to generate the mocks
  5. Cleans up the temporary .mockery.yml file

Example:
  draft mock:generate
  draft mock:generate -w /path/to/project
`,
	Run: run,
}

func run(cmd *cobra.Command, _ []string) {
	// Change to working directory if specified
	common.ChDir(cmd)

	// Create input DTO
	// Use "." since ChDir has already moved us to the correct directory
	input := dtos.MockGenerateInput{
		WorkingDir: ".",
	}

	// Execute action
	if err := generate.New(input).Exec(); err != nil {
		log.Exitf(1, "failed to generate mocks: %s", err.Error())
	}

	log.Success("🎉 Mock generation completed successfully")
}

func GetCmd() *cobra.Command {
	return mockGenerateCmd
}
