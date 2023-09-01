package structs

import (
	"testing"

	. "reflect"

	"github.com/sshelll/sinfra/util"
	"github.com/stretchr/testify/assert"
)

func TestGenMap(t *testing.T) {
	mt := TypeOf(map[string]int{})
	m := Gen(mt, &GenOption{MapLen: util.Ptr(10)}).(map[string]int)
	assert.Equal(t, 10, len(m))
}

func TestGenSlice(t *testing.T) {
	st := TypeOf([]string{})
	_ = Gen(st, nil)
}

func TestGenNum(t *testing.T) {

	_ = Gen(TypeOf(1), nil).(int)
	_ = Gen(TypeOf(int8(1)), nil).(int8)
	_ = Gen(TypeOf(int16(1)), nil).(int16)
	_ = Gen(TypeOf(int32(1)), nil).(int32)
	_ = Gen(TypeOf(int64(1)), nil).(int64)

	_ = Gen(TypeOf(uint(1)), nil).(uint)
	_ = Gen(TypeOf(uint8(1)), nil).(uint8)
	_ = Gen(TypeOf(uint16(1)), nil).(uint16)
	_ = Gen(TypeOf(uint32(1)), nil).(uint32)
	_ = Gen(TypeOf(uint64(1)), nil).(uint64)
	_ = Gen(TypeOf(uintptr(1)), nil).(uintptr)

	_ = Gen(TypeOf(float32(1.0)), nil).(float32)
	_ = Gen(TypeOf(float64(1.0)), nil).(float64)

	_ = Gen(TypeOf((complex(1.0, 1.0))), nil).(complex128)
	_ = Gen(TypeOf(complex64(complex(1.0, 1.0))), nil).(complex64)

}

func TestGenStruct(t *testing.T) {

	type SubStruct struct {
		Name string
	}

	type TestStruct struct {
		Int         int
		IntPtr      *int
		Str         string
		StrPtr      *string
		StrSlice    []string
		StrPtrSlice []*string
		Map         map[string]int
		Sub         *SubStruct
		SubSlice    []*SubStruct
	}

	inst := Gen(TypeOf(TestStruct{}), nil).(TestStruct)
	t.Logf("%v", inst)
	t.Log(inst.Sub)
	t.Log(inst.SubSlice[0].Name)
	t.Log(inst.SubSlice[1].Name)

	inst = Gen(TypeOf(TestStruct{}), &GenOption{MaxDepth: util.Ptr(1)}).(TestStruct)
	t.Logf("%+v", inst)
	t.Log(inst.Sub)

}
