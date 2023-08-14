package util

func AllowPanic(fn func()) (r interface{}) {
	defer func() {
		r = recover()
	}()
	fn()
	return
}
