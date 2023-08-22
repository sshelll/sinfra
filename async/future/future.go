package future

import (
	"time"
)

type Future interface {
	Get() any
	GetWithTimeout(timeout time.Duration) any
	Cancel()
	IsDone() bool
	State() State
}

type taskFuture struct {
	task *Task
}

func (tf *taskFuture) Get() any {
	select {
	case <-tf.task.done:
		return tf.task.Result()
	}
}

func (tf *taskFuture) GetWithTimeout(timeout time.Duration) any {
	select {
	case <-tf.task.done:
		return tf.task.Result()
	case <-time.After(timeout):
		return nil
	}
}

func (tf *taskFuture) Cancel() {
	tf.task.Cancel()
}

func (tf *taskFuture) IsDone() bool {
	select {
	case <-tf.task.done:
		return true
	default:
		return false
	}
}

func (tf *taskFuture) State() State {
	return tf.task.State()
}
