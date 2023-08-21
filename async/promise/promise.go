package promise

import (
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
	catchMu  *sync.Mutex
	final    func(any)
	finalMu  *sync.Mutex
}

func New(fn func(resolve, reject func(v any))) *Promise {
	p := &Promise{
		state:   PENDING,
		done:    make(chan struct{}),
		catchMu: &sync.Mutex{},
		finalMu: &sync.Mutex{},
	}

	go func() {
		defer close(p.done)
		defer func() {
			if r := recover(); r != nil {
				p.rejected = r
				p.state = REJECTED
			}
			if p.state == REJECTED {
				p.catchMu.Lock()
				if p.catch != nil {
					p.catch(p.rejected)
					p.state = FULFILLED
				}
				p.catchMu.Unlock()
			}
			p.finalMu.Lock()
			if p.final != nil {
				p.final(p.Result())
			}
			p.finalMu.Unlock()
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

func (p *Promise) Result() any {
	if p.state == REJECTED {
		return p.rejected
	}
	return p.resolved
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

func (p *Promise) Then(onResolve func(any) any) *Promise {
	return New(func(resolve, reject func(v any)) {
		p.Await()
		if p.state == FULFILLED {
			resolve(onResolve(p.resolved))
		} else {
			reject(p.rejected)
		}
	})
}

func (p *Promise) Catch(onReject func(any)) *Promise {
	p.catchMu.Lock()
	defer p.catchMu.Unlock()
	select {
	case <-p.done:
		onReject(p.rejected)
		p.state = FULFILLED
	default:
	}
	p.catch = onReject
	return p
}

func (p *Promise) Final(fn func(any)) *Promise {
	p.finalMu.Lock()
	defer p.finalMu.Unlock()
	select {
	case <-p.done:
		fn(p.Result())
	default:
	}
	p.final = fn
	return p
}
