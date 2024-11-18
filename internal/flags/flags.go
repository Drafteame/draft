package flags

import (
	"github.com/spf13/pflag"
)

type flags struct {
	WorkingDir string
}

var Flags = flags{}

func init() {
	pflag.StringVarP(&Flags.WorkingDir, "working-dir", "w", "", "Working directory")
	pflag.Parse()
}
