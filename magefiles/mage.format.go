// nolint
package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mgechev/revive/cli"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/revivelib"
	"magefiles/reviveextrarules"
)

// Vet execute `go vet` checks.
func Vet() error {
	command := "go"
	args := []string{"vet", "./..."}

	out, err := sh.Output(command, args...)

	if out != "" {
		println(out)
	}

	return err
}

// Lint Runs revive checks over the code.
func Lint() error {
	mg.Deps(Revive, Vet)

	return nil
}

// Format Runs gofmt over the code.
func Format() error {
	outImp, err := sh.Output("goimports-reviser", "-format", "-excludes=bin,node_modules,tmp,.git,magefiles", "./...")
	if err != nil {
		return err
	}

	if outImp != "" {
		println(outImp)
	}

	return nil
}

// Revive runs revive checks over the code.
func Revive() {
	os.Args = []string{"revive", "-config=revive.toml", "-formatter=friendly", "./..."}
	extraRules := []revivelib.ExtraRule{
		{Rule: &reviveextrarules.IncorrectConfigImport{}, DefaultConfig: lint.RuleConfig{}},
	}
	cli.RunRevive(extraRules...)
}
