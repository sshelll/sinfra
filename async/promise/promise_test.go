package promise

import (
	"encoding/json"
	"net/http"
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

func TestPromiseCatchThenError(t *testing.T) {
	p := New(func(resolve, reject func(v any)) {
		time.Sleep(time.Millisecond * 500)
		reject("main failed")
	}).Catch(func(a any) {
		t.Log("catch 1st, reject:", a)
	}).Then(func(a any) any {
		t.Log("then 1st, resolve:", a)
		panic("then 1st panic")
	}).Catch(func(a any) {
		t.Log("catch 2nd, reject:", a)
	})
	p.Await()
	t.Log("final state =", p.State().String())
}

func TestPromiseNotCatchMainError(t *testing.T) {
	p0 := New(func(resolve, reject func(v any)) {
		time.Sleep(time.Millisecond * 500)
		reject("main failed")
	})
	p1 := p0.Then(func(a any) any {
		t.Log("then 1st, resolve:", a)
		return "then 1st"
	})
	p1.Await()
	t.Log("p0 final state =", p0.State().String())
	t.Log("p1 final state =", p1.State().String())
}

func TestPromiseNotCatchThenError(t *testing.T) {
	p0 := New(func(resolve, reject func(v any)) {
		time.Sleep(time.Millisecond * 500)
		resolve("main ok")
	}).Then(func(a any) any {
		t.Log("then 1st, resolve:", a)
		panic("then 1st panic")
	})
	p0.Await()
	t.Log("p0 final state =", p0.State().String())
}

func TestPromiseWithRealOp(t *testing.T) {
	p := New(func(resolve, reject func(v any)) {
		time.Sleep(time.Millisecond * 500)
		resolve([]byte(`{"user_id": "100001", "user_name": "sshelll"}`))
	}).Then(func(a any) any {
		resp := a.([]byte)
		t.Log("fetched user info, resp size =", len(resp))
		data := make(map[string]string, 2)
		if err := json.Unmarshal(resp, &data); err != nil {
			panic(err)
		}
		t.Log("resp.user_id =", data["user_id"], ", resp.user_name =", data["user_name"])
		return data["user_id"]
	}).Then(func(a any) any {
		t.Log("got user_id =", a)
		return "user_name = sshelll"
	}).Final(func(a any) {
		t.Log("final result =", a)
	})
	p.Await()
	t.Log("outer final state =", p.State().String())
	t.Log("outer final result =", p.Result())
}

func TestPromiseThenReturnAnotherPromise(t *testing.T) {
	p := New(func(resolve, reject func(v any)) {
		time.Sleep(time.Millisecond * 500)
		resolve("1st")
	}).Then(func(a any) any {
		t.Log("then 1st, resolve:", a)
		return New(func(resolve, reject func(v any)) {
			time.Sleep(time.Millisecond * 500)
			resolve("inside new promise")
		}).Final(func(a any) {
			t.Log("final result from inner promise =", a)
		})
	}).Then(func(a any) any {
		t.Log("then 2nd, resolve:", a)
		return "p final result from 2nd then"
	})
	p.Await()
	t.Log("outer final state =", p.State().String())
	t.Log("outer final result =", p.Result())
}

func TestPromiseWithRealHttpReq(t *testing.T) {
	api := func(url string) *Promise {
		return New(func(resolve, reject func(v any)) {
			resp, err := http.Get(url)
			if err != nil {
				reject(err)
			} else {
				if resp.StatusCode != http.StatusOK {
					reject(resp.Status)
				} else {
					resolve(resp.Status)
				}
			}
		}).Catch(func(a any) {
			t.Log("catch error:", a)
		})
	}

	seq := func() *Promise {
		return New(func(resolve, reject func(v any)) {
			resolve("seq")
		}).Then(func(a any) any {
			t.Log("seq 1st start, resove result =", a)
			return api("https://cdn.staticfile.org/jquery/2.0.3/jquery.min.js")
		}).Then(func(a any) any {
			t.Log("seq 2nd start, resove result =", a)
			return api("https://static.runoob.com/assets/upvotejs/dist/upvotejs/upvotejs.jquery.js")
		}).Then(func(a any) any {
			t.Log("seq 3rd start, resove result =", a)
			return api("https://www.qq.com")
		}).Then(func(a any) any {
			t.Log("seq 4th start, resove result =", a)
			return "done"
		}).Catch(func(a any) {
			t.Log("catch error:", a)
		}).Then(func(a any) any {
			t.Log("seq 5th start, resove result =", a)
			return "done"
		}).Final(func(a any) {
			t.Log("final func, result =", a)
		})
	}

	p := seq()
	p.Await()
	t.Log("outer final state =", p.State().String())
	t.Log("outer final result =", p.Result())
}
