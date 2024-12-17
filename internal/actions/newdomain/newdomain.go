package newdomain

import (
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/templates"
)

type NewDomain struct {
	input dtos.DomainInput
	tmpl  templates.Domains
}

func GetAction(input dtos.DomainInput) (*NewDomain, error) {
	tmpl, err := templates.NewDomains(input)
	if err != nil {
		return nil, err
	}

	return &NewDomain{
		input: input,
		tmpl:  tmpl,
	}, nil
}

func (nd *NewDomain) Exec() error {
	if err := nd.exec(); err != nil {
		return err
	}

	return nd.postCreate()
}
