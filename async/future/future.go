package future

import (
	"time"
)

type Future[T any] interface {
	ID() string
	Get() (T, error)
	GetWithTimeout(timeout time.Duration) (T, error)
	Cancel(error)
	IsDone() bool
	State() State
}

type taskFuture[T any] struct {
	task *Task[T]
}

func (tf *taskFuture[T]) ID() string {
	return tf.task.id
}

func (tf *taskFuture[T]) Get() (T, error) {
	select {
	case <-tf.task.done:
		return tf.task.Result()
	}
}

func (tf *taskFuture[T]) GetWithTimeout(timeout time.Duration) (T, error) {
	select {
	case <-tf.task.done:
		return tf.task.Result()
	case <-time.After(timeout):
		var result T
		return result, nil
	}
}

func (tf *taskFuture[T]) Cancel(err error) {
	tf.task.Cancel(err)
}

func (tf *taskFuture[T]) IsDone() bool {
	select {
	case <-tf.task.done:
		return true
	default:
		return false
	}
}

func (tf *taskFuture[T]) State() State {
	return tf.task.State()
}
