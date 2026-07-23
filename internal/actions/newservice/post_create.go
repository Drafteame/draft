package newservice

import (
	"errors"
)

func (ns *NewService) postCreate() error {
	switch ns.input.ServiceFramework {
	case "sls":
		if err := ns.fixOtel(); err != nil {
			return err
		}

		return ns.updateSlsPlugins()
	default:
		return errors.New("unknown service framework")
	}
}
