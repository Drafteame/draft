package newservice

import (
	"github.com/Drafteame/draft/internal/config"
	"github.com/Drafteame/draft/internal/sentry"
)

func (css *NewService) preCreate() error {
	return css.createSentryProject()
}

func (css *NewService) createSentryProject() error {
	cfg := config.Get()

	if cfg.Sentry.Token == "" {
		println("Sentry token not found, skipping project creation")
		return nil
	}

	projectID, err := sentry.CreateProject(css.input.ServiceName)
	if err != nil {
		return err
	}

	println("Sentry project created with ID", projectID)

	keys, err := sentry.GetClientKeys(projectID)
	if err != nil {
		return err
	}

	css.input.HasSentry = true
	css.input.SentryDSN = keys["dsn"]

	println("Sentry DSN:", css.input.SentryDSN)

	return nil
}
