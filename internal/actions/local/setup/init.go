package setup

import (
	"fmt"
	"os"
	"time"

	migrateup "github.com/Drafteame/draft/internal/actions/local/migrate/command"
	"github.com/Drafteame/draft/internal/pkg/exec"
)

func (a *Action) init() error {
	if !a.Input.BypassDocker {
		cmd := "docker compose -f ./docker/docker-compose.yml up -d"

		_, err := exec.Command(cmd, exec.WithStdout(os.Stdout), exec.WithStderr(os.Stderr))
		if err != nil {
			return fmt.Errorf("init failed: %w", err)
		}

		if err := a.postgresHealthCheck(); err != nil {
			return err
		}
	}

	input := migrateup.Input{
		Command:            "up",
		WorkingDir:         a.Input.WorkingDir,
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
	println("postgres health check")

	cmd := "docker exec -i nix-local-postgres pg_isready -U root"

	maxTries := 5

	var errtries error

	for i := 0; i < maxTries; i++ {
		_, err := exec.Command(cmd, exec.WithStdout(os.Stdout), exec.WithStderr(os.Stderr))
		if err == nil {
			time.Sleep(2 * time.Second)
			errtries = nil
			break
		}

		_, _ = fmt.Printf("retrying (%d/%d): %v\n", i+1, maxTries, err)

		sleep := time.Duration(2 * (i + 1))
		time.Sleep(sleep * time.Second)

		errtries = err
	}

	return errtries
}
