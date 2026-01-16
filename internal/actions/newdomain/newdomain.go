package newdomain

import (
	"context"

	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/templates"
)

type NewDomain struct {
	ctx   context.Context
	input dtos.DomainInput
	tmpl  templates.Domains
}

func New(ctx context.Context, input dtos.DomainInput) *NewDomain {
	return &NewDomain{
		ctx:   ctx,
		input: input,
	}
}

func (nd *NewDomain) Exec() error {
	tmpl, err := templates.NewDomains(nd.input)
	if err != nil {
		return err
	}

	nd.tmpl = tmpl

	if err := nd.preCreate(); err != nil {
		return err
	}

	if err := nd.exec(); err != nil {
		return err
	}

	return nd.postCreate()
}

// preCreate performs validation and setup before domain creation
func (nd *NewDomain) preCreate() error {
	// Future validations can be added here
	// Example: check if domain path already exists
	return nil
}
