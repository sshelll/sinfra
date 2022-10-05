package conv

func Ptr(i interface{}) interface{} {
	if i == nil {
		panic("cannot get ptr of nil")
	}
	return &i
}

func BoolPtr(b bool) *bool {
	return &b
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

func UintPtr(u uint) *uint {
	return &u
}

func Uint8Ptr(u uint8) *uint8 {
	return &u
}

func Uint16Ptr(u uint16) *uint16 {
	return &u
}

func Uint32Ptr(u uint32) *uint32 {
	return &u
}

func Uint64Ptr(u uint64) *uint64 {
	return &u
}

func StrPtr(s string) *string {
	return &s
}

func IntElemOrDefault(p *int, i int) int {
	if p != nil {
		return *p
	}
	return i
}

func Int64ElemOrDefault(p *int64, i int64) int64 {
	if p != nil {
		return *p
	}
	return i
}

func StrElemOrDefault(p *string, s string) string {
	if p != nil {
		return *p
	}
	return s
}
