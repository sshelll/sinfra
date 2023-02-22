package util

// any to ptr
func Ptr[T any](v T) *T {
	return &v
}
