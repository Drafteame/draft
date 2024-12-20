package inputs

import (
	"github.com/charmbracelet/huh"
)

func Text(title string, opts ...Option[string]) error {
	input := huh.NewInput().Title(title)

	inputOpts := options[string]{}

	for _, opt := range opts {
		opt(&inputOpts)
	}

	if inputOpts.description != "" {
		input.Description(inputOpts.description)
	}

	if inputOpts.value != nil {
		input.Value(inputOpts.value)
	}

	if inputOpts.validation != nil {
		input.Validate(inputOpts.validation)
	}

	return run(input)
}
