package stream

import (
	"fmt"
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
