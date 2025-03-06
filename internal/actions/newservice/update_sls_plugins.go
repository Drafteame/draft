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

		_, errCmd := exec.Command(cmd)
		if errCmd != nil {
			err = fmt.Errorf("cannot upgrade serverless plugins: %w", errCmd)
			return
		}
	}

	spinErr := spinner.New().
		Title("Updating serverless plugins...").
		Action(action).
		Run()

	return errors.Join(spinErr, err)
}
