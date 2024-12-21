package newdomain

func (nd *NewDomain) postCreate() error {
	actions := []func() error{
		nd.postgresModels,
		nd.mockery,
	}

	for _, action := range actions {
		if err := action(); err != nil {
			return err
		}
	}

	return nil
}
