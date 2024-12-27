package data

type flags struct {
	WorkingDir string
	Debug      bool
	TTY        bool
	NoSentry   bool
}

var Flags = flags{}
