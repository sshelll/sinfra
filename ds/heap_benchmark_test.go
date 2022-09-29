package ds

import (
	"math/rand"
	"testing"
)

var (
	heap = NewHeap(intComparator, intTypeChecker)
)

func BenchmarkOffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !heap.Offer(rand.Intn(b.N)) {
			b.Error("offer failed")
		}
	}
}
