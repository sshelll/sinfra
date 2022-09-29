package ds

import (
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
)

type LinkedMapTestSuite struct {
	suite.Suite
}


func TestLinkedMap(t *testing.T) {
	suite.Run(t, new(LinkedMapTestSuite))
}

func (suite *LinkedMapTestSuite) TestWorkInSerial() {
	linkedMap := NewLinkedMap(3)
	linkedMap.Set("1", 1) // 1 -> nil
	linkedMap.Set("2", 2) // 2 -> 1
	linkedMap.Set("3", 3) // 3 -> 2 -> 1
	suite.Equal(3, linkedMap.Size())

	linkedMap.Set("4", 4) // 4 -> 2 -> 1
	suite.Nil(linkedMap.Get("3"))
	suite.Equal(3, linkedMap.Size())

	linkedMap.Set("5", 5) // 5 -> 2 -> 1
	suite.Nil(linkedMap.Get("4"))
	suite.Equal(3, linkedMap.Size())

	_ = linkedMap.Get("5") // 2 -> 1 -> 5

	linkedMap.Set("6", 6) // 6 -> 1 -> 5
	suite.Equal(5, linkedMap.Get("5"))
	suite.Equal(3, linkedMap.Size())

	linkedMap.Get("5")    // 5 -> 6 -> 1
	linkedMap.Get("1")    // 5 -> 6 -> 1
	linkedMap.Set("6", 6) // 5 -> 1 -> 6

	linkedMap.Set("7", 7) // 7 -> 1 -> 6
	suite.Nil(linkedMap.Get("5"))
}

func (suite *LinkedMapTestSuite) TestWorkInParallel() {
	linkedMap := NewLinkedMap(3)
	wg := new(sync.WaitGroup)

	tmpl := func(do func()) {
		defer wg.Done()
		do()
	}

	wg.Add(9)

	for i := 0; i < 3; i++ {
		go tmpl(func() {
			linkedMap.Set("1", 1)
		})
	}

	for i := 0; i < 3; i++ {
		go tmpl(func() {
			linkedMap.Set("2", 1)
		})
	}

	go tmpl(func() {
		linkedMap.Set("4", 1)
		linkedMap.Get("4")
	})

	go tmpl(func() {
		linkedMap.Set("5", 1)
		linkedMap.Get("5")
	})

	go tmpl(func() {
		linkedMap.Set("6", 1)
		linkedMap.Get("6")
	})
	wg.Wait()
}
