package ds

import "github.com/sshelll/sinfra/util"

type Trie struct {
	root *trieNode
}

type trieNode struct {
	next   map[rune]*trieNode
	cnt    int
	endCnt int
}

func newTrieNode() *trieNode {
	return &trieNode{
		next:   make(map[rune]*trieNode),
		cnt:    0,
		endCnt: 0,
	}
}

func (tn *trieNode) traversal() []string {
	res := make([]string, 0, 4)
	buf := util.NewRuneBuffer()

	var dfs func(*trieNode)
	dfs = func(node *trieNode) {
		if node.endCnt > 0 {
			res = append(res, buf.String())
		}
		for r, n := range node.next {
			if n != nil {
				buf.WriteRune(r)
				dfs(n)
				buf.RemoveLastRune()
			}
		}
	}
	dfs(tn)

	return res
}

func (t *Trie) Add(s string) {
	if len(s) == 0 {
		return
	}
	if t.root == nil {
		t.root = newTrieNode()
	}
	node := t.root
	for _, r := range s {
		tn := node.next[r]
		if tn == nil {
			tn = newTrieNode()
		}
		tn.cnt++
		node.next[r] = tn
		node = tn
	}
	node.endCnt++
}

func (t *Trie) Contains(s string) bool {
	tn := t.searchStr(s)
	return tn != nil && tn.endCnt > 0
}

func (t *Trie) ContainsPrefix(s string) bool {
	return t.searchStr(s) != nil
}

func (t *Trie) Search(s string) []string {
	tn := t.searchStr(s)
	if tn == nil {
		return nil
	}
	res := tn.traversal()
	for i := range res {
		res[i] = s + res[i]
	}
	return res
}

func (t *Trie) Delete(s string) {
	if !t.Contains(s) {
		return
	}
	node := t.root
	for _, r := range s {
		nn := node.next[r]
		nn.cnt--
		if nn.cnt == 0 {
			delete(node.next, r)
			break
		}
		node = nn
	}
	node.endCnt--
}

func (t *Trie) searchStr(s string) *trieNode {
	if len(s) == 0 {
		return nil
	}
	if t.root == nil {
		return nil
	}
	node := t.root
	for _, r := range s {
		nn := node.next[r]
		if nn == nil {
			return nil
		}
		node = nn
	}
	return node
}
