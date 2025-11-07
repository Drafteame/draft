package createproject

import (
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/pkg/log"
	"github.com/Drafteame/draft/internal/pkg/sentry"
)

func (cp *CreateProject) sentry() error {
	if data.Flags.NoSentry {
		return nil
	}

	return cp.createSentryProject()
}

func (cp *CreateProject) createSentryProject() error {
	log.Info("Creating Sentry project:", cp.input.ProjectName)

	projectID, err := sentry.CreateProject(cp.input.ProjectName)
	if err != nil {
		return err
	}

	log.Info("Sentry project created with ID:", projectID)

	keys, err := sentry.GetClientKeys(projectID)
	if err != nil {
		return err
	}

	log.Info("Sentry DSN:", keys["dsn"])

	log.Info("Creating dev and prod stages...")

	if err := sentry.CreateStages(cp.input.ProjectName, keys["dsn"]); err != nil {
		return err
	}

	log.Info("Stages created successfully!")
	log.Info("Project setup complete!")

	return nil
}
