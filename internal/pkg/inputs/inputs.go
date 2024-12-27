package inputs

import (
	"os"

	"github.com/charmbracelet/huh"

	"github.com/Drafteame/draft/internal/data"
)

func run(input huh.Field) error {
	group := huh.NewGroup(input)
	form := huh.NewForm(group)

	if !data.Flags.TTY {
		form.WithInput(os.Stdin)
		form.WithOutput(os.Stdout)
	}

	return form.WithTheme(huh.ThemeCharm()).Run()
}
