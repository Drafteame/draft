package setup

type Input struct {
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
