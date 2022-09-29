package ds

/* linkedEntry 链式哈希表节点
 * 在哈希表中形成一个链表来维护节点顺序
 */
type linkedEntry struct {
	key     string
	prev    *linkedEntry
	next    *linkedEntry
	payload interface{}
}

func NewLinkedEntry(key string, payload interface{}) *linkedEntry {
	return &linkedEntry{
		key:     key,
		payload: payload,
	}
}

/* AddBefore 将当前节点插入到目标节点之前
 * Before: A <-> B
 * Call:   C.AddBefore(B)
 * After:  A <-> C <-> B
 */
func (e *linkedEntry) AddBefore(entry *linkedEntry) {
	ePrev := entry.prev
	if ePrev != nil {
		ePrev.next = e
	}
	e.prev = ePrev
	e.next = entry
	entry.prev = e
}
