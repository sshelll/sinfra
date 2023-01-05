package structs

import (
	"math/rand"
	. "reflect"
	"unsafe"

	"github.com/sshelll/sinfra/util"
)

type GenOption struct {
	Bool     *bool
	Intn     *int
	Int8n    *int8
	Int16n   *int16
	Int32n   *int32
	Int63n   *int64
	StrLen   *int
	SliceLen *int
	MapLen   *int
	ChanBuf  *int
	RandSeed *int64
	MaxDepth *int
}

func Gen(t Type, opt *GenOption) interface{} {

	var rd *rand.Rand

	if opt != nil && opt.RandSeed != nil {
		rd = rand.New(rand.NewSource(*opt.RandSeed))
	}

	rint, rintn := rand.Int, rand.Intn
	rint31, rint31n := rand.Int31, rand.Int31n
	rint63, rint63n := rand.Int63, rand.Int63n
	ruint32, ruint64 := rand.Uint32, rand.Uint64
	rfloat32, rfloat64 := rand.Float32, rand.Float64

	if rd != nil {
		rint, rintn = rd.Int, rd.Intn
		rint31, rint31n = rd.Int31, rd.Int31n
		rint63, rint63n = rd.Int63, rd.Int63n
		ruint32, ruint64 = rd.Uint32, rd.Uint64
		rfloat32, rfloat64 = rd.Float32, rd.Float64
	}

	maxDepth := 8
	if opt != nil && opt.MaxDepth != nil {
		maxDepth = *opt.MaxDepth
	}

	var gen func(Type, int) Value

	gen = func(t Type, depth int) Value {

		switch t.Kind() {

		case Bool:
			if opt != nil && opt.Bool != nil {
				return ValueOf(*opt.Bool)
			}
			return ValueOf(util.RandBool(rd))

		case Int:
			if opt != nil && opt.Intn != nil {
				return ValueOf(rintn(*opt.Intn))
			}
			return ValueOf(rint())

		case Int8:
			if opt != nil && opt.Int8n != nil {
				return ValueOf(int8(rintn(int(*opt.Int8n))))
			}
			return ValueOf(int8(rint()))

		case Int16:
			if opt != nil && opt.Int16n != nil {
				return ValueOf(int16(rintn(int(*opt.Int16n))))
			}
			return ValueOf(int16(rint()))

		case Int32:
			if opt != nil && opt.Int32n != nil {
				return ValueOf(rint31n(*opt.Int32n))
			}
			return ValueOf(rint31())

		case Int64:
			if opt != nil && opt.Int63n != nil {
				return ValueOf(rint63n(*opt.Int63n))
			}
			return ValueOf(rint63())

		case Uint:
			return ValueOf(uint(ruint32()))

		case Uint8:
			return ValueOf(uint8(ruint32()))

		case Uint16:
			return ValueOf(uint16(ruint32()))

		case Uint32:
			return ValueOf(ruint32())

		case Uint64:
			return ValueOf(ruint64())

		case Uintptr:
			return ValueOf(uintptr(rand.Uint32()))

		case Float32:
			return ValueOf(rfloat32())

		case Float64:
			return ValueOf(rfloat64())

		case Complex64:
			return ValueOf(complex64(complex(float64(rfloat32()), float64(rfloat32()))))

		case Complex128:
			return ValueOf(complex(rfloat64(), rfloat64()))

		case Chan:
			buf := 0
			if opt != nil && opt.ChanBuf != nil {
				buf = *opt.ChanBuf
			}
			return ValueOf(MakeChan(t, buf).Interface())

		case Func:
			return ValueOf(MakeFunc(t, func(args []Value) (results []Value) {
				numOut := t.NumOut()
				results = make([]Value, numOut)
				for i := 0; i < numOut; i++ {
					results[i] = ValueOf(gen(t.Out(i), depth))
				}
				return
			}).Interface())

		case Interface:
			return ValueOf("interface{}")

		case Map:
			size := 2
			if opt != nil && opt.MapLen != nil {
				size = *opt.MapLen
			}
			m := MakeMap(t)
			for i := 0; i < size; i++ {
				m.SetMapIndex(gen(t.Key(), depth), gen(t.Elem(), depth))
			}
			return ValueOf(m.Interface())

		case Array, Slice:
			size := 2
			if opt != nil && opt.SliceLen != nil {
				size = *opt.SliceLen
			}
			slice := MakeSlice(t, size, size)
			for i := 0; i < size; i++ {
				slice.Index(i).Set(gen(t.Elem(), depth))
			}
			return ValueOf(slice.Interface())

		case String:
			if opt != nil && opt.StrLen != nil {
				return ValueOf(util.RandAsciiStr(*opt.StrLen, rd))
			}
			return ValueOf(util.RandAsciiStr(16, rd))

		case UnsafePointer:
			i := 0
			return ValueOf(unsafe.Pointer(&i))

		case Struct:
			v := New(t).Elem()
			if depth <= maxDepth {
				for i, n := 0, t.NumField(); i < n; i++ {
					f := v.Field(i)
					f.Set(gen(f.Type(), depth+1))
				}
			}
			return ValueOf(v.Interface())

		case Ptr:
			v := New(t.Elem())
			v.Elem().Set(gen(t.Elem(), depth))
			return v

		default:
			panic("unexpected reflect.Kind " + t.Kind().String())

		}
	}

	return gen(t, 1).Interface()

}
