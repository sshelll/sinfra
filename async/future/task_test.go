package future

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "test_name", "TestTask")

	task := NewTask(ctx, time.Second*5, func(ctx context.Context) any {
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
		task.Cancel()
	}()
	task.Run(&sync.WaitGroup{})

	t.Log("result =", task.Result())

	select {
	case <-ctx.Done():
		t.Log("ctx done:", ctx.Err())
	default:
		t.Log("ctx not done")
	}
}
