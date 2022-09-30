package structs

import "reflect"

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
