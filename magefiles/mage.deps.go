package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Deps namespace for dependency management
type Deps mg.Namespace

// InstallHooks install Husky hooks.
func (Deps) InstallHooks() error {
	return sh.Run("husky", "install")
}
