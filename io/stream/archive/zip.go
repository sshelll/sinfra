package archive

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/SCU-SJL/sinfra/io/stream"
)

type Zipper struct {
	zw          *zip.Writer
	pr          *io.PipeReader
	pw          *io.PipeWriter
	getFilename func(extra interface{}) string
}

func NewZipper() *Zipper {
	zipper := &Zipper{}
	zipper.pr, zipper.pw = io.Pipe()
	zipper.zw = zip.NewWriter(zipper.pw)
	return zipper
}

func (z *Zipper) SetGetFileNameFromExtraFn(fn func(extra interface{}) string) {
	z.getFilename = fn
}

// SafeZip note that outputStream only contains one datapack.
func (z *Zipper) SafeZip(inputStream *stream.IOStream, inputErr *stream.ErrorPasser) (
	outputStream *stream.IOStream, outputErr *stream.ErrorPasser) {

	safeHandler := stream.NewSafeIOStreamHandler(inputStream, inputErr, z.datapackHandleFn, z.finalizeFn)

	outputStream, outputErr = safeHandler.BuildStream()
	outputStream.Write(stream.NewSimpleDatapack(z.pr, nil))

	safeHandler.Start()

	return

}

func (z *Zipper) datapackHandleFn(rc io.ReadCloser, extra interface{}) error {

	if rc == nil {
		return nil
	}

	if extra == nil {
		return errors.New("extra should be a string represents filename, but got nil")
	}

	var filename string

	if z.getFilename != nil {
		filename = z.getFilename(extra)
	} else {
		ok := false
		filename, ok = extra.(string)
		if !ok || len(strings.TrimSpace(filename)) == 0 {
			return fmt.Errorf("extra should be a string represents filename, but got type = %s, value = %v",
				reflect.TypeOf(extra).Name(), extra)
		}
	}

	fw, err := z.zw.Create(filename)
	if err != nil {
		return err
	}

	if _, err := io.Copy(fw, rc); err != nil {
		return err
	}

	return rc.Close()

}

func (z *Zipper) finalizeFn() {
	z.zw.Close()
	z.pw.Close()
}
