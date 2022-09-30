package structs

import (
	"reflect"
	"strings"
)

func ExtractStructTags(t reflect.Type, tagKey string, keyCb func(string) string) (tagList []string) {

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		panic("kind invalid, expected: reflect.Struct, actual: " + t.Kind().String())
	}

	for i, n := 0, t.NumField(); i < n; i++ {

		f := t.Field(i)

		// is struct
		if f.Type.Kind() == reflect.Struct && f.Anonymous {
			tagList = append(tagList, ExtractStructTags(f.Type, tagKey, keyCb)...)
			continue
		}

		// is struct ptr
		if f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct && f.Anonymous {
			tagList = append(tagList, ExtractStructTags(f.Type.Elem(), tagKey, keyCb)...)
			continue
		}

		// is basic type
		tag := f.Tag.Get(tagKey)
		if tag != "" {
			if keyCb != nil {
				tag = keyCb(tag)
			}
			tagList = append(tagList, tag)
		}

	}

	return

}

func ExtractGormColumnName(t reflect.Type) []string {
	return ExtractStructTags(t, "gorm", func(s string) string {
		if s == "primary_key" {
			return "id"
		}
		return strings.Split(s, "column:")[1]
	})
}

func ExtractJsonKey(t reflect.Type) []string {
	return ExtractStructTags(t, "json", func(s string) string {
		return strings.Split(s, ",")[0]
	})
}
