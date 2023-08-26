package ds

import "sync"

type (
	MetaForest[T any] struct {
		Children []*MetaTree[T]
		Extra    any
		nodeMap  *sync.Map
	}

	MetaTree[T any] struct {
		Root  *MetaNode[T]
		Extra any
	}

	MetaNode[T any] struct {
		// This means a MetaNode is a MetaData
		MetaData[T]
		Key       string
		Val       T
		ParentKey *string
		Parent    *MetaNode[T]
		Children  []*MetaNode[T]
	}

	MetaData[T any] interface {
		SetKey(string)
		GetKey() string
		SetVal(T)
		GetVal() T
		GetParent() MetaData[T]
		GetChildren() []MetaData[T]
		GetParentKey() *string
		GetChildrenKeys() []string
	}
)

// EnableHashIndex enables hash index for MetaForest.
// This is useful when you want to search a node by key.
// NOTE: this function just init the index, it won't build the index,
// the index will be built when you call SearchNode.
func (f *MetaForest[T]) EnableHashIndex() {
	if f.nodeMap != nil {
		return
	}
	f.nodeMap = &sync.Map{}
}

func (f *MetaForest[T]) ResetHashIndex() {
	if f.nodeMap == nil {
		return
	}
	f.nodeMap = &sync.Map{}
}

func (f *MetaForest[T]) DisableHashIndex() {
	f.nodeMap = nil
}

func (f *MetaForest[T]) RemoveHashIndex(key string) {
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

func (node *MetaNode[T]) GetParentKey() *string {
	var t *MetaNode[T]
	if node.Parent != t {
		k := node.Parent.GetKey()
		return &k
	}
	return node.ParentKey
}

func (node *MetaNode[T]) GetChildrenKeys() []string {
	var r []string
	for _, c := range node.Children {
		r = append(r, c.GetKey())
	}
	return r
}

// ConvToMetaForest builds a MetaForest from a list of MetaData.
// This function just conv a list of MetaData to a list of MetaTree,
// so you can use a list of functions provided by MetaTree and MetaForest.
func ConvToMetaForest[T any](metaDataList []MetaData[T]) *MetaForest[T] {
	if len(metaDataList) == 0 {
		return nil
	}

	treeList := make([]*MetaTree[T], 0, len(metaDataList))

	for i := range metaDataList {
		metaData := metaDataList[i]
		treeList = append(treeList, ConvToMetaTree[T](metaData))
	}

	return &MetaForest[T]{Children: treeList}
}

func ConvToMetaTree[T any](metaData MetaData[T]) *MetaTree[T] {
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

func isNil[T comparable](arg T) bool {
	var t T
	return arg == t
}

func ConvToMetaDataList[T any](metaNodes []*MetaNode[T]) []MetaData[T] {
	var r []MetaData[T]
	for _, node := range metaNodes {
		r = append(r, node)
	}
	return r
}

// BuildMetaForestBottomUp builds a MetaForest from a list of MetaData.
// This function builds the forest bottom up, so each MetaData must have a 'parent key',
// which means you have to impl GetParentKey(), GetKey(), GetVal() methods.
// NOTE: GetParent() returns "" means this node is a root.
func BuildMetaForestBottomUp[T any](metaDataList []MetaData[T]) *MetaForest[T] {
	var roots []*MetaNode[T]
	var nodeMap = &sync.Map{}

	// init node map
	for _, metaData := range metaDataList {
		key := metaData.GetKey()
		nodeMap.Store(key, &MetaNode[T]{
			Key:       metaData.GetKey(),
			Val:       metaData.GetVal(),
			ParentKey: metaData.GetParentKey(),
		})
	}

	// link nodes to their parents
	nodeMap.Range(func(key, value any) (continue_ bool) {
		continue_ = true
		node := value.(*MetaNode[T])

		// find parent
		parentKey := node.GetParentKey()
		if parentKey == nil {
			roots = append(roots, node)
			return
		}
		p, ok := nodeMap.Load(*parentKey)
		parent := p.(*MetaNode[T])

		// no parent found, this node is a root
		if !ok {
			roots = append(roots, node)
			return
		}

		// do link
		node.Parent = parent
		parent.Children = append(parent.Children, node)
		return
	})

	trees := make([]*MetaTree[T], 0, len(roots))
	for _, root := range roots {
		trees = append(trees, &MetaTree[T]{Root: root})
	}

	return &MetaForest[T]{
		Children: trees,
		nodeMap:  nodeMap,
	}
}
