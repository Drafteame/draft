package newservice

import (
	"github.com/Drafteame/draft/internal/config"
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/pkg/log"
	"github.com/Drafteame/draft/internal/pkg/sentry"
)

func (ns *NewService) sentry() error {
	if data.Flags.NoSentry {
		return nil
	}

	if err := ns.createSentryProject(); err != nil {
		return err
	}

	return ns.setupSentryStages()
}

func (ns *NewService) createSentryProject() error {
	cfg, errCfg := config.Get()
	if errCfg != nil {
		return errCfg
	}

	if cfg.Sentry.Token == "" {
		log.Warn("Sentry token not found, skipping project creation")
		return nil
	}

	projectID, err := sentry.CreateProject(ns.input.ServiceName)
	if err != nil {
		return err
	}

	log.Info("Sentry project created with ID", projectID)

	keys, err := sentry.GetClientKeys(projectID)
	if err != nil {
		return err
	}

	ns.input.HasSentry = true
	ns.input.SentryDSN = keys["dsn"]

	log.Info("Sentry DSN:", ns.input.SentryDSN)

	return nil
}

func (ns *NewService) setupSentryStages() error {
	if !ns.input.HasSentry {
		return nil
	}

	dsn := ns.input.SentryDSN
	serviceName := ns.input.ServiceName

	return sentry.CreateStages(serviceName, dsn)
}
