package future

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type State int

const (
	PENDING State = iota
	INCOMPLETE
	COMPLETED
)

func (s State) String() string {
	switch s {
	case INCOMPLETE:
		return "INCOMPLETE"
	case COMPLETED:
		return "COMPLETED"
	default:
		return "UNDEFINED"
	}
}

type Callback[T any] func(context.Context) T

// Task T is the return type of the task, see Callback.
type Task[T any] struct {
	id        string
	fn        Callback[T]
	ctx       context.Context
	canceler  context.CancelCauseFunc
	tcanceler context.CancelFunc

	state      State
	done       chan struct{}
	result     T
	failedInfo error
}

func NewTask[T any](ctx context.Context, id string, timeout time.Duration, callback Callback[T]) *Task[T] {
	if ctx == nil {
		ctx = context.Background()
	}
	// set timeout ctx first to make sure we can get the timeout error
	tctx, tcancelFn := context.WithTimeoutCause(ctx, timeout, fmt.Errorf("task timeout after %v", timeout))
	cctx, cancelFn := context.WithCancelCause(tctx)
	return &Task[T]{
		id:        id,
		fn:        callback,
		ctx:       cctx,
		canceler:  cancelFn,
		tcanceler: tcancelFn,
		done:      make(chan struct{}),
		state:     PENDING,
	}
}

func (t *Task[T]) Run(wg *sync.WaitGroup) {
	if t.IsDone() {
		return
	}

	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}

	defer t.finalize()

	// precheck, maybe the task is canceled before running
	if err := context.Cause(t.ctx); err != nil {
		return
	}

	// run task
	t.result = t.fn(t.ctx)
	t.state = COMPLETED
}

func (t *Task[T]) Cancel(err error) {
	t.canceler(err)
}

func (t *Task[T]) IsDone() bool {
	select {
	case <-t.done:
		return true
	default:
		return false
	}
}

func (t *Task[T]) State() State {
	return t.state
}

func (t *Task[T]) Result() (T, error) {
	return t.result, t.failedInfo
}

func (t *Task[T]) finalize() {
	// mark task as done at last
	defer close(t.done)

	// cancel all context secondly
	defer func() {
		t.canceler(nil)
		t.tcanceler()
	}()

	// panicked
	if r := recover(); r != nil {
		t.state = INCOMPLETE
		err := fmt.Errorf("task panic: %v", r)
		t.canceler(err)
		t.failedInfo = err
		return
	}

	// already canceled manually or by timeout
	if err := context.Cause(t.ctx); err != nil {
		t.state = INCOMPLETE
		t.failedInfo = err
	}
}
