package setup

import (
	"fmt"
	"os"

	"github.com/Drafteame/draft/internal/pkg/exec"
)

func (a *Action) prune() error {
	println("pruning env")

	if !a.Input.Prune || a.Input.BypassDocker {
		return nil
	}

	commands := []string{
		"docker compose -f ./docker/docker-compose.yml down",
		"rm -rf ~/.draftea/docker/data/postgres",
	}

	for _, cmd := range commands {
		println("cmd:", cmd)
		_, err := exec.Command(cmd, exec.WithStdout(os.Stdout), exec.WithStderr(os.Stderr))
		if err != nil {
			return fmt.Errorf("prune failed: %w", err)
		}
	}

	return nil
}
