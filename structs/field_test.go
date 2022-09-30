package structs

import (
	"reflect"
	"testing"

	"github.com/SCU-SJL/sinfra/conv"
	"github.com/stretchr/testify/assert"
)

func TestExtractStructFieldValues(t *testing.T) {

	type base struct {
		ID string
	}

	s := struct {
		*base
		Name     *string
		Age      int
		Birthday int64
		Map      map[string]int
	}{
		base: &base{
			ID: "id",
		},
		Name:     conv.StrPtr("name"),
		Age:      1,
		Birthday: 1,
		Map: map[string]int{
			"k1": 1,
			"k2": 2,
			"k3": 3,
		},
	}

	fvs := ExtractStructFieldValues(reflect.ValueOf(s))
	assert.Equal(t, "id", fvs[0])
	assert.Equal(t, "name", fvs[1])
	assert.Equal(t, 1, fvs[2])
	assert.Equal(t, int64(1), fvs[3])
	assert.Equal(t,
		map[string]int{
			"k1": 1,
			"k2": 2,
			"k3": 3,
		}, fvs[4])

}

func TestExtractNormalTypeValue(t *testing.T) {
	fvs := ExtractStructFieldValues(reflect.ValueOf(3))
	assert.Equal(t, 1, len(fvs))
	assert.Equal(t, 3, fvs[0])
}
