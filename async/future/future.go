package future

import (
	"time"
)

type Future interface {
	ID() string
	Get() (any, error)
	GetWithTimeout(timeout time.Duration) (any, error)
	Cancel(error)
	IsDone() bool
	State() State
}

type taskFuture struct {
	task *Task
}

func (tf *taskFuture) ID() string {
	return tf.task.id
}

func (tf *taskFuture) Get() (any, error) {
	select {
	case <-tf.task.done:
		return tf.task.Result()
	}
}

func (tf *taskFuture) GetWithTimeout(timeout time.Duration) (any, error) {
	select {
	case <-tf.task.done:
		return tf.task.Result()
	case <-time.After(timeout):
		return nil, nil
	}
}

func (tf *taskFuture) Cancel(err error) {
	tf.task.Cancel(err)
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
