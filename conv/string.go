package conv

import (
	"bytes"
	"strings"
)

type sptrConvFlag byte

const (
	DefaultSptrConvFlag sptrConvFlag = 1 << iota
	WithNilSptrConvFlag
	WithTrimSptrConvFlag
)

func StrPtr(s string) *string {
	return &s
}

func StrVal(p *string, convFlag sptrConvFlag) string {
	if p == nil {
		if WithNilSptrConvFlag&convFlag > 0 {
			return "<nil>"
		}
		return ""
	}
	if WithTrimSptrConvFlag&convFlag > 0 {
		return strings.TrimSpace(*p)
	}
	return *p
}

func StrConcat(sep string, strList ...string) string {
	if len(strList) == 0 {
		return ""
	}
	buf := bytes.Buffer{}
	for _, s := range strList[:len(strList)-1] {
		buf.WriteString(s)
		buf.WriteString(sep)
	}
	buf.WriteString(strList[len(strList)-1])
	return buf.String()
}
