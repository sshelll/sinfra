package util

import "github.com/sshelll/sinfra/gtype"

func Min[T gtype.Sortable](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T gtype.Sortable](a, b T) T {
	if a > b {
		return a
	}
	return b
}
