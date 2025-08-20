package setup

import (
	"fmt"
	"os"
	"time"

	migrateup "github.com/Drafteame/draft/internal/actions/local/migrate/command"
	"github.com/Drafteame/draft/internal/pkg/exec"
	"github.com/Drafteame/draft/internal/pkg/log"
)

func (a *Action) init() error {
	if !a.Input.BypassDocker {
		cmd := "docker compose -f ./docker/docker-compose.yml up -d"

		_, err := exec.Command(cmd, exec.WithStdout(os.Stdout), exec.WithStderr(os.Stderr))
		if err != nil {
			return fmt.Errorf("init failed: %w", err)
		}

		if errHealth := a.postgresHealthCheck(); errHealth != nil {
			return errHealth
		}
	}

	input := migrateup.Input{
		Command:            "up",
		LocalMigrateConfig: ".local-migrate-config.yml",
		All:                true,
	}

	errMigrate := migrateup.New(input).Exec()
	if errMigrate != nil {
		return fmt.Errorf("migrate up failed: %w", errMigrate)
	}

	input.Group = "test"

	errMigrateTest := migrateup.New(input).Exec()
	if errMigrateTest != nil {
		return fmt.Errorf("test migrate up failed: %w", errMigrateTest)
	}

	return nil
}

func (a *Action) postgresHealthCheck() error {
	log.Debug("postgres health check")

	cmd := "docker exec -i nix-local-postgres pg_isready -U root"

	maxTries := 5

	var errTries error

	for i := 0; i < maxTries; i++ {
		_, err := exec.Command(cmd, exec.WithStdout(os.Stdout), exec.WithStderr(os.Stderr))
		if err == nil {
			time.Sleep(2 * time.Second)
			errTries = nil
			break
		}

		log.Debugf("retrying postgres health check (%d/%d): %v\n", i+1, maxTries, err)

		sleep := time.Duration(2 * (i + 1))
		time.Sleep(sleep * time.Second)

		errTries = err
	}

	return errTries
}
