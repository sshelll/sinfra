package promise

import (
	"sync"
)

// NOTE: Promise will not support generic type, or the chain call will be very ugly.

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
	state       State
	done        chan struct{}
	resolved    interface{}
	rejected    interface{}
	confirmedCh chan struct{}
	resultMu    *sync.Mutex
	catch       func(any)
	catchMu     *sync.Mutex
	final       func(any)
	finalMu     *sync.Mutex
}

func New(fn func(resolve, reject func(v any))) *Promise {
	p := &Promise{
		state:       PENDING,
		done:        make(chan struct{}),
		resultMu:    &sync.Mutex{},
		catchMu:     &sync.Mutex{},
		finalMu:     &sync.Mutex{},
		confirmedCh: make(chan struct{}),
	}

	go func() {
		// close chan must be the last step
		defer close(p.done)

		defer func() {
			// panic protection
			if r := recover(); r != nil {
				p.reject(r)
			}

			// waiting for confirm
			<-p.confirmedCh

			// check state
			if p.state == PENDING {
				panic("promise: promise is still pending")
			}

			// catch error
			if p.state == REJECTED {
				p.catchMu.Lock()
				if p.catch != nil {
					p.catch(p.rejected)
					p.state = FULFILLED
				}
				p.catchMu.Unlock()
			}

			// final
			p.finalMu.Lock()
			if p.final != nil {
				p.final(p.Result())
			}
			p.finalMu.Unlock()
		}()

		fn(p.resolve, p.reject)
	}()

	return p
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
	if p.state == REJECTED {
		panic(p.rejected)
	}
}

func (p *Promise) Done() <-chan struct{} {
	return p.done
}

func (p *Promise) Then(onResolve func(any) any) *Promise {
	return New(func(resolve, reject func(v any)) {
		// wait for the last promise to be done
		p.Await()
		// if the result of the last promise is another promise,
		// wait for it to be done
		lastResult := p.Result()
		if newp, ok := lastResult.(*Promise); ok {
			newp.Await()
			lastResult = newp.Result()
		}
		if p.state == FULFILLED {
			// handle the result first, then mark as resolved
			resolve(onResolve(lastResult))
		} else {
			// mark as rejected immediately
			reject(lastResult)
		}
	})
}

func (p *Promise) Catch(onReject func(any)) *Promise {
	p.catchMu.Lock()
	defer p.catchMu.Unlock()
	p.catch = onReject
	select {
	case <-p.done:
		if p.state == REJECTED {
			onReject(p.rejected)
			p.state = FULFILLED
		}
	default:
	}
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

func (p *Promise) resolve(v any) {
	defer p.confirm()
	p.resultMu.Lock()
	defer p.resultMu.Unlock()
	if p.state != PENDING {
		return
	}
	p.state = FULFILLED
	p.resolved = v
	p.rejected = nil
}

func (p *Promise) reject(v any) {
	defer p.confirm()
	p.resultMu.Lock()
	defer p.resultMu.Unlock()
	if p.state != PENDING {
		return
	}
	p.state = REJECTED
	p.rejected = v
	p.resolved = nil
}

func (p *Promise) confirm() {
	select {
	case <-p.confirmedCh:
	default:
		close(p.confirmedCh)
	}
}
