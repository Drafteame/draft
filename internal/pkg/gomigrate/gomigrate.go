package gomigrate

import (
	"fmt"
	"os"
	"strings"

	"github.com/Drafteame/draft/internal/pkg/exec"
)

type Action string

const (
	ActionUp    Action = "up"
	ActionDown  Action = "down"
	ActionForce Action = "force"
)

type Config struct {
	Source   string
	Database string
	Action   string
	Args     []string
}

func Exec(action Action, config Config) error {
	args := []string{"migrate"}

	if config.Source != "" {
		args = append(args, "-source", fmt.Sprintf("file://%s", config.Source))
	}

	if config.Database != "" {
		args = append(args, "-database", config.Database)
	}

	args = append(args, string(action))

	if len(config.Args) > 0 {
		args = append(args, config.Args...)
	}

	script := strings.Join(args, " ")

	fmt.Printf("Running command: %s", script)

	_, errExec := exec.Command(script, exec.WithStdout(os.Stdout), exec.WithStderr(os.Stderr))
	return errExec
}
