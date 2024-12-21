package newdomain

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/internal/actions/newdomain"
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/forms"
)

var cmd = &cobra.Command{
	Use:   "new:domain",
	Short: "Create a new domain",
	Long:  "Create a new configurable domain, creates models, services, repositories and any other needed config to work",
	Run:   run,
}

func run(_ *cobra.Command, _ []string) {
	if data.Flags.WorkingDir != "" {
		if err := os.Chdir(data.Flags.WorkingDir); err != nil {
			panic(err)
		}
	}

	data.LoadMeta()

	input := dtos.DomainInput{}

	if err := forms.NewDomain(&input); err != nil {
		panic(err)
	}

	action, err := newdomain.GetAction(input)
	if err != nil {
		panic(err)
	}

	if err := action.Exec(); err != nil {
		panic(err)
	}

	println("Domain", input.DomainName, "created successfully")
}

func GetCmd() *cobra.Command {
	return cmd
}
