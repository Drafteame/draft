package newservice

import (
	"os"

	"github.com/charmbracelet/huh/spinner"

	"github.com/Drafteame/draft/internal/exec"
)

func (css *NewService) postCreate() error {
	if css.input.ServiceFramework == "sls" {
		css.postCreateSls()
		return nil
	}

	return nil
}

func (css *NewService) postCreateSls() {
	_ = spinner.New().
		Title("Updating serverless plugins...").
		Action(css.updateSlsPlugins).
		Run()
}

func (css *NewService) updateSlsPlugins() {
	if err := os.Chdir(css.input.ServicePath); err != nil {
		println("Error changing directory to", css.input.ServicePath, ":", err.Error())
		return
	}

	cmd := "npx npm-check-updates '/serverless-.*/' -u && npm install"

	out, code, err := exec.Command(cmd)
	if err != nil {
		println("Error updating serverless plugins", code, "-", err.Error(), ":")
		println(out)
	}
}
