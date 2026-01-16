package newdomain

import (
	"errors"

	"github.com/charmbracelet/huh/spinner"

	"github.com/Drafteame/draft/internal/actions/mockery"
)

func (nd *NewDomain) mockery() error {
	var err error

	spin := spinner.New().Title("Generating mocks")

	action := func() {
		spin.Update("Running mockery to create mocks")
		err = nd.runMockery()
		if err != nil {
			return
		}
	}

	spinErr := spin.Action(action).Run()

	return errors.Join(spinErr, err)
}

func (nd *NewDomain) runMockery() error {
	// Build the config file paths for both service and repository
	serviceMockeryPath := nd.input.DomainPath + "/service/.mockery.pkg.yml"
	repositoryMockeryPath := nd.input.DomainPath + "/repository/.mockery.pkg.yml"

	configFiles := []string{serviceMockeryPath, repositoryMockeryPath}

	// Run mockery action with config files
	// Using jobsNum=2 (one for service, one for repository)
	// dry=false (actually execute mockery)
	// gitMod=false (not using git diff mode)
	m := mockery.New(nd.ctx, configFiles, 2, false, false)

	return m.Exec()
}
