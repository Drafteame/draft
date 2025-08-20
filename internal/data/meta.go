package data

import (
	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/pkg/log"
	"github.com/Drafteame/draft/internal/project"
)

type meta struct {
	PackageName string
}

var Meta = meta{}

func LoadMeta() {
	if !files.Exists("go.mod") {
		log.Exit(1, "go.mod file not found")
	}

	name, err := project.GetPackage("go.mod")
	if err != nil {
		log.Exitf(1, "failed to get package name: %s", err.Error())
	}

	Meta.PackageName = name
}
