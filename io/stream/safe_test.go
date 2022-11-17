package stream

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"testing"
)

func TestDownstreamPanic(t *testing.T) {
	stream, ep := NewSafeIOStreamWriter(&StringProducer{}).Start()
	outputStream, outputErr := NewProcWithPanic().Proc(stream, ep)
	for {
		data, closed := outputStream.Read()
		if closed {
			break
		}
		t.Logf("[outter] datapack.extra = %v", data.Extra())
		bs, _ := ioutil.ReadAll(data.ReadCloser())
		t.Logf("[outter] read data: %v", string(bs))
	}
	for {
		err, done := outputErr.Check()
		if done {
			break
		}
		t.Logf("[outter] err from output: %v", err)
	}
}

type StringProducer struct {
	idx int
}

func (p *StringProducer) Next() (datapack Datapack, hasNext bool, err error) {
	log.Println("[StringProducer] Next is called")
	p.idx++
	str := fmt.Sprintf("this is the %dth str", p.idx)
	r := bytes.NewBufferString(str)
	rc := ioutil.NopCloser(r)
	datapack = NewSimpleDatapack(rc, "string producer")
	hasNext = p.idx < 30
	return
}

type ProcWithPanic struct {
	idx int
	pr  *io.PipeReader
	pw  *io.PipeWriter
}

func NewProcWithPanic() *ProcWithPanic {
	inst := &ProcWithPanic{}
	inst.pr, inst.pw = io.Pipe()
	return inst
}

func (p *ProcWithPanic) Proc(inputStream *IOStream, inputErr *ErrorPasser) (
	outputStream *IOStream, outputErr *ErrorPasser) {
	safeHandler := NewSafeIOStreamHandler(inputStream, inputErr, p.datapackHandleFn, p.finalizeFn)
	outputStream, outputErr = safeHandler.BuildStream()
	outputStream.Write(NewSimpleDatapack(p.pr, "proc with panic"))
	safeHandler.Start()
	return
}

func (p *ProcWithPanic) datapackHandleFn(rc io.ReadCloser, extra interface{}) error {

	log.Println("[ProcWithPanic] datapackHandleFn is running")

	p.idx++
	if p.idx > 1 {
		log.Println("start throwing panic")
		panic("proc panic")
	}

	bs, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}
	n, err := p.pw.Write(bs)
	log.Printf("[ProcWithPanic] %d time handle, write %d bytes\n", p.idx, n)
	return err

}

func (p *ProcWithPanic) finalizeFn() {
	log.Println("[ProcWithPanic] finalizeFn is running")
	p.pw.Close()
	p.pr.Close()
}
