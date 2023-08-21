package promise

import (
	"testing"
	"time"
)

func TestPromiseThenWithoutChain(t *testing.T) {
	p := New(func(resolve, reject func(any)) {
		resolve("1st")
	})

	time.Sleep(time.Millisecond * 50)
	p.Await()
	t.Log("promise result =", p.resolved)

	p.Then(func(v any) any {
		t.Log("I.", v)
		return "2nd"
	}).Await()

	p.Then(func(v any) any {
		t.Log("II.", v)
		return "3rd"
	}).Await()

	p.Then(func(v any) any {
		t.Log("III.", v)
		return "4th"
	}).Await()
}

func TestPromiseThenWithChain(t *testing.T) {
	p := New(func(resolve, reject func(any)) {
		resolve("1st")
	})

	p.Then(func(v any) any {
		t.Log("I.", v)
		return "2nd"
	}).Then(func(v any) any {
		t.Log("II.", v)
		return "3rd"
	}).Then(func(v any) any {
		t.Log("III.", v)
		return "4th"
	}).Await()
}

func TestPromiseCatch(t *testing.T) {
	p := New(func(resolve, reject func(v any)) {
		time.Sleep(time.Millisecond * 500)
		reject("main failed")
	}).Catch(func(v any) {
		t.Log("reject", v)
	})
	p.Await()
	time.Sleep(time.Millisecond * 500)
}

func TestPromiseComplex(t *testing.T) {
	p := New(func(resolve, reject func(v any)) {
		time.Sleep(time.Millisecond * 500)
		reject("main failed")
	}).Catch(func(a any) {
		t.Log("catch 1st, reject:", a)
	}).Then(func(a any) any {
		t.Log("then 1st, resolve:", a)
		panic("then 1st panic")
	})
	p.Await()
	t.Log("final state =", p.State().String())
}
