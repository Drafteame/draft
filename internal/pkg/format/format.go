package format

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh/spinner"

	"github.com/Drafteame/draft/internal/pkg/exec"
)

// Code formats generated code using goimports-reviser
func Code() error {
	var err error

	spin := spinner.New().Title("Formatting generated code")

	action := func() {
		_, errExec := exec.Command("goimports-reviser ./...")
		if errExec != nil {
			err = fmt.Errorf("command 'goimports-reviser ./...' failed: %w", errExec)
		}
	}

	spinErr := spin.Action(action).Run()

	return errors.Join(spinErr, err)
}
