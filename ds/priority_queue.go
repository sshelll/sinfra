package ds

import (
	"sync"
)

type PriorityQueue struct {
	queue *Heap
	mu    sync.RWMutex
}

func NewPriorityQueue(comparator Comparator, typeChecker TypeChecker) *PriorityQueue {
	return &PriorityQueue{
		queue: NewHeap(comparator, typeChecker),
		mu:    sync.RWMutex{},
	}
}

func (pq *PriorityQueue) Offer(elem interface{}) (ok bool) {
	pq.mu.Lock()
	ok = pq.queue.Offer(elem)
	pq.mu.Unlock()
	return
}

func (pq *PriorityQueue) BatchOffer(elemList ...interface{}) (ok bool) {
	pq.mu.Lock()
	ok = pq.queue.BatchOffer(elemList)
	pq.mu.Unlock()
	return
}

func (pq *PriorityQueue) BatchOfferWithArray(elemList []interface{}) (ok bool) {
	pq.mu.Lock()
	ok = pq.queue.BatchOfferWithArray(elemList)
	pq.mu.Unlock()
	return
}

func (pq *PriorityQueue) Poll() (elem interface{}) {
	pq.mu.Lock()
	elem = pq.queue.Poll()
	pq.mu.Unlock()
	return
}

func (pq *PriorityQueue) Remove(target interface{}) (ok bool) {
	pq.mu.Lock()
	ok = pq.queue.Remove(target)
	pq.mu.Unlock()
	return
}

/* RemoveWithComparator remove all targets, this method could cost much time and cause blocking */
func (pq *PriorityQueue) RemoveWithComparator(target interface{}, comparator Comparator, rmAll bool) (ok bool) {
	pq.mu.Lock()
	ok = pq.queue.RemoveWithComparator(target, comparator, rmAll)
	pq.mu.Unlock()
	return
}

func (pq *PriorityQueue) Size() (size int) {
	pq.mu.RLock()
	size = pq.queue.Size()
	pq.mu.RUnlock()
	return
}

func (pq *PriorityQueue) IsEmpty() (isEmpty bool) {
	pq.mu.RLock()
	isEmpty = pq.queue.IsEmpty()
	pq.mu.RUnlock()
	return
}

func (pq *PriorityQueue) Clear() {
	pq.mu.Lock()
	pq.queue.Clear()
	pq.mu.Unlock()
}

func (pq *PriorityQueue) Peek() (head interface{}) {
	pq.mu.RLock()
	head = pq.queue.Peek()
	pq.mu.RUnlock()
	return
}

func (pq *PriorityQueue) IsElemTypeValid(elem interface{}) bool {
	return pq.queue.IsTypeValid(elem)
}

func (pq *PriorityQueue) Clone() *PriorityQueue {
	return &PriorityQueue{
		queue: pq.queue.Clone(),
		mu:    sync.RWMutex{},
	}
}