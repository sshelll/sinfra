package util

import (
	. "reflect"
)

func GetSliceElemKind(st Type) Kind {
	mustBeSliceOrArray(st)
	return MakeSlice(st, 1, 1).Index(0).Kind()
}

func GetSliceElemType(st Type) Type {
	mustBeSliceOrArray(st)
	return MakeSlice(st, 1, 1).Index(0).Type()
}

func mustBeSliceOrArray(t Type) {
	if k := t.Kind(); k != Slice && k != Array {
		panic("type invalid, expected: reflect.Array or reflect.Slice, actual: " + k.String())
	}
}
