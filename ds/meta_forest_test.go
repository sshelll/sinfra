package ds

import (
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

// This is a wrapper of real biz entity.
type TestUserWrapper struct {
	info     *TestUser
	id       string
	parentID *string
}

func (u *TestUserWrapper) SetKey(key string) {
	panic("implement me")
}

func (u *TestUserWrapper) GetKey() string {
	return u.id
}

func (u *TestUserWrapper) SetVal(val *TestUser) {
	panic("implement me")
}

func (u *TestUserWrapper) GetVal() *TestUser {
	return u.info
}

func (u *TestUserWrapper) GetParent() MetaData[*TestUser] {
	panic("implement me")
}

func (u *TestUserWrapper) GetChildren() []MetaData[*TestUser] {
	panic("implement me")
}

func (u *TestUserWrapper) GetParentKey() *string {
	return u.parentID
}

func (u *TestUserWrapper) GetChildrenKeys() []string {
	panic("implement me")
}

func TestRealCase(t *testing.T) {
	leaderA := &TestUserWrapper{
		id: "a",
	}
	leaderB := &TestUserWrapper{
		id: "b",
	}
	leaderC := &TestUserWrapper{
		id: "c",
	}
	workerA1 := &TestUserWrapper{
		id:       "a1",
		parentID: util.Ptr(leaderA.id),
	}
	workerA2 := &TestUserWrapper{
		id:       "a2",
		parentID: util.Ptr(leaderA.id),
	}
	workerB1 := &TestUserWrapper{
		id:       "b1",
		parentID: util.Ptr(leaderB.id),
	}
	workerB11 := &TestUserWrapper{
		id:       "b11",
		parentID: util.Ptr(workerB1.id),
	}
	users := []MetaData[*TestUser]{
		leaderA,
		leaderB,
		leaderC,
		workerA1,
		workerA2,
		workerB1,
		workerB11,
	}
	forest := BuildMetaForestBottomUp(users)
	t.Logf("forest tree cnt: %v", len(forest.Children))
	t.Log(forest.Children[0].Root.ExpandKey())
	t.Log(forest.Children[1].Root.ExpandKey())
	t.Log(forest.Children[2].Root.ExpandKey())
	t.Log(forest.SearchNode("b1").ExpandKey())
}
