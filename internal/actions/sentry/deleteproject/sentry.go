package deleteproject

import (
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/pkg/log"
	"github.com/Drafteame/draft/internal/pkg/sentry"
)

func (dl *DeleteProject) sentry() error {
	if data.Flags.NoSentry {
		return nil
	}

	return dl.deleteSentryProject()
}

func (dl *DeleteProject) deleteSentryProject() error {
	if err := sentry.DeleteProject(dl.input.ProjectID); err != nil {
		return err
	}

	log.Info("Sentry project with ID has been removed:", dl.input.ProjectID)

	return nil
}
