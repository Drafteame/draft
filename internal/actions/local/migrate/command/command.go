package command

type Input struct {
	Command            string
	WorkingDir         string
	Database           string
	LocalMigrateConfig string
	Group              string
	All                bool
	Version            int64
}

type Action struct {
	Input Input
}

func New(input Input) *Action {
	return &Action{input}
}
