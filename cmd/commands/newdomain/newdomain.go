package newdomain

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/internal/actions/newdomain"
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/forms"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var newDomainCmd = &cobra.Command{
	Use:   "new:domain",
	Short: "Create a new domain",
	Long:  "Create a new configurable domain, creates models, services, repositories and any other needed config to work",
	Run:   run,
}

func run(cmd *cobra.Command, _ []string) {
	common.ChDir(cmd)

	data.LoadMeta()

	input := dtos.DomainInput{}

	if err := forms.NewDomain(&input); err != nil {
		log.Exitf(1, "Failed colect domain info: %v", err)
	}

	if errExec := newdomain.New(input).Exec(); errExec != nil {
		log.Exitf(1, "Failed to create domain: %v", errExec)
	}

	log.Successf("Domain %s created successfully", input.DomainName)
}

func GetCmd() *cobra.Command {
	return newDomainCmd
}
