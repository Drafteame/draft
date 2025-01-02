package deleteproject

import (
	"github.com/Drafteame/draft/internal/config"
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/pkg/sentry"
)

func (dl *DeleteProject) sentry() error {
	if data.Flags.NoSentry {
		return nil
	}

	if dl.input.ProjectID == "" {
		println("Process canceled")
		return nil
	}

	return dl.deleteSentryProject()
}

func (dl *DeleteProject) deleteSentryProject() error {
	cfg := config.Get()

	if cfg.Sentry.Token == "" {
		println("Sentry token not found, skipping project deletion")
		return nil
	}

	if err := sentry.DeleteProject(dl.input.ProjectID); err != nil {
		return err
	}

	println("Sentry project with ID has been removed:", dl.input.ProjectID)

	return nil
}
