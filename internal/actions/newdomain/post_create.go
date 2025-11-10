package newdomain

func (nd *NewDomain) postCreate() error {
	actions := []func() error{
		nd.mockery,
		nd.format,
	}

	// Only add postgres models for postgres DB
	if nd.input.DBType == "postgres" {
		actions = append([]func() error{nd.postgresModels}, actions...)
	}

	for _, action := range actions {
		if err := action(); err != nil {
			return err
		}
	}

	return nil
}
