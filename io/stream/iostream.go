package stream

import (
	"io"
)

// IOStream is a stream of Datapack.
type IOStream struct {
	dataCh chan Datapack
	ctrlCh chan struct{}
}

func NewIOStream() *IOStream {
	return &IOStream{
		dataCh: make(chan Datapack, 1),
		ctrlCh: make(chan struct{}),
	}
}

func (s *IOStream) Write(data Datapack) {
	s.dataCh <- data
}

func (s *IOStream) Read() (data Datapack, streamClosed bool) {
	dp, ok := <-s.dataCh
	return dp, !ok
}

func (s *IOStream) TryRead() (data Datapack, streamClosed bool) {
	select {
	case data, ok := <-s.dataCh:
		return data, !ok
	default:
		return nil, false
	}
}

func (s *IOStream) Close() {
	if !s.isClosed() {
		close(s.ctrlCh)
		close(s.dataCh)
	}
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
