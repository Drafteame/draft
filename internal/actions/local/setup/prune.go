package setup

import (
	"fmt"
	"os"

	"github.com/Drafteame/draft/internal/pkg/exec"
	"github.com/Drafteame/draft/internal/pkg/log"
)

func (a *Action) prune() error {
	log.Debug("Pruning env")

	if !a.Input.Prune || a.Input.BypassDocker {
		return nil
	}

	commands := []string{
		"docker compose -f ./docker/docker-compose.yml down",
		"rm -rf ~/.draftea/docker/data/postgres",
	}

	for _, cmd := range commands {
		log.Debug("Prune command:", cmd)
		_, err := exec.Command(cmd, exec.WithStdout(os.Stdout), exec.WithStderr(os.Stderr))
		if err != nil {
			return fmt.Errorf("prune failed: %w", err)
		}
	}

	return nil
}
