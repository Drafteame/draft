package inputs

import "github.com/charmbracelet/huh"

func Select[T comparable](title string, opts ...Option[T]) error {
	input := huh.NewSelect[T]().Title(title)

	inputOpts := options[T]{}

	for _, opt := range opts {
		opt(&inputOpts)
	}

	if inputOpts.description != "" {
		input.Description(inputOpts.description)
	}

	selectOpts := make([]huh.Option[T], 0, len(opts))

	for optKey, optVal := range inputOpts.options {
		selectOpts = append(selectOpts, huh.NewOption(optKey, optVal))
	}

	input.Options(selectOpts...)

	if inputOpts.value != nil {
		input.Value(inputOpts.value)
	}

	if inputOpts.validation != nil {
		input.Validate(inputOpts.validation)
	}

	return run(input)
}
