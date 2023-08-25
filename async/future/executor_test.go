package future

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFutureOK(t *testing.T) {
	newTask := func(id int) *Task[int] {
		idstr := fmt.Sprintf("task-%d", id)
		var callback Callback[int] = func(ctx context.Context) int {
			for i := 0; i < 2; i++ {
				select {
				case <-ctx.Done():
					panic(context.Cause(ctx))
				default:
				}
				time.Sleep(time.Millisecond * 100)
			}
			return id
		}
		task := NewTask(nil, idstr, time.Second*100, callback)
		return task
	}

	// submit tasks and get futures
	futures := make([]Future[int], 0, 6)
	for i := 0; i < 6; i++ {
		f := Submit(newTask(i))
		futures = append(futures, f)
	}

	// wait for futures
	for _, f := range futures {
		result, err := f.Get()
		assert.Nil(t, err)
		assert.Equal(t, f.State(), COMPLETED)
		assert.Equal(t, f.ID(), fmt.Sprintf("task-%d", result))
	}
}

func TestFutureCanceledByUser(t *testing.T) {
	newTask := func(id int) *Task[any] {
		idstr := fmt.Sprintf("task-%d", id)
		return NewTask(nil, idstr, time.Second*5, func(ctx context.Context) any {
			for i := 0; i < 10; i++ {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}
				time.Sleep(time.Second)
			}
			return id
		})
	}
	exec := NewExecutor[any](4)
	futures := make([]Future[any], 0, 6)
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
		_, err := f.Get()
		assert.NotNil(t, err)
		assert.Equal(t, f.State(), INCOMPLETE)
	}
}

func TestClose(t *testing.T) {
	exec := NewExecutor[any](4)
	exec.Close()
	future := exec.Submit(NewTask(nil, "task-0", time.Second*5, func(ctx context.Context) any {
		panic("should not call me")
	}))
	result, err := future.Get()
	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestCloseTaskBeforeSubmit(t *testing.T) {
	exec := NewExecutor[any](4)
	task := NewTask(context.Background(), "closed_task", time.Second, func(ctx context.Context) any {
		panic("should not call me")
	})
	task.Cancel(nil)
	f := exec.Submit(task)
	result, err := f.Get()
	assert.Nil(t, result)
	assert.NotNil(t, err)
}
