package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"strings"

	"github.com/SCU-SJL/sinfra/io/stream"
)

type Zipper struct {
	zw *zip.Writer
	pr *io.PipeReader
	pw *io.PipeWriter
}

func NewZipper() *Zipper {
	zipper := &Zipper{}
	zipper.pr, zipper.pw = io.Pipe()
	zipper.zw = zip.NewWriter(zipper.pw)
	return zipper
}

func (z *Zipper) Zip(inputStream *stream.IOStream, inputErr *stream.ErrorPasser) (
	outputStream *stream.IOStream, outputErr *stream.ErrorPasser) {

	if z.zw == nil || z.pr == nil || z.pw == nil {
		panic("zipper is not initialized")
	}

	outputStream = stream.NewIOStream()
	outputErr = stream.NewErrorPasser()

	go outputStream.Write(stream.NewSimpleDatapack(z.pr, nil))

	go func() {

		defer func() {
			if r := recover(); r != nil {
				outputErr.Put(fmt.Errorf("zipper panicked, err = %v", r))
			}
			outputErr.Close()
			outputStream.Close()
			inputStream.Close()
			z.zw.Close()
			z.pw.Close()
		}()

		z.doZip(inputStream, inputErr, outputErr)

	}()

	return

}

func (z *Zipper) doZip(inputStream *stream.IOStream, inputErr, outputErr *stream.ErrorPasser) {

	for {

		datapack, closed := inputStream.Read()
		if closed {
			break
		}

		dataRC := datapack.ReadCloser()
		filename := datapack.Extra().(string)

		// ignore empty data
		if dataRC == nil || len(strings.TrimSpace(filename)) == 0 {
			continue
		}

		fw, err := z.zw.Create(filename)
		if err != nil {
			outputErr.Put(err)
			break
		}

		if _, err := io.Copy(fw, dataRC); err != nil {
			outputErr.Put(err)
			break
		}

		if err := dataRC.Close(); err != nil {
			outputErr.Put(err)
			break
		}

	}

	if err, _ := inputErr.Check(); err != nil {
		outputErr.Put(err)
	}

}
