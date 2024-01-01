package rusty

type (
	Option[T any] struct {
		some *some[T]
		none none[T]
	}

	some[T any] struct {
		value T
	}

	none[T any] struct{}
)

func Some[T any](value T) Option[T] {
	return Option[T]{
		some: &some[T]{value: value},
	}
}

func None[T any]() Option[T] {
	return Option[T]{}
}

func (opt Option[T]) Unwrap() T {
	if opt.some == nil {
		panic("unwrap on none")
	}
	return opt.some.value
}

func (opt Option[T]) UnwrapOr(def T) T {
	if opt.some == nil {
		return def
	}
	return opt.some.value
}

func (opt Option[T]) Take() Option[T] {
	if opt.some == nil {
		return Option[T]{}
	}
	opt.some = nil
	return Option[T]{
		some: &some[T]{value: opt.some.value},
	}
}

func (opt Option[T]) IsSome() bool {
	return opt.some != nil
}

func (opt Option[T]) IsNone() bool {
	return opt.some == nil
}
