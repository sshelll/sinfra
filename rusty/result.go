package rusty

type (
	Result[T any] struct {
		ok  ok[T]
		err error
	}

	ok[T any] struct {
		value T
	}
)

func OK[T any](value T) Result[T] {
	return Result[T]{
		ok: ok[T]{value: value},
	}
}

func Err[T any](err error) Result[T] {
	return Result[T]{
		err: err,
	}
}

func (res *Result[T]) Unwrap() T {
	if res.err != nil {
		panic(res.err)
	}
	return res.ok.value
}

func (res *Result[T]) UnwrapOr(def T) T {
	if res.err != nil {
		return def
	}
	return res.ok.value
}

func (res *Result[T]) IsOk() bool {
	return res.err == nil
}

func (res *Result[T]) IsErr() bool {
	return res.err != nil
}

func (res *Result[T]) Ok() Option[T] {
	if res.err != nil {
		return Option[T]{}
	}
	return Option[T]{
		some: &some[T]{value: res.ok.value},
	}
}

func (res *Result[T]) Err() Option[error] {
	if res.err == nil {
		return Option[error]{}
	}
	return Option[error]{
		some: &some[error]{value: res.err},
	}
}
