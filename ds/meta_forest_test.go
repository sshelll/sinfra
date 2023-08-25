package ds

import (
	"testing"
)

type UnitTestNode struct {
	k string
	v string
	p *UnitTestNode
	c []*UnitTestNode
}

func (n *UnitTestNode) SetKey(k string) {
	n.k = k
}

func (n *UnitTestNode) GetKey() string {
	return n.k
}

func (n *UnitTestNode) SetVal(v string) {
	n.v = v
}

func (n *UnitTestNode) GetVal() string {
	return n.v
}

func (n *UnitTestNode) GetParent() MetaData[string] {
	return n.p
}

func (n *UnitTestNode) GetChildren() []MetaData[string] {
	var r []MetaData[string]
	for _, c := range n.c {
		r = append(r, c)
	}
	return r
}

func TestBuildTree(t *testing.T) {
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
	utRoot := &UnitTestNode{
		k: "root",
		v: "root-val",
		c: []*UnitTestNode{
			{
				k: "child-1",
				v: "child-1-val",
				c: []*UnitTestNode{
					{
						k: "child-1-1",
						v: "child-1-1-val",
						c: []*UnitTestNode{
							{
								k: "child-1-1-1",
								v: "child-1-1-1-val",
							},
							{
								k: "child-1-1-2",
								v: "child-1-1-2-val",
							},
						},
					},
					{
						k: "child-1-2",
						v: "child-1-2-val",
						c: []*UnitTestNode{
							{
								k: "child-1-2-1",
								v: "child-1-2-1-val",
							},
						},
					},
				},
			},
			{
				k: "child-2",
				v: "child-2-val",
			},
			{
				k: "child-3",
				v: "child-3-val",
			},
		},
	}
	tree := BuildMetaTreeTopDown(utRoot)
	keys := tree.Root.ExpandKey()
	vals := tree.Root.ExpandVal()
	t.Logf("keys: %v", keys)
	t.Logf("vals: %v", vals)
}
