package format

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh/spinner"

	"github.com/Drafteame/draft/internal/pkg/exec"
)

// Code formats generated code using goimports-reviser.
// If no paths are provided, formats all files (./...).
// Otherwise, formats only the specified paths.
func Code(paths ...string) error {
	var err error

	spin := spinner.New().Title("Formatting generated code")

	action := func() {
		target := "./..."
		if len(paths) > 0 {
			target = ""
			for i, path := range paths {
				if i > 0 {
					target += " "
				}
				target += path
			}
		}

		_, errExec := exec.Command(fmt.Sprintf("goimports-reviser %s", target))
		if errExec != nil {
			err = fmt.Errorf("command 'goimports-reviser %s' failed: %w", target, errExec)
		}
	}

	spinErr := spin.Action(action).Run()

	return errors.Join(spinErr, err)
}
