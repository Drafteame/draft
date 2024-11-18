package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Test mg.Namespace

// All Execute all unit testing.
func (Test) All() error {
	out, err := sh.Output("go", "test", "-v", "-race", "./...", "-cover", "-coverprofile=coverage.out")
	if err != nil {
		return err
	}

	if out != "" {
		println(out)
	}

	return err
}

// Cover Show HTML coverage output.
func (t Test) Cover() error {
	mg.Deps(t.All)
	return sh.Run("go", "tool", "cover", "-html", "coverage.out")
}
