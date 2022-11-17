package archive

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/SCU-SJL/sinfra/io/stream"
)

type Zipper struct {
	downstreamCtx    context.Context
	zw               *zip.Writer
	pr               *io.PipeReader
	pw               *io.PipeWriter
	zipFilename      string
	ctxKeyOfFileName string
}

func NewZipper() *Zipper {
	zipper := &Zipper{}
	zipper.pr, zipper.pw = io.Pipe()
	zipper.zw = zip.NewWriter(zipper.pw)
	return zipper
}

func (z *Zipper) SetDownstreamCtx(ctx context.Context) {
	z.downstreamCtx = ctx
}

func (z *Zipper) SetCtxKeyOfFileName(key string) {
	z.ctxKeyOfFileName = key
}

// Proc note that outputStream only contains one datapack.
func (z *Zipper) Proc(inputStream *stream.IOStream, inputErr *stream.ErrorPasser) (
	outputStream *stream.IOStream, outputErr *stream.ErrorPasser) {

	safeHandler := stream.NewSafeIOStreamHandler(inputStream, inputErr, z.datapackHandleFn, z.finalizeFn)

	outputStream, outputErr = safeHandler.BuildStream()

	ctx := z.downstreamCtx
	if ctx == nil {
		ctx = context.Background()
	}

	outputStream.Write(stream.NewSimpleDatapack(ctx, z.pr))

	safeHandler.Start()

	return

}

func (z *Zipper) datapackHandleFn(ctx context.Context, rc io.ReadCloser) error {

	if rc == nil {
		return errors.New("extra should be a string represents filename, but got nil")
	}

	v := ctx.Value(z.ctxKeyOfFileName)
	if v == nil {
		return fmt.Errorf("cannot find file name with key '%s'", z.ctxKeyOfFileName)
	}

	filename, ok := v.(string)
	if !ok || len(strings.TrimSpace(filename)) == 0 {
		return fmt.Errorf("ctx value with key '%s' should be a string represents filename, but got type = %s, value = %v",
			z.ctxKeyOfFileName, reflect.TypeOf(v).Name(), v)
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
