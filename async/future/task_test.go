package future

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestTaskCancelByUser(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "test_name", "TestTask")

	task := NewTask(ctx, "task-0", time.Second*5, func(ctx context.Context) any {
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				t.Log("task still running, loop:", i)
			}
			time.Sleep(time.Second)
		}
		return "done"
	})

	go func() {
		time.Sleep(time.Second * 2)
		task.Cancel(errors.New("cancel by user"))
	}()

	wg := &sync.WaitGroup{}
	task.Run(wg)
	wg.Wait()

	result, err := task.Result()
	t.Logf("result = %v, %v, state = %v", result, err, task.State())

	select {
	case <-ctx.Done():
		t.Log("ctx done:", ctx.Err())
		// whether the task is completed or not, the original ctx should not be influenced
		t.Fail()
	default:
		t.Log("ctx not done")
	}
}

func TestTaskCancelAfterTimeout(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "test_name", "TestTask")

	task := NewTask(ctx, "task-0", time.Second*2, func(ctx context.Context) any {
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				t.Log("task still running, loop:", i)
			}
			time.Sleep(time.Second)
		}
		return "done"
	})

	wg := &sync.WaitGroup{}
	task.Run(wg)
	wg.Wait()

	result, err := task.Result()
	t.Logf("result = %v, %v, state = %v", result, err, task.State())

	select {
	case <-ctx.Done():
		t.Log("ctx done:", ctx.Err())
		// whether the task is completed or not, the original ctx should not be influenced
		t.Fail()
	default:
		t.Log("ctx not done")
	}
}
