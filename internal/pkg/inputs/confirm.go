package inputs

import "github.com/charmbracelet/huh"

func Confirm(title string, opts ...Option[bool]) error {
	input := huh.NewConfirm().Title(title)

	inputOpts := options[bool]{}

	for _, opt := range opts {
		opt(&inputOpts)
	}

	if inputOpts.description != "" {
		input.Description(inputOpts.description)
	}

	if inputOpts.value != nil {
		input.Value(inputOpts.value)
	}

	return run(input)
}
