package future

import (
	"sync"
)

type Executor struct {
	wg     *sync.WaitGroup
	taskCh chan *Task
	sema   chan struct{}
}

func NewExecutor(limit int) *Executor {
	f := &Executor{
		wg:     &sync.WaitGroup{},
		taskCh: make(chan *Task, limit),
		sema:   make(chan struct{}, limit),
	}
	go f.run()
	return f
}

func (exec *Executor) Submit(task *Task) Future {
	go func() {
		exec.taskCh <- task
	}()
	return &taskFuture{task: task}
}

func (exec *Executor) Close() {
	close(exec.taskCh)
}

func (exec *Executor) Wait() {
	exec.wg.Wait()
}

func (exec *Executor) run() {
	for {
		task, ok := <-exec.taskCh
		if !ok {
			break
		}
		if task == nil {
			continue
		}
		exec.sema <- struct{}{}
		go func(task *Task) {
			task.Run(exec.wg)
			<-exec.sema
		}(task)
	}
}
