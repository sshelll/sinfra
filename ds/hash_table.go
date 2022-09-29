package ds

import (
	"reflect"
	"sync"
)

// Deprecated
// use sync.Map instead
type HashTable struct {
	keyType reflect.Type
	valType reflect.Type
	mu      sync.RWMutex
	table   map[interface{}]interface{}
}

func NewHashTable(keyType, valType reflect.Type, initSize int) *HashTable {
	return &HashTable{
		keyType: keyType,
		valType: valType,
		mu:      sync.RWMutex{},
		table:   make(map[interface{}]interface{}, initSize),
	}
}

func (ht *HashTable) Set(key, val interface{}) (ok bool) {
	ht.mu.Lock()
	if !ht.isTypeValid(key, ht.keyType) || !ht.isTypeValid(val, ht.valType) {
		ht.mu.Unlock()
		return false
	}
	ht.table[key] = val
	ht.mu.Unlock()
	return true
}

func (ht *HashTable) Get(key interface{}) (val interface{}) {
	ht.mu.RLock()
	defer ht.mu.RUnlock()
	if !ht.isTypeValid(key, ht.keyType) {
		return nil
	}
	return ht.table[key]
}

func (ht *HashTable) Size() (size int) {
	ht.mu.RLock()
	size = len(ht.table)
	ht.mu.RUnlock()
	return
}

func (ht *HashTable) Keys() (keys []interface{}) {
	ht.mu.RLock()
	mapKeys := reflect.ValueOf(ht.table).MapKeys()
	for i := range mapKeys {
		keys = append(keys, mapKeys[i].Interface())
	}
	ht.mu.RUnlock()
	return
}

func (ht *HashTable) Iterator() (iter *reflect.MapIter) {
	ht.mu.RLock()
	iter = reflect.ValueOf(ht.table).MapRange()
	ht.mu.RUnlock()
	return
}

func (ht *HashTable) isTypeValid(elem interface{}, targetType reflect.Type) bool {
	return reflect.TypeOf(elem) == targetType
}
