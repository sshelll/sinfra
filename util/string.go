package util

import "strings"

func IsStrBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func IsStrPtrBlank(sp *string) bool {
	return sp == nil || len(strings.TrimSpace(*sp)) == 0
}
