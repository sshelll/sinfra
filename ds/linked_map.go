package ds

import (
	"sync"
)

/* LinkedMap 并发安全的 LRU 链式哈希表 */
type LinkedMap struct {
	entries  map[string]*linkedEntry
	head     *linkedEntry
	tail     *linkedEntry
	capacity int
	size     int
	mapMu    sync.RWMutex
	listMu   sync.Mutex
}

func NewLinkedMap(capacity int) *LinkedMap {
	return &LinkedMap{
		// set 'cap = cap + 1' to avoid resizing map
		entries:  make(map[string]*linkedEntry, capacity+1),
		head:     nil,
		tail:     nil,
		capacity: capacity,
		size:     0,
		mapMu:    sync.RWMutex{},
		listMu:   sync.Mutex{},
	}
}

func (lm *LinkedMap) Set(key string, val interface{}) {
	lm.mapMu.Lock()
	// key exists
	if entry := lm.entryOf(key); entry != nil {
		entry.payload = val
		lm.afterEntryAccess(entry)
		lm.mapMu.Unlock()
		return
	}

	// add new entry, possibly trigger lru
	newEntry := NewLinkedEntry(key, val)
	lm.addNewEntry(key, newEntry)
	lm.mapMu.Unlock()
	return
}

func (lm *LinkedMap) Get(key string) (result interface{}) {
	lm.mapMu.RLock()
	if entry, ok := lm.entries[key]; ok {
		result = entry.payload
		lm.afterEntryAccess(entry)
	}
	lm.mapMu.RUnlock()
	return
}

func (lm *LinkedMap) Size() (size int) {
	lm.mapMu.RLock()
	size = lm.size
	lm.mapMu.RUnlock()
	return
}

/* move target entry to last */
func (lm *LinkedMap) afterEntryAccess(entry *linkedEntry) {
	lm.listMu.Lock()
	// do nothing when cur entry is tail
	if lm.tail != entry {
		// entry.next mustn't be nil
		if entry.prev == nil { // entry is head
			entry.next.prev = nil
			lm.head = entry.next
		} else { // entry is not head
			entry.prev.next = entry.next
			entry.next.prev = entry.prev
		}
		entry.next = nil
		entry.prev = lm.tail
		lm.tail.next = entry
		lm.tail = entry
	}
	lm.listMu.Unlock()
}

/* possibly remove eldest entry before add */
func (lm *LinkedMap) addNewEntry(key string, entry *linkedEntry) {
	if lm.size > lm.capacity {
		panic("size of linked map is greater than its capacity")
	}
	if lm.size == lm.capacity {
		lm.removeEldestEntry()
	}
	lm.entries[key] = entry
	if lm.head != nil {
		entry.AddBefore(lm.head)
	} else { // first entry
		lm.tail = entry
	}
	lm.head = entry
	lm.size++
	return
}

func (lm *LinkedMap) removeEldestEntry() {
	first := lm.head
	if first != nil {
		lm.head = first.next
		lm.head.prev = nil
		delete(lm.entries, first.key)
		lm.size--
	}
}

func (lm *LinkedMap) entryOf(key string) *linkedEntry {
	entry := lm.entries[key]
	return entry
}
