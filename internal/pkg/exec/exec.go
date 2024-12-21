package exec

import (
	"bytes"
	"errors"
	"os/exec"
)

func Command(script string) (string, int, error) {
	// Create the command, using "bash" with the "-c" flag to execute cmdStr
	cmd := exec.Command("bash", "-c", script)

	// Buffers to capture stdout and stderr
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()

	// Get the exit code
	exitCode := 0
	var exitErr *exec.ExitError

	if errors.As(err, &exitErr) {
		exitCode = exitErr.ExitCode()
	}

	// Return the combined output, exit code, and error
	return stdout.String() + stderr.String(), exitCode, err
}
