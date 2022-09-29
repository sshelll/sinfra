package ds

type List interface {
	Size() int
	IsEmpty() bool
	Contains(elem interface{}) bool
	ContainsAll(elemList ...interface{}) bool
	ToArray() []interface{}
	Add(elem interface{}) bool
	AddAll(elemList ...interface{}) bool
	Remove(elem interface{}) bool
	RemoveAll(elem ...interface{}) bool
	RemoveAt(idx int) (interface{}, bool)
	SubList(from, to int) (List, error)
	Get(idx int) (interface{}, error)
	Set(idx int, elem interface{}) error
	IndexOf(elem interface{}) int
	LastIndexOf(elem interface{}) int
	Clear()
}
