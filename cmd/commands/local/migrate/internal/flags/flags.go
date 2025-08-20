package flags

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/pkg/log"
)

type Flags struct {
	Database string
	Config   string
	All      bool
	Group    string
}

// Register sets up flags for the given cobra command to manage database migrations and configuration.
func Register(cmd *cobra.Command) {
	cmd.Flags().StringP("database", "D", "", "database name")
	cmd.Flags().StringP("local-migrate-config", "c", ".local-migrate-config.yml", "path to the migrations config file")
	cmd.Flags().Bool("all", false, "migrate all databases")
	cmd.Flags().String("group", "", "DB migrations group name")
}

// GetFlags extracts and returns command-line flags from the provided cobra.Command as a Flags struct.
func GetFlags(cmd *cobra.Command) *Flags {
	database, err := cmd.Flags().GetString("database")
	if err != nil {
		log.Exitf(1, "invalid database: %s", err.Error())
	}

	localMigrateConfig, err := cmd.Flags().GetString("local-migrate-config")
	if err != nil {
		log.Exitf(1, "invalid local migrate config: %s", err.Error())
	}

	group, err := cmd.Flags().GetString("group")
	if err != nil {
		log.Exitf(1, "invalid group flag: %s", err.Error())
	}

	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		log.Exitf(1, "invalid all flag: %s", err.Error())
	}

	return &Flags{
		Database: database,
		Config:   localMigrateConfig,
		All:      all,
		Group:    group,
	}
}
