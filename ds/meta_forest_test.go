package ds

import (
	"reflect"
	"testing"

	"github.com/sshelll/sinfra/util"
)

func TestConvToTree(t *testing.T) {
	/*
		root
		├── child-1
		│   ├── child-1-1
		│   │   ├── child-1-1-1
		│   │   └── child-1-1-2
		│   └── child-1-2
		│       └── child-1-2-1
		├── child-2
		└── child-3
	*/
	utRoot := &MetaNode[string]{
		Key: "root",
		Val: "root-val",
		Children: []*MetaNode[string]{
			{
				Key: "child-1",
				Val: "child-1-val",
				Children: []*MetaNode[string]{
					{
						Key: "child-1-1",
						Val: "child-1-1-val",
						Children: []*MetaNode[string]{
							{
								Key: "child-1-1-1",
								Val: "child-1-1-1-val",
							},
							{
								Key: "child-1-1-2",
								Val: "child-1-1-2-val",
							},
						},
					},
					{
						Key: "child-1-2",
						Val: "child-1-2-val",
						Children: []*MetaNode[string]{
							{
								Key: "child-1-2-1",
								Val: "child-1-2-1-val",
							},
						},
					},
				},
			},
			{
				Key: "child-2",
				Val: "child-2-val",
			},
			{
				Key: "child-3",
				Val: "child-3-val",
			},
		},
	}
	tree := ConvToMetaTree(utRoot)
	keys := tree.Root.ExpandKey()
	vals := tree.Root.ExpandVal()
	t.Logf("keys: %v", keys)
	t.Logf("vals: %v", vals)
}

// This refers to real biz entity.
type TestUser struct {
	name string
	age  int
}

func TestRealCase(t *testing.T) {
	// wrap biz entity to MetaNode
	leaderA := NewMetaNode[*TestUser]("a", &TestUser{name: "a"}, nil)
	leaderB := NewMetaNode[*TestUser]("b", &TestUser{name: "b"}, nil)
	leaderC := NewMetaNode[*TestUser]("c", &TestUser{name: "c"}, nil)
	leaderD := NewMetaNode[*TestUser]("d", &TestUser{name: "d"}, nil)
	workerA1 := NewMetaNode[*TestUser]("a1", &TestUser{name: "a1"}, &leaderA.Key)
	workerA2 := NewMetaNode[*TestUser]("a2", &TestUser{name: "a2"}, &leaderA.Key)
	workerB1 := NewMetaNode[*TestUser]("b1", &TestUser{name: "b1"}, &leaderB.Key)
	workerB11 := NewMetaNode[*TestUser]("b11", &TestUser{name: "b11"}, &workerB1.Key)
	workerC1 := NewMetaNode[*TestUser]("c1", &TestUser{name: "c1"}, &leaderC.Key)

	users := []MetaData[*TestUser]{
		leaderA,
		leaderB,
		leaderC,
		leaderD,
		workerA1,
		workerA2,
		workerB1,
		workerB11,
		workerC1,
	}

	forest := BuildMetaForestBottomUp(users)

	t.Logf("forest tree cnt: %v", len(forest.Children))
	t.Log(forest.Children[0].Root.ExpandKey())
	t.Log(forest.Children[1].Root.ExpandKey())
	t.Log(forest.Children[2].Root.ExpandKey())
	t.Log(forest.Children[3].Root.ExpandKey())
	t.Log(forest.SearchNode("b1").ExpandKey())

	t.Log("-----------------")

	var empty MetaData[*TestUser]
	var emptyNode *MetaNode[*TestUser]
	t.Log(leaderA.GetParent(),
		leaderA.Parent == nil,
		leaderA.NoParent(),
		leaderA.GetParent() == nil,
		leaderA.GetParent() == empty,
		leaderA.GetParent() == emptyNode,
		leaderA.Nil(leaderA.GetParent()),
	)

	var dataA MetaData[*TestUser] = leaderA
	t.Log(
		dataA.GetParent() == nil,                     // false
		util.IsNilT(dataA.GetParent()),               // true
		reflect.ValueOf(leaderA.GetParent()).IsNil(), // true
		reflect.ValueOf(dataA.GetParent()).IsNil(),   // true
	)
}
