package conv

func Ptr(i interface{}) interface{} {
	if i == nil {
		panic("cannot get ptr of nil")
	}
	return &i
}

func IntPtr(i int) *int {
	return &i
}

func Int8Ptr(i int8) *int8 {
	return &i
}

func BytePtr(b byte) *byte {
	return &b
}

func Int16Ptr(i int16) *int16 {
	return &i
}

func Int32Ptr(i int32) *int32 {
	return &i
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func StrPtr(s string) *string {
	return &s
}
