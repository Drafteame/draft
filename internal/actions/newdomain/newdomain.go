package newdomain

import (
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/templates"
)

type NewDomain struct {
	input dtos.DomainInput
	tmpl  templates.Domains
}

func New(input dtos.DomainInput) *NewDomain {
	return &NewDomain{
		input: input,
	}
}

func (nd *NewDomain) Exec() error {
	tmpl, err := templates.NewDomains(nd.input)
	if err != nil {
		return err
	}

	nd.tmpl = tmpl

	if errExec := nd.exec(); errExec != nil {
		return errExec
	}

	return nd.postCreate()
}
