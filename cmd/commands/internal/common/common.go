package common

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/pkg/log"
)

// ChDir changes the current working directory to the value of the "working-dir" flag of the parent command, if set.
// Logs an error and exits if the directory change fails.
func ChDir(cmd *cobra.Command) {
	workDir := cmd.Parent().Flag("working-dir").Value.String()
	if workDir == "" {
		return
	}

	if err := os.Chdir(workDir); err != nil {
		log.Exitf(1, "failed to change working directory: %v", err)
	}
}

// GetBoolFlag retrieves a boolean flag value from a cobra command.
// Logs an error and exits if the flag cannot be retrieved.
func GetBoolFlag(cmd *cobra.Command, name string) bool {
	val, err := cmd.Flags().GetBool(name)
	if err != nil {
		log.Exitf(1, "failed to obtain %s flag: %s", name, err.Error())
	}
	return val
}

// GetStringFlag retrieves a string flag value from a cobra command.
// Logs an error and exits if the flag cannot be retrieved.
func GetStringFlag(cmd *cobra.Command, name string) string {
	val, err := cmd.Flags().GetString(name)
	if err != nil {
		log.Exitf(1, "failed to obtain %s flag: %s", name, err.Error())
	}
	return val
}
