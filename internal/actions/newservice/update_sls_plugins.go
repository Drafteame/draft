package newservice

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/huh/spinner"

	"github.com/Drafteame/draft/internal/pkg/exec"
)

func (ns *NewService) updateSlsPlugins() error {
	var err error

	action := func() {
		if errChdir := os.Chdir(ns.input.ServicePath); errChdir != nil {
			err = errors.Join(errors.New("could not change to service directory"), errChdir)
			return
		}

		cmd := "npx npm-check-updates '/serverless-.*/' -u && npm install"

		out, code, errCmd := exec.Command(cmd)
		if errCmd != nil {
			err = fmt.Errorf("%w: cannot upgrade serverless plugins [code %d]: %s", errCmd, code, out)
			return
		}
	}

	spinErr := spinner.New().
		Title("Updating serverless plugins...").
		Action(action).
		Run()

	return errors.Join(spinErr, err)
}
