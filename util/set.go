package util

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](list ...T) Set[T] {
	s := Set[T]{}
	for _, i := range list {
		s[i] = struct{}{}
	}
	return s
}

func (s Set[T]) Add(item T) Set[T] {
	s[item] = struct{}{}
	return s
}

func (s Set[T]) Contains(i T) bool {
	_, ok := s[i]
	return ok
}

func (s Set[T]) ToArray() []T {
	arr := make([]T, 0, len(s))
	for i := range s {
		arr = append(arr, i)
	}
	return arr
}
