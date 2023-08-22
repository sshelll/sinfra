package promise

import (
	"testing"
	"time"
)

func TestAllOK(t *testing.T) {
	p0 := New(func(resolve, reject func(v any)) {
		resolve("ok_0")
	})

	p1 := New(func(resolve, reject func(v any)) {
		resolve("ok_1")
	})

	p2 := New(func(resolve, reject func(v any)) {
		resolve("ok_2")
	})

	all := All(p0, p1, p2).Then(func(v any) any {
		values := v.([]any)
		if len(values) != 3 {
			t.Errorf("expected 3 values, got %d", len(values))
		}
		return values
	})

	all.Await()
	t.Log("done, stat:", all.State())
	t.Log("done, result:", all.Result())
}

func TestAllOneFailed(t *testing.T) {
	p0 := New(func(resolve, reject func(v any)) {
		resolve("ok_0")
	})

	p1 := New(func(resolve, reject func(v any)) {
		resolve("ok_1")
	})

	p2 := New(func(resolve, reject func(v any)) {
		reject("failed_2")
	})

	all := All(p0, p1, p2).Then(func(v any) any {
		values := v.([]any)
		return values
	}).Catch(func(a any) {
		t.Log("catch:", a)
	})

	all.Await()
	t.Log("done, stat:", all.State())
	t.Log("done, result:", all.Result())
}

func TestRaceOK(t *testing.T) {
	p0 := New(func(resolve, reject func(v any)) {
		time.Sleep(1 * time.Second)
		resolve("ok_0")
	})

	p1 := New(func(resolve, reject func(v any)) {
		time.Sleep(2 * time.Second)
		resolve("ok_1")
	})

	p2 := New(func(resolve, reject func(v any)) {
		time.Sleep(3 * time.Second)
		resolve("ok_2")
	})

	all := Race(p0, p1, p2).Then(func(v any) any {
		t.Log("race finished:", v)
		return v
	})

	all.Await()
	t.Log("done, stat:", all.State())
	t.Log("done, result:", all.Result())
	time.Sleep(time.Second)
}

func TestRaceOneFailed(t *testing.T) {
	p0 := New(func(resolve, reject func(v any)) {
		time.Sleep(1 * time.Second)
		reject("fail_0")
	})

	p1 := New(func(resolve, reject func(v any)) {
		time.Sleep(2 * time.Second)
		resolve("ok_1")
	})

	p2 := New(func(resolve, reject func(v any)) {
		time.Sleep(3 * time.Second)
		resolve("ok_2")
	})

	all := Race(p0, p1, p2).Then(func(v any) any {
		t.Log("race finished:", v)
		return v
	}).Catch(func(a any) {
		t.Log("catch:", a)
	})

	all.Await()
	t.Log("done, stat:", all.State())
	t.Log("done, result:", all.Result())
	time.Sleep(time.Second)
}
