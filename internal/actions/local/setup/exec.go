package setup

import "os"

func (a *Action) Exec() error {
	if a.Input.WorkingDir != "" {
		if err := os.Chdir(a.Input.WorkingDir); err != nil {
			return err
		}
	}

	if err := a.prune(); err != nil {
		return err
	}

	return a.init()
}
