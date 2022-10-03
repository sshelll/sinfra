package structs

import (
	"reflect"
	"testing"
)

func TestXxx(t *testing.T) {

	s := struct {
		Ints []int
	}{}

	v := reflect.New(reflect.TypeOf(s)).Elem()

	f := v.Field(0)
	arrV := reflect.MakeSlice(f.Type(), 2, 2)
	t.Logf("field type = %s", f.Type().Name())
	t.Logf("arr elem type = %s", arrV.Index(0).Kind().String())

}
