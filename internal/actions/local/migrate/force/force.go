package force

type Input struct {
	WorkingDir         string
	Database           string
	LocalMigrateConfig string
	Group              string
	Version            int64
}

type Action struct {
	Input Input
}

func New(input Input) *Action {
	return &Action{input}
}
