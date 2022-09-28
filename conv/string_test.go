package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrVal(t *testing.T) {
	var v string
	var flag = WithNilSptrConvFlag | WithTrimSptrConvFlag
	v = StrVal(nil, flag)
	assert.Equal(t, "<nil>", v)
	v = StrVal(StrPtr(" v "), flag)
	assert.Equal(t, "v", v)
}
