package error

func Drop1(args ...interface{}) interface{} {
	return args[0]
}

func DropN(n int, args ...interface{}) []interface{} {
	out := make([]interface{}, n)
	for i := 0; i < n; i++ {
		out[i] = args[i]
	}
	return out
}
