package newservice

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh/spinner"

	"github.com/Drafteame/draft/internal/pkg/exec"
)

// fixOtel runs the monorepo's `mage otel:fix` to sync the newly created
// service's otel-layer collector config and environment.yml otel section
// against the repo's single source of truth. Must run before updateSlsPlugins
// changes the working directory away from the repo root.
func (ns *NewService) fixOtel() error {
	var err error

	action := func() {
		_, errCmd := exec.Command("mage otel:fix")
		if errCmd != nil {
			err = fmt.Errorf("cannot sync otel configuration: %w", errCmd)
			return
		}
	}

	spinErr := spinner.New().
		Title("Syncing OTel configuration...").
		Action(action).
		Run()

	return errors.Join(spinErr, err)
}
