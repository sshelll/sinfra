package ds

import "sync"

type (
	MetaForest struct {
		Children []*MetaTree
		Extra    interface{}
		mu       sync.Mutex
		nodeMap  *sync.Map
	}

	MetaTree struct {
		Root  *MetaNode
		Extra interface{}
	}

	MetaNode struct {
		Key      string
		Val      interface{} // real meta data
		Children []*MetaNode
	}

	MetaData interface {
		GetKey() string
		GetVal() interface{}
		GetChildren() []MetaData
	}
)

func (f *MetaForest) EnableHashIndex() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.nodeMap != nil {
		return
	}
	f.nodeMap = &sync.Map{}
}

func (f *MetaForest) ResetHashIndex() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.nodeMap == nil {
		return
	}
	f.nodeMap = &sync.Map{}
}

func (f *MetaForest) DisableHashIndex() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.nodeMap = nil
}

func (f *MetaForest) RemoveHashIndex(key string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.nodeMap.Delete(key)
}

func (f *MetaForest) SearchNode(key string) *MetaNode {

	if f == nil {
		return nil
	}

	if node, ok := f.nodeMap.Load(key); ok {
		return node.(*MetaNode)
	}

	for _, tree := range f.Children {
		if node := tree.SearchNode(key); node != nil {
			f.nodeMap.Store(node.Key, node)
			return node
		}
	}

	return nil

}

func (t *MetaTree) SearchNode(key string) *MetaNode {

	if t == nil || t.Root == nil {
		return nil
	}

	return t.Root.SearchNode(key)

}

func (node *MetaNode) SearchNode(key string) *MetaNode {

	if node == nil {
		return nil
	}

	if node.Key == key {
		return node
	}

	for i := range node.Children {
		if x := node.Children[i].SearchNode(key); x != nil {
			return x
		}
	}

	return nil

}

func (node *MetaNode) ExpandKey() []string {

	keyList := make([]string, 0, 8)

	node.walk(func(mn *MetaNode) {
		keyList = append(keyList, mn.Key)
	})

	return keyList

}

func (node *MetaNode) ExpandVal() []interface{} {

	valList := make([]interface{}, 0, 8)

	node.walk(func(mn *MetaNode) {
		valList = append(valList, mn.Val)
	})

	return valList

}

func (node *MetaNode) walk(cb func(*MetaNode)) {

	if node == nil {
		return
	}

	cb(node)

	for i := range node.Children {
		c := node.Children[i]
		c.walk(cb)
	}

}

func BuildMetaForest(metaDataList []MetaData) *MetaForest {

	if len(metaDataList) == 0 {
		return nil
	}

	treeList := make([]*MetaTree, 0, len(metaDataList))

	for i := range metaDataList {
		metaData := metaDataList[i]
		treeList = append(treeList, BuildMetaTree(metaData))
	}

	return &MetaForest{Children: treeList}

}

func BuildMetaTree(metaData MetaData) *MetaTree {

	if metaData == nil {
		return nil
	}

	root := &MetaNode{}
	toQueue := []*MetaNode{root}
	fromQueue := []MetaData{metaData}

	for len(fromQueue) > 0 {

		fromNode := fromQueue[0]
		fromQueue = fromQueue[1:]

		toNode := toQueue[0]
		toQueue = toQueue[1:]

		toNode.Key = fromNode.GetKey()
		toNode.Val = fromNode.GetVal()

		for i := range fromNode.GetChildren() {
			fromQueue = append(fromQueue, fromNode.GetChildren()[i])
			child := &MetaNode{}
			toNode.Children = append(toNode.Children, child)
			toQueue = append(toQueue, child)
		}

	}

	return &MetaTree{Root: root}

}
