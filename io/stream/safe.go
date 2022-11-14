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

type SafeIOStreamReader struct {
	inputStream     *IOStream
	datapackHandler func(rc io.ReadCloser, extra interface{}) error
	finalizer       func()
}

func NewSafeIOStreamReader(inputStream *IOStream, handler func(io.ReadCloser, interface{}) error,
	finalizer func()) *SafeIOStreamReader {
	return &SafeIOStreamReader{
		inputStream:     inputStream,
		datapackHandler: handler,
		finalizer:       finalizer,
	}
}

func (s *SafeIOStreamReader) Start() (*IOStream, *ErrorPasser) {

	outputStream := NewIOStream()
	outputErr := NewErrorPasser()

	go func() {

		defer func() {
			if r := recover(); r != nil {
				outputErr.Put(fmt.Errorf("zipper panicked, err = %v", r))
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
			}

		}

	}()

	return outputStream, outputErr

}
