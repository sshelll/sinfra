package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sshelll/sinfra/gtype"
)

func IsStrBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func IsStrPtrBlank(sp *string) bool {
	return sp == nil || len(strings.TrimSpace(*sp)) == 0
}

func StrToInteger[T gtype.Int](s string) T {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return T(i64)
}

func StrToUnsigned[T gtype.Uint](s string) T {
	u64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return T(u64)
}

func StrToFloat[T gtype.Float](s string) T {
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return T(f64)
}

func NumToStr[T gtype.Uint | gtype.Int | gtype.Float](num T) string {
	return fmt.Sprintf("%v", num)
}
