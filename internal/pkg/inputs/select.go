package inputs

import (
	"sort"

	"github.com/charmbracelet/huh"
)

func Select[T comparable](title string, opts ...Option[T]) error {
	input := huh.NewSelect[T]().Title(title)

	inputOpts := options[T]{}

	for _, opt := range opts {
		opt(&inputOpts)
	}

	if inputOpts.description != "" {
		input.Description(inputOpts.description)
	}

	keys := make([]string, 0, len(inputOpts.options))
	for k := range inputOpts.options {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	selectOpts := make([]huh.Option[T], 0, len(inputOpts.options))
	for _, optKey := range keys {
		selectOpts = append(selectOpts, huh.NewOption(optKey, inputOpts.options[optKey]))
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
