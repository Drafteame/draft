package up

type Input struct {
	WorkingDir         string
	Database           string
	LocalMigrateConfig string
	Group              string
	All                bool
}

type Action struct {
	Input Input
}

func New(input Input) *Action {
	return &Action{input}
}
