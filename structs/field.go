package structs

import (
	"reflect"
	"strings"
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

func CheckRequired(v reflect.Value, tag string) (missing []string) {
	if !v.IsValid() {
		return
	}

	var check func(name string, v reflect.Value, tag string) (missing []string)
	check = func(name string, v reflect.Value, tag string) (missing []string) {
		for v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return
			}
			v = v.Elem()
		}

		if v.Kind() != reflect.Struct {
			return
		}

		for i, n := 0, v.NumField(); i < n; i++ {
			fv := v.Field(i)
			ft := v.Type().Field(i)

			if tagv := ft.Tag.Get(tag); tagv != "true" {
				continue
			}

			if ft.Anonymous || ft.Type.Kind() == reflect.Struct {
				missing = append(missing, check(name+"."+ft.Name, fv, tag)...)
				continue
			}

			switch fv.Kind() {
			case reflect.Invalid:
				missing = append(missing, name+"."+ft.Name)
			case reflect.Interface:
				if fv.IsNil() {
					missing = append(missing, name+"."+ft.Name)
				}
			case reflect.Ptr:
				if fv.IsNil() {
					missing = append(missing, name+"."+ft.Name)
				}
				missing = append(missing, check(name+"."+ft.Name, fv, tag)...)
			case reflect.Slice, reflect.Array, reflect.Map:
				if fv.Len() == 0 {
					missing = append(missing, name+"."+ft.Name)
				}
			}
		}

		return
	}

	typeName := v.Type().String()
	splited := strings.Split(typeName, ".")
	if len(splited) > 0 {
		typeName = splited[len(splited)-1]
	}
	if len(typeName) == 0 {
		typeName = "$"
	}
	missing = check(typeName, v, tag)

	return
}
