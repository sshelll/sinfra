package ds

/* Heap a thread-unsafe heap with customized comparator and type checker.
 * attention: nil element is not allowed.
 */
type Heap struct {
	data        []interface{}
	comparator  Comparator
	typeChecker TypeChecker
	size        int
}

func NewHeap(comparator Comparator, typeChecker TypeChecker) *Heap {
	return &Heap{
		data:        make([]interface{}, 1),
		comparator:  comparator,
		size:        0,
		typeChecker: typeChecker,
	}
}

func (heap *Heap) IsEmpty() bool {
	return heap.size == 0
}

func (heap *Heap) Size() int {
	return heap.size
}

func (heap *Heap) Peek() interface{} {
	if heap.IsEmpty() {
		return nil
	}
	return heap.data[0]
}

func (heap *Heap) IgnoredBatchOffer(elemList ...interface{}) (ok bool, invalidIdxList []int) {
	return heap.IgnoredBatchOfferWithArray(elemList[:])
}

/* IgnoredBatchOfferWithArray offer elem list and allow failure */
func (heap *Heap) IgnoredBatchOfferWithArray(elemArr []interface{}) (ok bool, invalidIdxList []int) {
	for i := range elemArr {
		if !heap.Offer(elemArr[i]) {
			invalidIdxList = append(invalidIdxList, i)
		}
	}
	return len(invalidIdxList) == 0, invalidIdxList
}

func (heap *Heap) BatchOffer(elemList ...interface{}) (ok bool) {
	return heap.BatchOfferWithArray(elemList[:])
}

func (heap *Heap) BatchOfferWithArray(elemArr []interface{}) (ok bool) {
	if len(elemArr) == 0 {
		return false
	}

	for i := range elemArr {
		if !heap.IsTypeValid(elemArr[i]) {
			return false
		}
	}
	for i := range elemArr {
		heap.Offer(elemArr[i])
	}
	return true
}

func (heap *Heap) Offer(elem interface{}) (ok bool) {
	if elem == nil {
		return false
	}
	if !heap.IsTypeValid(elem) {
		return false
	}
	if heap.size >= len(heap.data) {
		heap.grow()
	}
	if heap.size == 0 {
		heap.data[0] = elem
	} else {
		heap.siftUp(heap.size, elem)
	}
	heap.size++
	return true
}

func (heap *Heap) Poll() interface{} {
	if heap.size == 0 {
		return nil
	}
	lastIdx := heap.size - 1
	result := heap.data[0]
	lastElem := heap.data[lastIdx]
	heap.data[lastIdx] = nil
	heap.size--
	if lastIdx > 0 {
		heap.siftDown(0, lastElem)
	}
	return result
}

func (heap *Heap) Remove(target interface{}) (ok bool) {
	return heap.RemoveWithComparator(target, nil, false)
}

/* RemoveWithComparator remove elements with customized comparator.
 * 'target' is the first parameter of the comparator.
 */
func (heap *Heap) RemoveWithComparator(target interface{}, comparator Comparator, rmAll bool) (ok bool) {
	if comparator == nil {
		comparator = heap.comparator
	}

	if rmAll {
		/* remove one by one until no target exists */
		for heap.RemoveWithComparator(target, comparator, false) {
			ok = true
		}
		return
	}

	for i := range heap.data {
		if comparator(target, heap.data[i]) == 0 {
			return heap.RemoveAt(i) != nil
		}
	}
	return false
}

func (heap *Heap) RemoveAt(idx int) interface{} {
	if idx >= heap.size || idx < 0 {
		return nil
	}
	lastIdx := heap.size - 1
	result := heap.data[idx]
	lastElem := heap.data[lastIdx]
	heap.data[lastIdx] = nil
	heap.size--
	if lastIdx > 0 {
		heap.siftDown(idx, lastElem)
	}
	return result
}

func (heap *Heap) Clear() {
	newCapacity := heap.size / 2
	if newCapacity == 0 {
		newCapacity = 1
	}
	heap.data = make([]interface{}, newCapacity, newCapacity)
	heap.size = 0
}

func (heap *Heap) IsTypeValid(elem interface{}) bool {
	if heap.typeChecker == nil {
		return true
	}
	return heap.typeChecker(elem)
}

func (heap *Heap) Clone() *Heap {
	newData := make([]interface{}, len(heap.data))
	copy(newData, heap.data)
	return &Heap{
		data:        newData,
		comparator:  heap.comparator,
		typeChecker: heap.typeChecker,
		size:        heap.size,
	}
}

/* siftUp adjust heap from bottom to top
 * 'idx' is the index of the element to be inserted
 * 'elem' is the element to be inserted
 */
func (heap *Heap) siftUp(idx int, elem interface{}) {
	for idx > 0 {
		parentIdx := (idx - 1) >> 1
		if heap.compare(heap.data[parentIdx], elem) {
			break
		}
		heap.data[idx] = heap.data[parentIdx]
		idx = parentIdx
	}
	heap.data[idx] = elem
}

/* siftDown adjust heap from top to bottom
 * 'idx' is the index of the element to be removed
 * 'elem' is the element to be removed
 */
func (heap *Heap) siftDown(idx int, elem interface{}) {
	halfIdx := heap.size >> 1
	for idx < halfIdx {
		leftIdx := (idx << 1) + 1
		rightIdx := leftIdx + 1
		targetIdx := leftIdx

		if rightIdx < heap.size && heap.compare(heap.data[rightIdx], heap.data[leftIdx]) {
			targetIdx = rightIdx
		}

		if heap.compare(elem, heap.data[targetIdx]) {
			break
		}
		heap.data[idx] = heap.data[targetIdx]
		idx = targetIdx
	}
	heap.data[idx] = elem
}

func (heap *Heap) compare(a, b interface{}) (priorityHigher bool) {
	return heap.comparator(a, b) >= 0
}

func (heap *Heap) grow() {
	var newCapacity int
	if heap.size < 64 {
		newCapacity = heap.size * 2
	} else {
		newCapacity = int(float64(heap.size) * 1.5)
	}
	newData := make([]interface{}, newCapacity)
	copy(newData, heap.data)
	heap.data = newData
}

func (heap *Heap) reHeapify() {
	newHeap := NewHeap(heap.comparator, heap.typeChecker)
	for _, elem := range heap.data {
		newHeap.Offer(elem)
	}
	heap.data = newHeap.data
}
