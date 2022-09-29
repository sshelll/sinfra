package ds

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TrieTestSuite struct {
	suite.Suite
}

func TestTrie(t *testing.T) {
	suite.Run(t, new(TrieTestSuite))
}

func (suite *TrieTestSuite) TestTrieContains() {
	trie := new(Trie)
	strList := []string{"hello", "world"}
	for _, s := range strList {
		trie.Add(s)
	}
	for _, s := range strList {
		suite.True(trie.Contains(s), "expected: contains '%s' == true, actual: false", s)
	}
	suite.False(trie.Contains("hell"))
	suite.False(trie.Contains("w"))
}

func (suite *TrieTestSuite) TestTrieContainsPrefix() {
	trie := new(Trie)
	strList := []string{"hello", "world"}
	for _, s := range strList {
		trie.Add(s)
	}
	prefixList := []string{"h", "w", "he", "wo", "hell", "wor"}
	for _, s := range prefixList {
		suite.False(trie.Contains(s))
		suite.True(trie.ContainsPrefix(s))
	}
}

func (suite *TrieTestSuite) TestTrieDelete() {
	trie := new(Trie)
	strList := []string{"1", "12", "123", "1234"}
	for _, s := range strList {
		trie.Add(s)
	}
	for _, s := range strList {
		suite.True(trie.Contains(s), s)
		suite.True(trie.ContainsPrefix(s))
	}

	trie.Delete("1")
	suite.True(trie.ContainsPrefix("1"))
	suite.False(trie.Contains("1"))
	for _, s := range []string{"12", "123", "1234"} {
		suite.True(trie.Contains(s))
	}

	trie.Delete("12")
	suite.True(trie.ContainsPrefix("1"))
	suite.True(trie.ContainsPrefix("12"))
	suite.False(trie.Contains("1"))
	suite.False(trie.Contains("12"))
	for _, s := range []string{"123", "1234"} {
		suite.True(trie.Contains(s))
	}

	trie.Delete("123")
	suite.True(trie.ContainsPrefix("1"))
	suite.True(trie.ContainsPrefix("12"))
	suite.True(trie.ContainsPrefix("123"))
	suite.True(trie.Contains("1234"))
	suite.False(trie.Contains("1"))
	suite.False(trie.Contains("12"))
	suite.False(trie.Contains("123"))
}

func (suite *TrieTestSuite) TestSearch() {
	trie := new(Trie)
	strList := []string{"1", "12", "123", "1234", "你", "你我", "你我他"}
	for _, s := range strList {
		trie.Add(s)
	}
	suite.Equal([]string{"1", "12", "123", "1234"}, trie.Search("1"))
	suite.Equal([]string{"12", "123", "1234"}, trie.Search("12"))
	suite.Equal([]string{"123", "1234"}, trie.Search("123"))
	suite.Equal([]string{"1234"}, trie.Search("1234"))
	suite.Equal(0, len(trie.Search("0")))
	suite.Equal([]string{"你", "你我", "你我他"}, trie.Search("你"))
	suite.Equal([]string{"你我", "你我他"}, trie.Search("你我"))
	suite.Equal([]string{"你我他"}, trie.Search("你我他"))
}
