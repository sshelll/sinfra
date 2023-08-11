package structs

import (
	"reflect"
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
func CheckRequired(v reflect.Value, tag string) (missing []string) {
	var callback func(name string, fv reflect.Value, ft reflect.StructField)
	callback = func(name string, fv reflect.Value, ft reflect.StructField) {
		if ft.Tag.Get(tag) != "true" {
			return
		}
		switch fv.Kind() {
		case reflect.Invalid:
			missing = append(missing, name)
		case reflect.Interface:
			if fv.IsNil() {
				missing = append(missing, name)
			}
		case reflect.Ptr:
			if fv.IsNil() {
				missing = append(missing, name)
			}
			missing = append(missing, CheckRequired(fv, tag)...)
		case reflect.Slice, reflect.Array, reflect.Map:
			if fv.Len() == 0 {
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

	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
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
