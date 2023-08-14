package util

import "runtime"

func GetStackInfo() (stackInfo []byte) {
	stackInfo = make([]byte, 1<<16)
	n := runtime.Stack(stackInfo, true)
	return stackInfo[:n]
}
