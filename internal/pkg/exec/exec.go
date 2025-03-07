package exec

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type commandOptions struct {
	envs   map[string]string
	stdout io.Writer
	stderr io.Writer
}

type CommandOption func(*commandOptions)

func WithEnvs(envs map[string]string) CommandOption {
	return func(o *commandOptions) {
		o.envs = envs
	}
}

func WithStdout(w io.Writer) CommandOption {
	return func(o *commandOptions) {
		o.stdout = w
	}
}

func WithStderr(w io.Writer) CommandOption {
	return func(o *commandOptions) {
		o.stderr = w
	}
}

func Command(script string, opts ...CommandOption) (string, error) {
	options := commandOptions{}

	for _, opt := range opts {
		opt(&options)
	}

	// Create the command, using "bash" with the "-c" flag to execute cmdStr
	cmd := exec.Command("bash", "-c", script)

	// Buffers to capture stdout and stderr
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	outWriters := []io.Writer{&stdout}
	errWriters := []io.Writer{&stderr}

	if options.stdout != nil {
		outWriters = append(outWriters, options.stdout)
	}

	if options.stderr != nil {
		errWriters = append(errWriters, options.stderr)
	}

	cmd.Stdout = io.MultiWriter(outWriters...)
	cmd.Stderr = io.MultiWriter(errWriters...)

	// Set the environment variables
	cmd.Env = []string{}

	for key, value := range options.envs {
		cmd.Env = append(cmd.Env, key+"="+value)
	}

	cmd.Env = append(cmd.Env, os.Environ()...)

	// Run the command
	err := cmd.Run()

	// Get the exit code
	exitCode := 0
	var exitErr *exec.ExitError

	if errors.As(err, &exitErr) {
		exitCode = exitErr.ExitCode()
	}

	out := strings.TrimSpace(stdout.String() + stderr.String())

	if err != nil {
		return out, fmt.Errorf("exec: %s [code %d]: %s, %w", script, exitCode, out, err)
	}

	// Return the combined output, exit code, and error
	return out, err
}
