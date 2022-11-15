package stream

import (
	"fmt"
	"io"
)

type DatapackProducer interface {
	Next() (datapack Datapack, hasNext bool, err error)
}

type SafeIOStreamWriter struct {
	datapackProducer DatapackProducer
}

func NewSafeIOStreamWriter(p DatapackProducer) *SafeIOStreamWriter {
	return &SafeIOStreamWriter{
		datapackProducer: p,
	}
}

func (s *SafeIOStreamWriter) Start() (*IOStream, *ErrorPasser) {

	outputStream := NewIOStream()
	outputErr := NewErrorPasser()

	go func() {

		defer func() {

			if r := recover(); r != nil {
				err := fmt.Errorf("SafeIOStreamWriter panicked, panic info = %v", r)
				outputErr.Put(err)
			}

			outputErr.Close()
			outputStream.Close()

		}()

		for {
			datapack, hasNext, err := s.datapackProducer.Next()
			if err != nil {
				outputErr.Put(err)
				break
			}
			if datapack == nil {
				continue
			}
			outputStream.Write(datapack)
			if !hasNext {
				break
			}
		}

	}()

	return outputStream, outputErr

}

type SafeIOStreamHandler struct {
	inputStream, outputStream *IOStream
	inputErr, outputErr       *ErrorPasser
	datapackHandler           func(rc io.ReadCloser, extra interface{}) error
	finalizer                 func()
}

func NewSafeIOStreamHandler(
	inputStream *IOStream,
	inputErr *ErrorPasser,
	handler func(io.ReadCloser, interface{}) error,
	finalizer func(),
) *SafeIOStreamHandler {

	return &SafeIOStreamHandler{
		inputStream:     inputStream,
		inputErr:        inputErr,
		datapackHandler: handler,
		finalizer:       finalizer,
	}

}

func (s *SafeIOStreamHandler) BuildStream() (*IOStream, *ErrorPasser) {

	if s.inputStream == nil || s.inputErr == nil {
		return nil, nil
	}

	if s.datapackHandler == nil {
		return s.inputStream, s.inputErr
	}

	s.outputStream = NewIOStream()
	s.outputErr = NewErrorPasserWithCap(s.inputErr.Cap() + 2)

	return s.outputStream, s.outputErr

}

func (s *SafeIOStreamHandler) Start() {

	outputStream, outputErr := s.outputStream, s.outputErr

	if outputStream == nil || outputErr == nil {
		s.BuildStream()
	}

	go func() {

		defer func() {
			if r := recover(); r != nil {
				outputErr.Put(fmt.Errorf("SafeIOStreamHandler panicked, err = %v", r))
			}
			outputErr.Close()
			outputStream.Close()
			s.finalizer()
		}()

		for {
			datapack, closed := s.inputStream.Read()
			if closed {
				break
			}

			rc, extra := datapack.ReadCloser(), datapack.Extra()
			if rc == nil {
				continue
			}

			if err := s.datapackHandler(rc, extra); err != nil {
				outputErr.Put(err)
				break
			}
		}

		// handle input err
		for {
			err, done := s.inputErr.Check()
			if done {
				break
			}
			if err != nil {
				outputErr.Put(err)
			}
		}

	}()

}
