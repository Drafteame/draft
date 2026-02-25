package mockery

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/internal/actions/mockery"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var (
	dry    bool
	gitMod bool
)

var mockeryCmd = &cobra.Command{
	Use:   "mockery [config-files...]",
	Short: "Run mockery with base and package-specific configurations",
	Long: `Run mockery by merging all .mockery.pkg.yml files into a single invocation.

This command will:
1. Load base configuration from .mockery.base.yml
2. Search for all .mockery.pkg.yml files (or use provided list)
3. Merge all package configs into a single unified config
4. Execute mockery once for all packages
5. Clean up temporary files and exit with code 1 if failed

Examples:
  # Run mockery for all .mockery.pkg.yml files found in project
  draft mockery

  # Run mockery for specific package config files
  draft mockery services/user/.mockery.pkg.yml services/auth/.mockery.pkg.yml

  # Dry run - validate configs without executing mockery
  draft mockery --dry

  # Run mockery only for packages with modified files (compare HEAD with main)
  draft mockery --git-mod`,
	Run: run,
}

func init() {
	mockeryCmd.Flags().BoolVar(&dry, "dry", false, "Dry run - validate and prepare configs without executing mockery")
	mockeryCmd.Flags().BoolVar(&gitMod, "git-mod", false, "Only run mockery for packages with modified files (compares HEAD with main branch)")
}

func run(cmd *cobra.Command, args []string) {
	common.ChDir(cmd)

	if err := mockery.New(cmd.Context(), args, dry, gitMod).Exec(); err != nil {
		log.Exitf(1, "Failed to run mockery: %v", err)
	}

	if dry {
		log.Success("Mockery dry run completed successfully")
	} else {
		log.Success("Mockery executed successfully")
	}
}

func GetCmd() *cobra.Command {
	return mockeryCmd
}
