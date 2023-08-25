package ds

import "sync"

type (
	MetaForest[T any] struct {
		Children []*MetaTree[T]
		Extra    any
		mu       sync.Mutex
		nodeMap  *sync.Map
	}

	MetaTree[T any] struct {
		Root  *MetaNode[T]
		Extra any
	}

	MetaNode[T any] struct {
		// This means a MetaNode is a MetaData
		MetaData[T]
		Key      string
		Val      T
		Parent   *MetaNode[T]
		Children []*MetaNode[T]
	}

	MetaData[T any] interface {
		SetKey(string)
		GetKey() string
		SetVal(T)
		GetVal() T
		GetParent() MetaData[T]
		GetChildren() []MetaData[T]
	}
)

func (f *MetaForest[T]) EnableHashIndex() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.nodeMap != nil {
		return
	}
	f.nodeMap = &sync.Map{}
}

func (f *MetaForest[T]) ResetHashIndex() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.nodeMap == nil {
		return
	}
	f.nodeMap = &sync.Map{}
}

func (f *MetaForest[T]) DisableHashIndex() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.nodeMap = nil
}

func (f *MetaForest[T]) RemoveHashIndex(key string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.nodeMap.Delete(key)
}

func (f *MetaForest[T]) SearchNode(key string) *MetaNode[T] {

	if f == nil {
		return nil
	}

	if node, ok := f.nodeMap.Load(key); ok {
		return node.(*MetaNode[T])
	}

	for _, tree := range f.Children {
		if node := tree.SearchNode(key); node != nil {
			f.nodeMap.Store(node.GetKey(), node)
			return node
		}
	}

	return nil

}

func (t *MetaTree[T]) SearchNode(key string) *MetaNode[T] {

	if t == nil || t.Root == nil {
		return nil
	}

	return t.Root.SearchNode(key)

}

func (node *MetaNode[T]) SearchNode(key string) *MetaNode[T] {

	if node == nil {
		return nil
	}

	if node.GetKey() == key {
		return node
	}

	for i := range node.Children {
		if x := node.Children[i].SearchNode(key); x != nil {
			return x
		}
	}

	return nil

}

func (node *MetaNode[T]) ExpandKey() []string {

	keyList := make([]string, 0, 8)

	node.Walk(func(mn *MetaNode[T]) {
		keyList = append(keyList, mn.GetKey())
	})

	return keyList

}

func (node *MetaNode[T]) ExpandVal() []T {

	valList := make([]T, 0, 8)

	node.Walk(func(mn *MetaNode[T]) {
		valList = append(valList, mn.GetVal())
	})

	return valList

}

func (node *MetaNode[T]) Walk(cb func(*MetaNode[T])) {

	if node == nil {
		return
	}

	cb(node)

	for i := range node.Children {
		c := node.Children[i]
		c.Walk(cb)
	}

}

func (node *MetaNode[T]) SetKey(key string) {
	node.Key = key
}

func (node *MetaNode[T]) GetKey() string {
	return node.Key
}

func (node *MetaNode[T]) SetVal(val T) {
	node.Val = val
}

func (node *MetaNode[T]) GetVal() T {
	return node.Val
}

func (node *MetaNode[T]) GetParent() MetaData[T] {
	return node.Parent
}

func (node *MetaNode[T]) GetChildren() []MetaData[T] {
	var r []MetaData[T]
	for _, c := range node.Children {
		r = append(r, c)
	}
	return r
}

func BuildMetaForestTopDown[T any](metaDataList []MetaData[T]) *MetaForest[T] {

	if len(metaDataList) == 0 {
		return nil
	}

	treeList := make([]*MetaTree[T], 0, len(metaDataList))

	for i := range metaDataList {
		metaData := metaDataList[i]
		treeList = append(treeList, BuildMetaTreeTopDown[T](metaData))
	}

	return &MetaForest[T]{Children: treeList}

}

func BuildMetaTreeTopDown[T any](metaData MetaData[T]) *MetaTree[T] {

	if metaData == nil {
		return nil
	}

	root := &MetaNode[T]{}
	toQueue := []*MetaNode[T]{root}
	fromQueue := []MetaData[T]{metaData}

	for len(fromQueue) > 0 {

		fromNode := fromQueue[0]
		fromQueue = fromQueue[1:]

		toNode := toQueue[0]
		toQueue = toQueue[1:]

		toNode.SetKey(fromNode.GetKey())
		toNode.SetVal(fromNode.GetVal())

		for i := range fromNode.GetChildren() {
			fromQueue = append(fromQueue, fromNode.GetChildren()[i])
			child := &MetaNode[T]{
				Parent: toNode,
			}
			toNode.Children = append(toNode.Children, child)
			toQueue = append(toQueue, child)
		}

	}

	return &MetaTree[T]{Root: root}

}
