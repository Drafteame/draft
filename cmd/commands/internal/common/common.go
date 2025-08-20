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
		log.Exitf(1, "Failed to change working directory: %v", err)
	}
}
