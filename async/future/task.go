package future

import (
	"context"
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

type Task struct {
	fn       func(ctx context.Context) any
	ctx      context.Context
	canceler context.CancelFunc

	state     State
	done      chan struct{}
	result    any
	panicInfo any
}

func NewTask(ctx context.Context, timeout time.Duration, fn func(ctx context.Context) any) *Task {
	if ctx == nil {
		ctx = context.Background()
	}
	cctx, cancelFn := context.WithTimeout(ctx, timeout)
	return &Task{
		fn:       fn,
		ctx:      cctx,
		canceler: cancelFn,
		done:     make(chan struct{}),
		state:    PENDING,
	}
}

func (t *Task) Run(wg *sync.WaitGroup) {
	wg.Add(1)
	defer close(t.done)
	defer func() {
		if r := recover(); r != nil {
			t.panicInfo = r
			t.state = INCOMPLETE
		}
		t.canceler()
		wg.Done()
	}()
	t.result = t.fn(t.ctx)
	t.state = COMPLETED
}

func (t *Task) Cancel() {
	t.canceler()
}

func (t *Task) State() State {
	return t.state
}

func (t *Task) Result() any {
	if t.state == INCOMPLETE {
		return t.panicInfo
	}
	return t.result
}
