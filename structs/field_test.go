package structs

import (
	"reflect"
	"testing"

	"github.com/sshelll/sinfra/conv"
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

type S1 struct {
	ID        *string           `required:"true"`
	Name      string            `required:"true,allowEmpty"`
	Interfase interface{}       `required:"true"`
	Slice     []int             `required:"true"`
	Map       map[string]string `required:"true,allowEmpty"`
	Anon      struct {
		Anon1 *int  `required:"true"`
		Anon2 []int `required:"true"`
		Anon3 string
	} `required:"true"`
	S2 *S2 `required:"true"`
}

type S2 struct {
	Age *int `required:"true"`
}

func TestCheckRequired(t *testing.T) {
	s := &S1{
		ID:   conv.StrPtr("id"),
		Name: "name",
		Anon: struct {
			Anon1 *int  `required:"true"`
			Anon2 []int `required:"true"`
			Anon3 string
		}{
			Anon1: conv.IntPtr(1),
			Anon2: []int{1, 2, 3},
			Anon3: "anon3",
		},
		S2: &S2{},
	}

	missing := CheckRequired(reflect.ValueOf(s), "required")
	t.Log(missing)

	s.ID = nil
	missing = CheckRequired(reflect.ValueOf(s), "required")
	t.Log(missing)

	s.Anon.Anon1 = nil
	missing = CheckRequired(reflect.ValueOf(s), "required")
	t.Log(missing)

	s.Anon = struct {
		Anon1 *int  `required:"true"`
		Anon2 []int `required:"true"`
		Anon3 string
	}{}
	missing = CheckRequired(reflect.ValueOf(s), "required")
	t.Log(missing)

	n := reflect.ValueOf(nil)
	missing = CheckRequired(n, "required")
}

func TestCheckRequiredWithEmptyField(t *testing.T) {
	s := &S1{
		ID:        conv.StrPtr("id"),
		Name:      "",
		Interfase: nil,
		Slice:     []int{},
		Map:       map[string]string{},
		Anon: struct {
			Anon1 *int  `required:"true"`
			Anon2 []int `required:"true"`
			Anon3 string
		}{
			Anon1: conv.IntPtr(1),
			Anon2: []int{1, 2, 3},
			Anon3: "anon3",
		},
		S2: &S2{},
	}

	missing := CheckRequired(reflect.ValueOf(s), "required")
	t.Log(missing)

	s.Map = nil
	missing = CheckRequired(reflect.ValueOf(s), "required")
	t.Log(missing)
}
