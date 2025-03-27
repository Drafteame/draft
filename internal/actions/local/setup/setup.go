package setup

type Input struct {
	WorkingDir   string
	Prune        bool
	BypassDocker bool
}

type Action struct {
	Input Input
}

func New(input Input) *Action {
	return &Action{
		Input: input,
	}
}
