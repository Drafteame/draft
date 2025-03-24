package newdomain

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh/spinner"

	"github.com/Drafteame/draft/internal/pkg/exec"
)

func (nd *NewDomain) format() error {
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
