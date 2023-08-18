package structs

import (
	"reflect"
	"strings"

	"github.com/sshelll/sinfra/util"
)

// ExtractStructFieldValues this method will extract real val if a field is a non-nil ptr.
func ExtractStructFieldValues(v reflect.Value) (fvs []interface{}) {
	if !v.IsValid() {
		panic("reflect value is invalid: " + v.Kind().String())
	}

	switch v.Kind() {

	case reflect.Invalid:
		panic("there is a field with invalid type in model: " + v.Type().Name())

	case reflect.Ptr:
		if v.IsNil() {
			fvs = append(fvs, nil)
		} else {
			fvs = append(fvs, ExtractStructFieldValues(v.Elem())...)
		}

	case reflect.Struct:
		for i, n := 0, v.NumField(); i < n; i++ {
			fv := v.Field(i)
			fvs = append(fvs, ExtractStructFieldValues(fv)...)
		}

	default:
		fvs = append(fvs, v.Interface())

	}

	return

}

// CheckRequired checks if a ptr field was set in the give struct, it returns a slice of missing field names.
// The tag name can be customized, like "required", and the tag value must be 'true' when a field is required.
// Also, you can use "allowEmpty" on string / slice / map / array fields to allow empty value, like `required:"true,allowEmpty"`.
func CheckRequired(v reflect.Value, tag string) (missing []string) {
	callback := func(name string, fv reflect.Value, ft reflect.StructField) {
		// check tag
		tag := ft.Tag.Get(tag)
		tags := strings.Split(tag, ",")
		if n := len(tags); n == 0 || n > 2 || tags[0] != "true" {
			return
		}
		allowEmpty := len(tags) == 2 && tags[1] == "allowEmpty"

		// precheck if the field is absolutly nil
		if util.IsNil(fv) {
			missing = append(missing, name)
			return
		}

		// check field value with opts
		switch fv.Kind() {
		case reflect.Invalid:
			missing = append(missing, name)
		case reflect.Ptr:
			missing = append(missing, CheckRequired(fv, tag)...)
		case reflect.String:
			if !allowEmpty && util.IsStrBlank(fv.String()) {
				missing = append(missing, name)
			}
		case reflect.Slice:
			if !allowEmpty && fv.Len() == 0 {
				missing = append(missing, name)
			}
		case reflect.Map:
			if !allowEmpty && fv.Len() == 0 {
				missing = append(missing, name)
			}
		case reflect.Array:
			if !allowEmpty && fv.Len() == 0 {
				missing = append(missing, name)
			}
		}
	}
	WalkStruct(v, callback)
	return
}

// WalkStruct walks each field of a struct recursively.
// The callback function will be called for each field.
// Param of callback:
// name is the full name of the field, starts with package name,
// fv is the reflect.Value of the field,
// ft is the type of the field.
func WalkStruct(v reflect.Value, callback func(name string, fv reflect.Value, ft reflect.StructField)) {
	if !v.IsValid() {
		return
	}

	var isNil bool
	if v, isNil = util.Indirect(v); isNil {
		return
	}

	if v.Kind() != reflect.Struct {
		return
	}

	var walk func(name string, v reflect.Value)
	walk = func(name string, v reflect.Value) {
		for i, n := 0, v.NumField(); i < n; i++ {
			fv := v.Field(i)
			ft := v.Type().Field(i)
			fname := name + "." + ft.Name
			if ft.Anonymous || ft.Type.Kind() == reflect.Struct {
				walk(fname, fv)
			}
			callback(fname, fv, ft)
		}
	}

	typeName := v.Type().String()
	if len(typeName) == 0 {
		typeName = "$unknown"
	}
	walk(typeName, v)
}
