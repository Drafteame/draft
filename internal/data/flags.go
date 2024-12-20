package data

type flags struct {
	WorkingDir string
	Debug      bool
	TTY        bool
}

var Flags = flags{}
