package util

type Set map[interface{}]struct{}

func NewSet(list ...interface{}) Set {
	s := Set{}
	for _, i := range list {
		s[i] = struct{}{}
	}
	return s
}

func (s Set) Add(item interface{}) Set {
	s[item] = struct{}{}
	return s
}

func (s Set) ToIntArray() []int {
	arr := make([]int, 0, len(s))
	for i := range s {
		arr = append(arr, i.(int))
	}
	return arr
}

func (s Set) Contains(i interface{}) bool {
	_, ok := s[i]
	return ok
}
