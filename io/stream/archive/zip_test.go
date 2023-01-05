package archive

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sshelll/sinfra/io/stream"
	"github.com/stretchr/testify/assert"
)

func TestZipWithFilepath(t *testing.T) {
	zipFile, _ := os.Create("tmp.zip")
	zw := zip.NewWriter(zipFile)
	dst, err := zw.Create("./dir/zip_test.go")
	assert.Nil(t, err)
	src, _ := os.Open("./zip_test.go")
	_, err = io.Copy(dst, src)
	assert.Nil(t, err)
	src.Close()
	zw.Close()
}

func TestSafeZip(t *testing.T) {

	fileStream, errPasser := stream.NewSafeIOStreamWriter(NewFileProducer()).Start()

	zipper := NewZipper()
	proc := stream.BuildProcChain(zipper.Proc)

	outputStream, outputErr := proc(fileStream, errPasser)

	f, err := os.Create("test_result.zip")
	assert.Nil(t, err)

	for {
		zipData, closed := outputStream.Read()
		if closed {
			t.Log("zip stream closed")
			break
		}
		t.Log("start handling zip data")
		zr := zipData.ReadCloser()
		_, err := io.Copy(f, zr)
		assert.Nil(t, err)
	}

	err, done := outputErr.Check()
	assert.Nil(t, err)
	assert.True(t, done)

}

func TestZipWithUpstreamErr(t *testing.T) {

	iostream, errPasser := stream.NewSafeIOStreamWriter(&ErrorProducer{}).Start()

	zipper := NewZipper()
	proc := stream.BuildProcChain(zipper.Proc)

	outputStream, outputErr := proc(iostream, errPasser)

	for {
		data, closed := outputStream.Read()
		if closed {
			t.Log("input stream closed")
			break
		}
		t.Log("start handling input")
		rc := data.ReadCloser()
		t.Logf("rc is nil = %v", rc == nil)
		bs, err := ioutil.ReadAll(rc)
		t.Logf("err = %v", err)
		t.Logf("input is %s", string(bs))
	}

	for {
		err, done := outputErr.Check()
		if done {
			break
		}
		t.Log("err is", err.Error())
	}

}

type FileProducer struct {
	idx   int
	files []string
}

func NewFileProducer() *FileProducer {
	files := make([]string, 0, 4)
	filepath.Walk("./", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, "go") {
			files = append(files, path)
		}
		return nil
	})
	return &FileProducer{
		idx:   0,
		files: files,
	}
}

func (p *FileProducer) Next() (datapack stream.Datapack, hasNext bool, err error) {

	time.Sleep(time.Second)

	f, err := os.Open(p.files[p.idx])
	if err != nil {
		return nil, false, err
	}

	datapack = stream.NewSimpleDatapack(context.WithValue(context.Background(), "filename", f.Name()), f)

	p.idx++
	hasNext = p.idx < len(p.files)

	return

}

type ErrorProducer struct {
	idx int
}

func (p *ErrorProducer) Next() (datapack stream.Datapack, hasNext bool, err error) {
	if p.idx > 0 {
		err = errors.New("some error from ErrorProducer")
		return
	}
	p.idx++
	r := bytes.NewBufferString("hello world")
	rc := io.NopCloser(r)
	datapack = stream.NewSimpleDatapack(context.WithValue(context.Background(), "filename", "hello.txt"), rc)
	hasNext = true
	return
}
