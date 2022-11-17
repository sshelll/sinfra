package stream

import (
	"context"
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

func NewIOStreamWithCap(maxDatapackCnt int) *IOStream {
	if maxDatapackCnt < 0 {
		maxDatapackCnt = 0
	}
	return &IOStream{
		mu:     &sync.Mutex{},
		dataCh: make(chan Datapack, maxDatapackCnt),
		ctrlCh: make(chan struct{}),
	}
}

func NewClosedIOStream(datapacks ...Datapack) *IOStream {
	iostream := NewIOStreamWithCap(len(datapacks))
	for i := range datapacks {
		iostream.Write(datapacks[i])
	}
	iostream.Close()
	return iostream
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
	Context() context.Context
	ReadCloser() io.ReadCloser
}

type simpleDatapack struct {
	r   io.ReadCloser
	ctx context.Context
}

func NewSimpleDatapack(ctx context.Context, r io.ReadCloser) *simpleDatapack {
	return &simpleDatapack{
		r:   r,
		ctx: ctx,
	}
}

func (s *simpleDatapack) ReadCloser() io.ReadCloser {
	return s.r
}

func (s *simpleDatapack) Context() context.Context {
	return s.ctx
}
