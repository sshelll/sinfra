package future

import (
	"context"
	"testing"
	"time"
)

func TestFuture(t *testing.T) {
	newTask := func(id int) *Task {
		return NewTask(nil, time.Second*5, func(ctx context.Context) any {
			for i := 0; i < 10; i++ {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					t.Logf("[%d]task still running, loop:%d", id, i)
				}
				time.Sleep(time.Second)
			}
			return id
		})
	}
	exec := NewExecutor(4)
	futures := make([]Future, 0, 6)
	for i := 0; i < 6; i++ {
		f := exec.Submit(newTask(i))
		futures = append(futures, f)
	}
	go func() {
		time.Sleep(time.Second)
		for _, f := range futures {
			f.Cancel()
		}
	}()
	for i, f := range futures {
		t.Logf("result of %dth future = %v, state = %s, done = %v",
			i, f.Get(), f.State().String(), f.IsDone())
	}
}
