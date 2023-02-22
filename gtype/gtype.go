package gtype

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Int_ interface {
	int | int8 | int16 | int32 | int64
}

type Uint interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Uint_ interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type Float interface {
	~float32 | ~float64
}

type Float_ interface {
	float32 | float64
}

type Complex interface {
	~complex64 | ~complex128
}

type Complex_ interface {
	complex64 | complex128
}

type Number interface {
	Int | Uint | Float | Complex
}

type Number_ interface {
	Int_ | Uint_ | Float_ | Complex_
}

type Integer interface {
	Int | Uint
}

type Integer_ interface {
	Int_ | Uint_
}

type Sortable interface {
	Integer | Float | ~string
}

type Sortable_ interface {
	Integer_ | Float_ | string
}

type Nonsense interface {
	int
	float32
}
