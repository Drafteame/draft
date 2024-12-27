package newservice

import (
	"errors"
)

func (ns *NewService) postCreate() error {
	switch ns.input.ServiceFramework {
	case "sls":
		return ns.updateSlsPlugins()
	default:
		return errors.New("unknown service framework")
	}
}
