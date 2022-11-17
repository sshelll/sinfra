package stream

import (
	"io"
	"sync"
)

// IOStream is a stream of Datapack.
type IOStream struct {
	mu     *sync.Mutex
	dataCh chan Datapack
	ctrlCh chan struct{}
}

func NewIOStream() *IOStream {
	return &IOStream{
		mu:     &sync.Mutex{},
		dataCh: make(chan Datapack, 1),
		ctrlCh: make(chan struct{}),
	}
}

func (s *IOStream) Write(data Datapack) (streamClosed bool) {
	s.mu.Lock()
	streamClosed = s.isClosed()
	s.mu.Unlock()
	if streamClosed {
		return
	}
	s.dataCh <- data
	return false
}

func (s *IOStream) Read() (data Datapack, streamClosed bool) {
	dp, ok := <-s.dataCh
	return dp, !ok
}

// TryRead try read datapack in a non-block way.
// NOTE: if streamClosed, data is nil
func (s *IOStream) TryRead() (data Datapack, streamClosed bool) {
	select {
	case data, ok := <-s.dataCh:
		return data, !ok
	default:
		return nil, false
	}
}

func (s *IOStream) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isClosed() {
		return
	}
	close(s.ctrlCh)
	close(s.dataCh)
}

func (s *IOStream) isClosed() bool {
	select {
	case <-s.ctrlCh:
		return true
	default:
		return false
	}
}

// Datapack is a io.ReadCloser with some extra info.
type Datapack interface {
	ReadCloser() io.ReadCloser
	Extra() interface{}
}

type simpleDatapack struct {
	r     io.ReadCloser
	extra interface{}
}

func NewSimpleDatapack(r io.ReadCloser, extra interface{}) *simpleDatapack {
	return &simpleDatapack{
		r:     r,
		extra: extra,
	}
}

func (s *simpleDatapack) ReadCloser() io.ReadCloser {
	return s.r
}

func (s *simpleDatapack) Extra() interface{} {
	return s.extra
}
