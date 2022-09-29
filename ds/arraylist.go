package ds

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/SCU-SJL/sinfra/util"
)

type ArrayList struct {
	size        int
	elementData []interface{}
	mu          sync.RWMutex
	comparator  Comparator
	typeChecker TypeChecker
}

func NewArrayList(initCapacity int, typeChecker TypeChecker, comparator Comparator) *ArrayList {
	var elementData []interface{}
	if initCapacity > 0 {
		elementData = make([]interface{}, 0, initCapacity)
	} else if initCapacity == 0 {
		elementData = make([]interface{}, 0, defaultCapacity)
	} else {
		panic(fmt.Sprintf("Illegal initial capacity: %d", initCapacity))
	}
	return &ArrayList{
		size:        0,
		elementData: elementData,
		mu:          sync.RWMutex{},
		comparator:  comparator,
		typeChecker: typeChecker,
	}
}

func (list *ArrayList) Size() (size int) {
	list.mu.RLock()
	size = list.size
	list.mu.RUnlock()
	return
}

func (list *ArrayList) IsEmpty() bool {
	return list.Size() == 0
}

func (list *ArrayList) Contains(elem interface{}) bool {
	list.mu.RLock()
	defer list.mu.RUnlock()
	return list.indexOf(elem, false) >= 0
}

func (list *ArrayList) ContainsAll(elemList ...interface{}) bool {
	for _, elem := range elemList {
		if !list.Contains(elem) {
			return false
		}
	}
	return true
}

func (list *ArrayList) ToArray() (arr []interface{}) {
	list.mu.RLock()
	arr = make([]interface{}, list.size)
	copy(arr, list.elementData)
	list.mu.RUnlock()
	return
}

func (list *ArrayList) Add(elem interface{}) bool {
	if !list.isTypeValid(elem) {
		return false
	}
	list.mu.Lock()
	if list.size == cap(list.elementData) {
		list.grow(list.size + 1)
	}
	list.elementData = append(list.elementData, elem)
	list.size++
	list.mu.Unlock()
	return true
}

func (list *ArrayList) AddAll(elemList ...interface{}) bool {
	for _, elem := range elemList {
		if !list.isTypeValid(elem) {
			return false
		}
	}
	list.mu.Lock()
	afterSize := list.size + len(elemList)
	if afterSize > cap(list.elementData) {
		list.grow(afterSize)
	}
	for i, elem := range elemList {
		list.elementData[list.size+i] = elem
	}
	list.size = afterSize
	list.mu.Unlock()
	return true
}

func (list *ArrayList) Remove(elem interface{}) bool {
	list.mu.Lock()
	defer list.mu.Unlock()
	idx := list.indexOf(elem, false)
	if idx < 0 {
		return false
	}
	list.removeWithIndexList([]int{idx})
	list.size--
	return true
}

func (list *ArrayList) RemoveAll(elemList ...interface{}) bool {
	list.mu.Lock()
	defer list.mu.Unlock()
	idxSet := util.NewSet()
	for _, elem := range elemList {
		if idx := list.indexOf(elem, false); idx >= 0 {
			idxSet.Add(idx)
		}
	}
	idxList := idxSet.ToIntArray()
	list.removeWithIndexList(idxList)
	list.size -= len(elemList)
	return true
}

func (list *ArrayList) RemoveAt(idx int) (interface{}, bool) {
	list.mu.Lock()
	if !list.isIndexValid(idx) {
		return nil, false
	}
	target := list.elementData[idx]
	list.removeWithIndexList([]int{idx})
	list.size--
	list.mu.Unlock()
	return target, true
}

func (list *ArrayList) SubList(from, to int) (List, error) {
	list.mu.RLock()
	if from > to || !list.isIndexValid(to) || !list.isIndexValid(from) {
		list.mu.RUnlock()
		errMsg := fmt.Sprintf("from: %d, to: %d are out of bounds", from, to)
		return nil, errors.New(errMsg)
	}
	data := make([]interface{}, to-from)
	copy(data, list.elementData[from:to])
	list.mu.RUnlock()
	return &ArrayList{
		size:        len(data),
		elementData: data,
		mu:          sync.RWMutex{},
		comparator:  list.comparator,
		typeChecker: list.typeChecker,
	}, nil
}

func (list *ArrayList) Get(idx int) (elem interface{}, err error) {
	list.mu.RLock()
	if !list.isIndexValid(idx) {
		list.mu.RUnlock()
		return nil, newOutOfBoundsErr(idx)
	}
	elem = list.elementData[idx]
	list.mu.RUnlock()
	return
}

func (list *ArrayList) Set(idx int, elem interface{}) error {
	list.mu.Lock()
	if !list.isTypeValid(elem) {
		list.mu.Unlock()
		return typeMismatchErr
	}
	if !list.isIndexValid(idx) {
		list.mu.RUnlock()
		return newOutOfBoundsErr(idx)
	}
	list.elementData[idx] = elem
	list.mu.Unlock()
	return nil
}

func (list *ArrayList) IndexOf(elem interface{}) int {
	list.mu.RLock()
	defer list.mu.RUnlock()
	return list.indexOf(elem, false)
}

func (list *ArrayList) LastIndexOf(elem interface{}) int {
	list.mu.RLock()
	defer list.mu.RUnlock()
	return list.indexOf(elem, true)
}

func (list *ArrayList) Clear() {
	list.mu.Lock()
	list.elementData = make([]interface{}, list.size/2)
	list.size = 0
	list.mu.Unlock()
}

func (list *ArrayList) indexOf(elem interface{}, reverse bool) int {
	if list.comparator == nil {
		panic("comparator of array list is nil")
	}
	if reverse {
		for i := list.size; i >= 0; i-- {
			if list.comparator(elem, list.elementData[i]) == 0 {
				return i
			}
		}
	} else {
		for i := range list.elementData {
			if list.comparator(elem, list.elementData[i]) == 0 {
				return i
			}
		}
	}
	return -1
}

func (list *ArrayList) removeWithIndexList(idxList []int) {
	if len(idxList) == 0 {
		return
	}
	sort.Ints(idxList)
	prevIdx, curIdx := 0, 0
	startIdx := 0
	for i := range idxList {
		curIdx = idxList[i]
		copy(list.elementData[startIdx:], list.elementData[prevIdx:curIdx])
		startIdx += curIdx - prevIdx
		prevIdx = curIdx + 1
	}
	copy(list.elementData[startIdx:], list.elementData[prevIdx:])
	for i := list.size - len(idxList); i < list.size; i++ {
		list.elementData[i] = nil
	}
}

func (list *ArrayList) grow(minSize int) {
	var newCap int
	if oldCap := cap(list.elementData); oldCap < 64 {
		newCap = oldCap << 1
	} else {
		newCap = int(float64(oldCap) * 1.5)
	}
	if newCap < minSize {
		newCap = minSize
	}
	newArr := make([]interface{}, newCap)
	copy(newArr, list.elementData)
	list.elementData = newArr
}

func (list *ArrayList) isTypeValid(elem interface{}) bool {
	if list.typeChecker == nil {
		return true
	}
	return list.typeChecker(elem)
}

func (list *ArrayList) isIndexValid(idx int) bool {
	return idx >= 0 && idx < list.size
}
