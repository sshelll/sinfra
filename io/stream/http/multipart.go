package http

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"reflect"
	"strings"

	"github.com/SCU-SJL/sinfra/io/stream"
)

type Multiparter struct {
	downstreamCtx       context.Context
	mw                  *multipart.Writer
	pr                  *io.PipeReader
	pw                  *io.PipeWriter
	fieldName           string
	formDataContentType string
	ctxKeyOfFileName    string
}

func NewMultiparter() *Multiparter {
	m := &Multiparter{}
	m.pr, m.pw = io.Pipe()
	m.mw = multipart.NewWriter(m.pw)
	m.formDataContentType = m.mw.FormDataContentType()
	return m
}

func (m *Multiparter) SetDownstreamCtx(ctx context.Context) {
	m.downstreamCtx = ctx
}

func (m *Multiparter) SetCtxKeyOfFileName(key string) {
	m.ctxKeyOfFileName = key
}

func (m *Multiparter) FormDataContentType() string {
	return m.formDataContentType
}

func (m *Multiparter) Proc(inputStream *stream.IOStream, inputErr *stream.ErrorPasser) (
	outputStream *stream.IOStream, outputErr *stream.ErrorPasser) {

	if m.pr == nil || m.pw == nil || m.mw == nil {
		panic("Multiparter is not initialized!")
	}

	safeHandler := stream.NewSafeIOStreamHandler(inputStream, inputErr, m.datapackHandleFn, m.finalizeFn)

	outputStream, outputErr = safeHandler.BuildStream()
	ctx := m.downstreamCtx
	if ctx == nil {
		ctx = context.Background()
	}

	outputStream.Write(stream.NewSimpleDatapack(ctx, m.pr))

	safeHandler.Start()

	return

}

func (m *Multiparter) datapackHandleFn(ctx context.Context, rc io.ReadCloser) error {

	if rc == nil {
		return nil
	}

	filename, err := m.extractFileName(ctx)
	if err != nil {
		return err
	}

	formFile, err := m.mw.CreateFormFile(m.fieldName, filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(formFile, rc)
	return err

}

func (m *Multiparter) finalizeFn() {
	m.mw.Close()
	m.pw.Close()
}

func (m *Multiparter) extractFileName(ctx context.Context) (string, error) {

	v := ctx.Value(m.ctxKeyOfFileName)
	if v == nil {
		return "", fmt.Errorf("cannot find file name with key '%s'", m.ctxKeyOfFileName)
	}

	filename, ok := v.(string)
	if !ok || len(strings.TrimSpace(filename)) == 0 {
		return "", fmt.Errorf("ctx value with key '%s' should be a string represents filename, but got type = %s, value = %v",
			m.ctxKeyOfFileName, reflect.TypeOf(v).Name(), v)
	}

	return filename, nil

}
