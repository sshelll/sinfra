package ds

import (
	"fmt"
	"sync"
)

type node struct {
	data interface{}
	next *node
	prev *node
}

/* LinkedList a thread-safe linked list */
type LinkedList struct {
	size        int
	head        *node
	tail        *node
	mu          sync.RWMutex
	comparator  Comparator
	typeChecker TypeChecker
}

func NewLinkedList(comparator Comparator, typeChecker TypeChecker) *LinkedList {
	return &LinkedList{
		mu:          sync.RWMutex{},
		comparator:  comparator,
		typeChecker: typeChecker,
	}
}

func (list *LinkedList) IsEmpty() bool {
	return list.Size() == 0
}

func (list *LinkedList) Size() int {
	list.mu.RLock()
	defer list.mu.RUnlock()
	return list.size
}

func (list *LinkedList) Get(idx int) (interface{}, error) {
	list.mu.RLock()
	defer list.mu.RUnlock()
	if err := list.checkIndex(idx); err != nil {
		return nil, err
	}
	res := list.nodeOf(idx)
	return res.data, nil
}

/* Set replaces the element at the specified position */
func (list *LinkedList) Set(idx int, data interface{}) error {
	list.mu.Lock()
	defer list.mu.Unlock()
	if err := list.checkIndex(idx); err != nil {
		return err
	}
	res := list.nodeOf(idx)
	res.data = data
	return nil
}

func (list *LinkedList) GetFirst() (interface{}, error) {
	list.mu.RLock()
	defer list.mu.RUnlock()
	if err := list.checkIndex(0); err != nil {
		return nil, err
	}
	return list.head.data, nil
}

func (list *LinkedList) GetLast() (interface{}, error) {
	list.mu.RLock()
	defer list.mu.RUnlock()
	if err := list.checkIndex(list.size - 1); err != nil {
		return nil, err
	}
	return list.tail.data, nil
}

func (list *LinkedList) Add(interface{}) bool {
	panic("do not call this method")
}

func (list *LinkedList) AddAll(...interface{}) bool {
	panic("do not call this method")
}

func (list *LinkedList) AddFirst(data interface{}) bool {
	if !list.isTypeValid(data) {
		return false
	}
	cur := &node{
		data: data,
		prev: nil,
	}
	list.mu.Lock()
	defer list.mu.Unlock()
	if list.head != nil {
		cur.next = list.head
		list.head.prev = cur
	} else {
		list.tail = cur
	}
	list.head = cur
	list.size++
	return true
}

func (list *LinkedList) AddLast(data interface{}) bool {
	if !list.isTypeValid(data) {
		return false
	}
	cur := &node{
		data: data,
		next: nil,
	}
	list.mu.Lock()
	defer list.mu.Unlock()
	if list.tail != nil {
		cur.prev = list.tail
		list.tail.next = cur
	} else {
		list.head = cur
	}
	list.tail = cur
	list.size++
	return true
}

func (list *LinkedList) Remove(elem interface{}) (ok bool) {
	list.mu.Lock()
	_, ok = list.removeAt(list.indexOf(elem, false))
	list.mu.Unlock()
	return
}

func (list *LinkedList) RemoveAt(idx int) (elem interface{}, ok bool) {
	list.mu.Lock()
	elem, ok = list.removeAt(idx)
	list.mu.Unlock()
	return
}

func (list *LinkedList) RemoveAll(elemList ...interface{}) (ok bool) {
	ok = false
	for _, elem := range elemList {
		ok = ok || list.Remove(elem)
	}
	return
}

func (list *LinkedList) RemoveFirst() (interface{}, error) {
	list.mu.Lock()
	defer list.mu.Unlock()
	if err := list.checkIndex(0); err != nil {
		return nil, err
	}

	res := list.head
	if list.head == list.tail {
		list.tail = nil
	}
	list.head = res.next
	if list.head != nil {
		list.head.prev = nil
	}
	list.size--
	return res.data, nil
}

func (list *LinkedList) RemoveLast() (interface{}, error) {
	list.mu.Lock()
	defer list.mu.Unlock()
	if err := list.checkIndex(list.size - 1); err != nil {
		return nil, err
	}
	res := list.tail
	if list.tail == list.head {
		list.head = nil
	}
	list.tail = res.prev
	if list.tail != nil {
		list.tail.next = nil
	}
	list.size--
	return res.data, nil
}

func (list *LinkedList) Offer(data interface{}) {
	list.AddLast(data)
}

func (list *LinkedList) Poll() (interface{}, error) {
	return list.RemoveFirst()
}

func (list *LinkedList) Push(data interface{}) {
	list.AddLast(data)
}

func (list *LinkedList) Pop() (interface{}, error) {
	return list.RemoveLast()
}

func (list *LinkedList) Contains(elem interface{}) bool {
	list.mu.RLock()
	idx := list.indexOf(elem, false)
	list.mu.RUnlock()
	return idx >= 0
}

func (list *LinkedList) ContainsAll(elemList ...interface{}) bool {
	for _, elem := range elemList {
		if !list.Contains(elem) {
			return false
		}
	}
	return true
}

/* SubList returns List instead of LinkedList for implementing interface */
func (list *LinkedList) SubList(from, to int) (List, error) {
	if from > to {
		return nil, fmt.Errorf("end idx: %d should not be smaller than start idx: %d", to, from)
	}
	list.mu.RLock()
	if err := list.checkIndex(from); err != nil {
		return nil, err
	}
	if err := list.checkIndex(to); err != nil {
		return nil, err
	}
	subList := NewLinkedList(list.comparator, list.typeChecker)
	headNode := list.nodeOf(from)
	tailNode := list.nodeOf(to)
	for headNode != tailNode {
		subList.AddLast(headNode.data)
		headNode = headNode.next
	}
	subList.AddLast(tailNode.data)
	list.mu.RUnlock()
	return subList, nil
}

func (list *LinkedList) IndexOf(elem interface{}) (idx int) {
	list.mu.RLock()
	idx = list.indexOf(elem, false)
	list.mu.RUnlock()
	return
}

func (list *LinkedList) LastIndexOf(elem interface{}) (idx int) {
	list.mu.RLock()
	idx = list.indexOf(elem, true)
	list.mu.RUnlock()
	return
}

func (list *LinkedList) ToArray() (arr []interface{}) {
	arr = make([]interface{}, list.Size())
	for node := list.head; node != nil; node = node.next {
		arr = append(arr, node.data)
	}
	return arr
}

func (list *LinkedList) Clear() {
	list.mu.Lock()
	list.clear()
	list.mu.Unlock()
}

func (list *LinkedList) nodeOf(idx int) *node {
	if err := list.checkIndex(idx); err != nil {
		return nil
	}
	var res *node
	if idx < list.size>>1 {
		res = list.head
		for i := 0; i < idx; i++ {
			res = res.next
		}
	} else {
		res = list.tail
		for i := list.size - 1; i > idx; i-- {
			res = res.prev
		}
	}
	return res
}

func (list *LinkedList) indexOf(elem interface{}, reverse bool) (idx int) {
	if list.comparator == nil {
		panic("comparator of linked list is nil")
	}
	if reverse {
		idx = list.size - 1
		for node := list.tail; node != nil; node = node.prev {
			if list.comparator(elem, node) == 0 {
				return idx
			}
			idx--
		}
	} else {
		idx = 0
		for node := list.head; node != nil; node = node.next {
			if list.comparator(elem, node.data) == 0 {
				return idx
			}
			idx++
		}
	}
	return -1
}

func (list *LinkedList) clear() {
	list.size = 0
	list.head = nil
	list.tail = nil
}

func (list *LinkedList) removeAt(idx int) (interface{}, bool) {
	node := list.nodeOf(idx)
	if node == nil {
		return nil, false
	} else if list.size == 1 {
		list.clear()
		return node.data, true
	} else if idx == 0 {
		list.head = node.next
		list.head.prev = nil
	} else if idx == list.size-1 {
		list.tail = node.prev
		list.tail.next = nil
	} else {
		node.prev.next = node.next
		node.next.prev = node.prev
	}
	list.size--
	return node.data, true
}

func (list *LinkedList) checkIndex(n int) error {
	if n > list.size-1 || n < 0 {
		return fmt.Errorf("index %d out of bounds for length %d", n, list.size)
	}
	return nil
}

func (list *LinkedList) isTypeValid(elem interface{}) bool {
	if list.typeChecker == nil {
		return true
	}
	return list.typeChecker(elem)
}
