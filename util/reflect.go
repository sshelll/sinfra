package util

import (
	"reflect"
)

func Indirect(v reflect.Value) (elem reflect.Value, isNil bool) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return v, true
		}
		v = v.Elem()
	}
	return v, false
}

func IsNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr,
		reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}
	return false
}

func IsNilT[T any](v T) bool {
	return IsNil(reflect.ValueOf(v))
}
