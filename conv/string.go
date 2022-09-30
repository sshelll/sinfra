package conv

import (
	"bytes"
)

func StrVal(p *string) string {
	if p == nil {
		return ""
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
