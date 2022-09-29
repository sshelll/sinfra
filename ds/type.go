package ds

import (
	"errors"
	"fmt"
	"reflect"
)

type Comparator func(a, b interface{}) int

type TypeChecker func(elem interface{}) bool

const defaultCapacity = 10

var (
	intType       = reflect.TypeOf(0)
	intComparator = func(a, b interface{}) int {
		aInt, bInt := a.(int), b.(int)
		return aInt - bInt
	}
	intTypeChecker = func(elem interface{}) bool {
		return intType == reflect.TypeOf(elem)
	}
	typeMismatchErr = errors.New("type mismatch")
)

func newOutOfBoundsErr(idx int) error {
	return errors.New(fmt.Sprintf("index: %d is out of bounds", idx))
}
