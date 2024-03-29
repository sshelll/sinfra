package functional

type iter[T, M, R any] struct {
	src    func() []T
	mapper func(T) T
	mapTo  func(T) M
	filter func(T) bool
}

// From, map to T, reduce to T, works with Collect and Reduce
func From[T any](arr []T) *iter[T, T, T] {
	return &iter[T, T, T]{
		src: func() []T {
			return arr
		},
	}
}

// FromM, map to M, works with CollectTo and MapTo
func FromM[T, M any](arr []T) *iter[T, M, M] {
	return &iter[T, M, M]{
		src: func() []T {
			return arr
		},
	}
}

// FromR, map to T, reduce to R, works with ReduceTo
func FromR[T, R any](arr []T) *iter[T, T, R] {
	return &iter[T, T, R]{
		src: func() []T {
			return arr
		},
	}
}

// FromMR, map to M, reduce to R, works with MapReduceTo
func FromMR[T, M, R any](arr []T) *iter[T, M, R] {
	return &iter[T, M, R]{
		src: func() []T {
			return arr
		},
	}
}

// Map, map to T, works with From, Collect and Reduce
func (it *iter[T, M, R]) Map(mapper func(T) T) *iter[T, M, R] {
	it.mapper = mapper
	return it
}

// MapTo, map to M, works with FromM, CollectTo and MapReduceTo
func (it *iter[T, M, R]) MapTo(mapper func(T) M) *iter[T, M, R] {
	it.mapTo = mapper
	return it
}

// Filter, filter elements
func (it *iter[T, M, R]) Filter(filter func(T) bool) *iter[T, M, R] {
	it.filter = filter
	return it
}

// Collect, map to T, collect to T
func (it *iter[T, M, R]) Collect() []T {
	var result []T
	if it.mapper == nil {
		it.mapper = func(i T) T { return i }
	}
	for _, v := range it.src() {
		if it.filter != nil && !it.filter(v) {
			continue
		}
		result = append(result, it.mapper(v))
	}
	return result
}

// CollectTo, map to M, collect to M
func (it *iter[T, M, R]) CollectTo() []M {
	var result []M
	if it.mapTo == nil {
		panic("mapTo function is not defined")
	}
	for _, v := range it.src() {
		if it.filter != nil && !it.filter(v) {
			continue
		}
		result = append(result, it.mapTo(v))
	}
	return result
}

// Reduce, map to T, reduce to T
func (it *iter[T, M, R]) Reduce(reducer func(T, T) T) T {
	var result T
	if it.mapper == nil {
		it.mapper = func(i T) T { return i }
	}
	arr := it.src()
	if len(arr) < 2 {
		panic("reduce requires at least 2 elements")
	}
	for _, v := range arr {
		if it.filter != nil && !it.filter(v) {
			continue
		}
		result = reducer(result, it.mapper(v))
	}
	return result
}

// ReduceTo, map to T, reduce to R
func (it *iter[T, M, R]) ReduceTo(reducer func(R, T) R) R {
	var result R
	if it.mapper == nil {
		it.mapper = func(i T) T { return i }
	}
	arr := it.src()
	if len(arr) < 2 {
		panic("reduce requires at least 2 elements")
	}
	for _, v := range arr {
		if it.filter != nil && !it.filter(v) {
			continue
		}
		result = reducer(result, it.mapper(v))
	}
	return result
}

// MapReduceTo, map to M, reduce to R
func (it *iter[T, M, R]) MapReduceTo(reducer func(R, M) R) R {
	var result R
	if it.mapTo == nil {
		panic("mapTo function is not defined")
	}
	arr := it.src()
	if len(arr) < 2 {
		panic("reduce requires at least 2 elements")
	}
	for _, v := range arr {
		if it.filter != nil && !it.filter(v) {
			continue
		}
		result = reducer(result, it.mapTo(v))
	}
	return result
}
