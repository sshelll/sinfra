package stream

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func TestProcessor(t *testing.T) {

	// build proc chain
	firstProc := newSimpleProcessor("1st proc")
	secondProc := newSimpleProcessor("2nd proc")
	thirdProc := newSimpleProcessor("3rd proc")

	proc := BuildProcChain(firstProc, secondProc, thirdProc)

	stream, errPasser := NewSafeIOStreamWriter(NewSimpleProducer(3)).Start()

	// start handling upstream data
	outputStream, outputErr := proc(stream, errPasser)

	for i := 1; ; i++ {

		t.Logf("block reading for the %d time", i)
		datapack, closed := outputStream.Read()

		if closed {
			t.Log("output stream closed, exit loop")
			break
		}

		if datapack == nil {
			t.Log("fetched an empty datapck, continue")
			continue
		}

		rc := datapack.ReadCloser()
		if rc == nil {
			t.Log("fetched an empty ReadCloser from datapack, continue")
			continue
		}

		bs, err := ioutil.ReadAll(rc)
		if err != nil {
			t.Logf("read from ReadCloser of datapack failed, err = %v", err)
			break
		}

		t.Logf("upstream data is '%v'", string(bs))

	}

	// check upstream err
	err, done := outputErr.Check()
	t.Logf("upstream err = %v, check done = %v", err, done)

}

func newSimpleProcessor(procName string) Processor {
	return func(inputStream *IOStream, inputErr *ErrorPasser) (
		outputStream *IOStream, outputErr *ErrorPasser) {

		outputStream = NewIOStream()
		outputErr = NewErrorPasser()

		go func() {

			defer outputErr.Close()

			for {
				// block get data from upstream
				data, upstreamClosed := inputStream.Read()
				if upstreamClosed {
					outputStream.Close()
				}

				if data == nil {
					continue
				}

				rc := data.ReadCloser()
				if rc == nil {
					continue
				}

				// start handle upstream data
				bs, err := ioutil.ReadAll(rc)
				if err != nil {
					outputErr.Put(err)
					outputStream.Close()
					return
				}
				rc.Close()

				// modify upstream data
				modifiedStr := string(bs) + " - modified by " + procName
				reader := bytes.NewBuffer([]byte(modifiedStr))

				// send to donwstream
				datapack := NewSimpleDatapack(ioutil.NopCloser(reader), procName)
				outputStream.Write(datapack)
			}

		}()

		return

	}

}

type simpleProducer struct {
	curCnt int
	endCnt int
}

func NewSimpleProducer(cnt int) *simpleProducer {
	return &simpleProducer{
		curCnt: 0,
		endCnt: cnt,
	}
}

func (s *simpleProducer) Next() (datapack Datapack, hasNext bool, err error) {

	// simulate http rtt
	time.Sleep(time.Second)

	if s.curCnt >= s.endCnt {
		return nil, false, nil
	}

	// mock data
	data := fmt.Sprintf("%dth data from upstream", s.curCnt)
	rc := ioutil.NopCloser(bytes.NewBuffer([]byte(data)))
	datapack = NewSimpleDatapack(rc, "upstream")

	// cal hasNext
	s.curCnt++
	hasNext = s.curCnt < s.endCnt

	return

}
