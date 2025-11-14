package mockery

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/internal/actions/mockery"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var serviceName string
var servicePath string

var mockeryCmd = &cobra.Command{
	Use:   "mockery",
	Short: "Run mockery for all services",
	Long: `Run mockery for specific service with .mockery.yml configuration file.

This command will:
1. Look for a .mockery.yml file in the service directory
2. Execute mockery with the specific configuration file

Examples:
  # Run mockery for all services in services directory
  draft mockery

  # Run mockery for the specific service
  draft mockery -s <servicepkg>

  # Run mockery using an absolute path from the root
  draft mockery -p <servicepath/servicepkg>`,
	Run: run,
}

func init() {
	mockeryCmd.Flags().StringVarP(&serviceName, "service", "s", "", "Service name (will look in services/<servicepkg>)")
	mockeryCmd.Flags().StringVarP(&servicePath, "path", "p", "", "Absolute path from root")
}

func run(cmd *cobra.Command, _ []string) {
	common.ChDir(cmd)

	if err := mockery.New(serviceName, servicePath).Exec(); err != nil {
		log.Exitf(1, "Failed to run mockery: %v", err)
	}

	log.Success("Mockery executed successfully")
}

func GetCmd() *cobra.Command {
	return mockeryCmd
}
