package functional

import "testing"

func TestMap(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	res := From(arr).
		Map(func(i int) int { return i * 2 }).
		Collect()
	for i := range arr {
		if arr[i]*2 != res[i] {
			t.Fatalf("expected %d, got %d", arr[i]*2, res[i])
		}
	}
}

func TestMapWithFilter(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	res := From(arr).
		Filter(func(i int) bool { return i > 2 }).
		Map(func(i int) int { return i * 2 }).
		Collect()
	if len(res) != 3 {
		t.Fatalf("expected %d, got %d", 3, len(res))
	}
	for i := range res {
		if arr[i+2]*2 != res[i] {
			t.Fatalf("expected %d, got %d", arr[i+2]*2, res[i])
		}
	}
}

func TestMapTo(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	res := FromM[int, string](arr).
		MapTo(func(i int) string { return string(rune((i + 64))) }).
		CollectTo()
	for i := range arr {
		if string(rune((arr[i] + 64))) != res[i] {
			t.Fatalf("expected %s, got %s", string(rune((arr[i] + 64))), res[i])
		}
	}
}

func TestMapToWithFilter(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	res := FromM[int, string](arr).
		Filter(func(i int) bool { return i > 2 }).
		MapTo(func(i int) string { return string(rune((i + 64))) }).
		CollectTo()
	if len(res) != 3 {
		t.Fatalf("expected %d, got %d", 3, len(res))
	}
	for i := range res {
		if string(rune((arr[i+2] + 64))) != res[i] {
			t.Fatalf("expected %s, got %s", string(rune((arr[i+2] + 64))), res[i])
		}
	}
}

func TestReduce(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	res := From(arr).
		Reduce(func(i, j int) int { return i + j })
	if res != 15 {
		t.Fatalf("expected %d, got %d", 15, res)
	}
}

func TestReduceTo(t *testing.T) {
	arr := []int{2, 3, 4, 5, 6}
	res := FromR[int, string](arr).
		Map(func(i int) int { return i - 1 }).
		ReduceTo(func(str string, i int) string { return str + string(rune(i+64)) })
	if res != "ABCDE" {
		t.Fatalf("expected %s, got %s", "ABCDE", res)
	}
}

func TestMapReduceTo(t *testing.T) {
	arr := []int{2, 3, 4, 5, 6}
	res := FromMR[int, string, int](arr).
		MapTo(func(i int) string { return string(rune(i + 64)) }).
		MapReduceTo(func(res int, str string) int { return res + len(str) })
	if res != 5 {
		t.Fatalf("expected %d, got %d", 5, res)
	}
}
