package future

import (
	"errors"
	"sync"
)

type Executor[T any] struct {
	wg     *sync.WaitGroup
	taskCh chan *Task[T]
	sema   chan struct{}
}

func NewExecutor[T any](limit int) *Executor[T] {
	f := &Executor[T]{
		wg:     &sync.WaitGroup{},
		taskCh: make(chan *Task[T], limit),
		sema:   make(chan struct{}, limit),
	}
	go f.run()
	return f
}

// Submit submits a task immediately without waiting for an executor to run.
func Submit[T any](task *Task[T]) Future[T] {
	go func() {
		task.Run(nil)
	}()
	return &taskFuture[T]{task: task}
}

func (exec *Executor[T]) Submit(task *Task[T]) Future[T] {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				task.Cancel(errors.New("executor closed"))
				task.finalize()
			}
		}()
		exec.taskCh <- task
	}()
	return &taskFuture[T]{task: task}
}

func (exec *Executor[T]) Close() {
	close(exec.taskCh)
}

func (exec *Executor[T]) Wait() {
	exec.wg.Wait()
}

func (exec *Executor[T]) run() {
	for {
		task, ok := <-exec.taskCh
		if task == nil {
			continue
		}
		exec.sema <- struct{}{}
		go func(task *Task[T]) {
			task.Run(exec.wg)
			<-exec.sema
		}(task)
		if !ok {
			break
		}
	}
}
