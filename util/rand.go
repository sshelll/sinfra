package util

import (
	"bytes"
	"math/rand"
)

func RandBool(rd ...*rand.Rand) bool {
	rintn := rand.Intn
	if len(rd) > 0 && rd[0] != nil {
		rintn = rd[0].Intn
	}
	return rintn(2) == 0
}

// RandAsciiStr gen rand string which consists of ascii characters from 32 to 126.
func RandAsciiStr(length int, rd ...*rand.Rand) string {
	rintn := rand.Intn
	if len(rd) > 0 && rd[0] != nil {
		rintn = rd[0].Intn
	}
	buf := bytes.Buffer{}
	for i := 0; i < length; i++ {
		buf.WriteString(string(rune(32 + rintn(95))))
	}
	return buf.String()
}

//func init() {
//rand.Seed(time.Now().UnixMilli())
//}
