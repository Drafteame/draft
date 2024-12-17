package data

import (
	"github.com/Drafteame/draft/internal/pkg/files"
	"github.com/Drafteame/draft/internal/project"
)

type meta struct {
	PackageName string
}

var Meta = meta{}

func LoadMeta() {
	setPackageName()
}

func setPackageName() {
	if !files.Exists("go.mod") {
		panic("go.mod file not found")
	}

	name, err := project.GetPackage("go.mod")
	if err != nil {
		panic(err)
	}

	Meta.PackageName = name
}
