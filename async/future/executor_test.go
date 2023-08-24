package future

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestFutureOK(t *testing.T) {
	newTask := func(id int) *Task {
		idstr := fmt.Sprintf("task-%d", id)
		return NewTask(nil, idstr, time.Second*5, func(ctx context.Context) any {
			for i := 0; i < 2; i++ {
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
	exec.Wait()
	for _, f := range futures {
		result, err := f.Get()
		t.Logf("result of %s, result = %v, err = %v, state = %s, done = %v",
			f.ID(), result, err, f.State().String(), f.IsDone())
	}
}

func TestFutureCanceledByUser(t *testing.T) {
	newTask := func(id int) *Task {
		idstr := fmt.Sprintf("task-%d", id)
		return NewTask(nil, idstr, time.Second*5, func(ctx context.Context) any {
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
			f.Cancel(errors.New("cancel by user"))
		}
	}()
	for _, f := range futures {
		result, err := f.Get()
		t.Logf("result of %s, result = %v, err = %v, state = %s, done = %v",
			f.ID(), result, err, f.State().String(), f.IsDone())
	}
}

func TestClose(t *testing.T) {
	exec := NewExecutor(4)
	exec.Close()
	future := exec.Submit(NewTask(nil, "task-0", time.Second*5, func(ctx context.Context) any {
		return nil
	}))
	result, err := future.Get()
	t.Logf("result = %v, err = %v", result, err)
}
