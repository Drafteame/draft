package mockery

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/internal/actions/mockery"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var jobsNum int

var mockeryCmd = &cobra.Command{
	Use:   "mockery [config-files...]",
	Short: "Run mockery with base and package-specific configurations",
	Long: `Run mockery by merging .mockery.base.yml with .mockery.pkg.yml files.

This command will:
1. Load base configuration from .mockery.base.yml
2. Search for all .mockery.pkg.yml files (or use provided list)
3. Merge base config with package config (package takes precedence)
4. Generate temporary config files and execute mockery concurrently
5. Report failures without stopping execution
6. Clean up temporary files and exit with code 1 if any failed

Examples:
  # Run mockery for all .mockery.pkg.yml files found in project
  draft mockery

  # Run mockery for specific package config files
  draft mockery services/user/.mockery.pkg.yml services/auth/.mockery.pkg.yml

  # Run with custom concurrent job limit (default: 3)
  draft mockery --jobs-num 5`,
	Run: run,
}

func init() {
	mockeryCmd.Flags().IntVarP(&jobsNum, "jobs-num", "j", 5, "Number of concurrent mockery jobs to run")
}

func run(cmd *cobra.Command, args []string) {
	common.ChDir(cmd)

	if err := mockery.New(args, jobsNum).Exec(); err != nil {
		log.Exitf(1, "Failed to run mockery: %v", err)
	}

	log.Success("Mockery executed successfully")
}

func GetCmd() *cobra.Command {
	return mockeryCmd
}
