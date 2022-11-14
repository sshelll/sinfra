package archive

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/SCU-SJL/sinfra/io/stream"
	"github.com/stretchr/testify/assert"
)

func TestZip(t *testing.T) {

	fileStream, errPasser := stream.NewSafeIOStreamWriter(NewFileProducer()).Start()

	zipper := NewZipper()
	proc := stream.BuildProcChain(zipper.Zip)

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

func TestSafeZip(t *testing.T) {

	fileStream, errPasser := stream.NewSafeIOStreamWriter(NewFileProducer()).Start()

	zipper := NewZipper()
	proc := stream.BuildProcChain(zipper.SafeZip)

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

	datapack = stream.NewSimpleDatapack(f, f.Name())

	p.idx++
	hasNext = p.idx < len(p.files)

	return

}
