package setup

func (a *Action) Exec() error {
	if err := a.prune(); err != nil {
		return err
	}

	return a.init()
}
