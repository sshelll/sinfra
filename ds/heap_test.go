package ds

import (
	"github.com/stretchr/testify/suite"
	"math/rand"
	"testing"
)

type HeapTestSuite struct {
	suite.Suite
	heap *Heap
}

func TestHeap(t *testing.T) {
	suite.Run(t, new(HeapTestSuite))
}

func (suite *HeapTestSuite) SetupSuite() {
	suite.heap = NewHeap(intComparator, intTypeChecker)
}

func (suite *HeapTestSuite) AfterTest(string, string) {
	suite.heap.Clear()
}

func (suite *HeapTestSuite) TestOffer() {
	arr := rand.Perm(3)
	for i := range arr {
		suite.True(suite.heap.Offer(arr[i]))
	}
	i := suite.heap.Size() - 1
	for !suite.heap.IsEmpty() {
		suite.Equal(i, suite.heap.Poll())
		i--
	}
}

func (suite *HeapTestSuite) TestPollFromEmptyHeap() {
	suite.True(suite.heap.IsEmpty())
	suite.Nil(suite.heap.Poll())
}

func (suite *HeapTestSuite) TestRemoveAt() {
	suite.True(suite.heap.BatchOffer(5, 4, 3, 2, 1, 0))
	suite.Equal(4, suite.heap.RemoveAt(1))
	suite.Equal(5, suite.heap.RemoveAt(0))
	suite.Equal(1, suite.heap.RemoveAt(2))
	suite.Equal(3, suite.heap.Poll())
	suite.Equal(2, suite.heap.Poll())
	suite.Equal(0, suite.heap.Poll())
}

func (suite *HeapTestSuite) TestRemoveAtBorderCase() {
	suite.True(suite.heap.BatchOffer(4, 3, 2, 1, 0))
	suite.Equal(2, suite.heap.RemoveAt(2))
}

func (suite *HeapTestSuite) TestRemove() {
	suite.True(suite.heap.BatchOffer(5, 4, 3, 2, 1, 0))
	suite.True(suite.heap.Remove(5))
	suite.Equal(5, suite.heap.Size())
	i := 4
	for !suite.heap.IsEmpty() {
		suite.Equal(i, suite.heap.Poll())
		i--
	}
}

func (suite *HeapTestSuite) TestRemoveWithComparator() {
	suite.True(suite.heap.BatchOffer(10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0))

	comparator := func(a, b interface{}) int {
		aInt, bInt := a.(int), b.(int)
		if bInt <= aInt {
			return 0
		}
		return -1
	}
	suite.True(suite.heap.RemoveWithComparator(5, comparator, true))
	i := 10
	for !suite.heap.IsEmpty() {
		suite.Equal(i, suite.heap.Poll())
		i--
	}
}

func (suite *HeapTestSuite) TestTypeChecker() {
	suite.True(suite.heap.Offer(0))
	suite.False(suite.heap.Offer(1.1))
	ok, invalidIdxList := suite.heap.IgnoredBatchOffer(1, 2, 1.15)
	suite.False(ok)
	suite.Equal(1, len(invalidIdxList))
	suite.Equal(2, invalidIdxList[0])
}
