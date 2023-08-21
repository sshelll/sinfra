package promise

import (
	"fmt"
	"sync"
)

type State int

func (s State) String() string {
	switch s {
	case PENDING:
		return "PENDING"
	case FULFILLED:
		return "FULFILLED"
	case REJECTED:
		return "REJECTED"
	default:
		return "UNDEFINED"
	}
}

const (
	PENDING State = iota
	FULFILLED
	REJECTED
)

type Promise struct {
	state    State
	done     chan struct{}
	resolved interface{}
	rejected interface{}
	catch    func(any)
	mu       *sync.Mutex
}

func New(fn func(resolve, reject func(v any))) *Promise {
	p := &Promise{
		state: PENDING,
		done:  make(chan struct{}),
		mu:    &sync.Mutex{},
	}

	go func() {
		defer close(p.done)
		defer func() {
			if r := recover(); r != nil {
				p.rejected = r
				p.state = REJECTED
			}
			if p.state == REJECTED {
				p.mu.Lock()
				if p.catch != nil {
					p.catch(p.rejected)
					p.state = FULFILLED
				}
				p.mu.Unlock()
			}
		}()
		fn(p.Resolve, p.Reject)
	}()

	return p
}

func (p *Promise) Resolve(v any) {
	p.state = FULFILLED
	p.resolved = v
}

func (p *Promise) Reject(v any) {
	p.state = REJECTED
	p.rejected = v
}

func (p *Promise) State() State {
	return p.state
}

func (p *Promise) Await() {
	select {
	case <-p.done:
	}
}

func (p *Promise) Panic() {
	p.Await()
	if r := p.rejected; r != nil {
		panic(r)
	}
}

func (p *Promise) Then(cb func(any) any) *Promise {
	return New(func(resolve, reject func(v any)) {
		fmt.Println("then called")
		p.Await()
		if p.state == FULFILLED {
			resolve(cb(p.resolved))
		}
	})
}

func (p *Promise) Catch(cb func(any)) *Promise {
	p.mu.Lock()
	defer p.mu.Unlock()
	select {
	case <-p.done:
		cb(p.rejected)
		p.state = FULFILLED
	default:
	}
	p.catch = cb
	return p
}
