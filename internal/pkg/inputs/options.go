package inputs

type options[T any] struct {
	value       *T
	options     map[string]T
	validation  func(val T) error
	description string
	placeholder string
	saveKey     bool
}

type Option[T any] func(*options[T])

func WithDescription[T any](desc string) Option[T] {
	return func(o *options[T]) {
		o.description = desc
	}
}

func WithValidation[T any](validation func(val T) error) Option[T] {
	return func(o *options[T]) {
		o.validation = validation
	}
}

func WithValue[T any](value *T) Option[T] {
	return func(o *options[T]) {
		o.value = value
	}
}

func WithOptions[T any](opts map[string]T) Option[T] {
	return func(o *options[T]) {
		o.options = opts
	}
}

func WithPlaceholder[T any](placeholder string) Option[T] {
	return func(o *options[T]) {
		o.placeholder = placeholder
	}
}

func WithSaveKey[T any]() Option[T] {
	return func(o *options[T]) {
		o.saveKey = true
	}
}
