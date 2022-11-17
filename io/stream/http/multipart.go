package http

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"

	"github.com/SCU-SJL/sinfra/io/stream"
)

type MultipartUploader struct {
	mw *multipart.Writer
	pr *io.PipeReader
	pw *io.PipeWriter

	url       string
	headers   map[string]string
	cookies   []*http.Cookie
	fieldName string

	getFilename func(extra interface{}) string
}

func NewMultipartUploader() *MultipartUploader {
	uploader := &MultipartUploader{}
	uploader.pr, uploader.pw = io.Pipe()
	uploader.mw = multipart.NewWriter(uploader.pw)
	return uploader
}

func (u *MultipartUploader) SetGetFileNameFromExtraFn(fn func(extra interface{}) string) {
	u.getFilename = fn
}

// SafeUpload note that outputStream only contains one datapack.
func (u *MultipartUploader) SafeUpload(inputStream *stream.IOStream, inputErr *stream.ErrorPasser) (
	outputStream *stream.IOStream, outputErr *stream.ErrorPasser) {

	if u.pr == nil || u.pw == nil || len(strings.TrimSpace(u.url)) == 0 {
		panic("MultipartUploader is not initialized!")
	}

	req, err := u.buildHttpRequest()
	if err != nil {
		outputErr = stream.NewErrorPasser()
		outputErr.Put(err)
		return
	}

	safeHandler := stream.NewSafeIOStreamHandler(inputStream, inputErr, u.datapackHandleFn, u.finalizeFn)
	outputStream, outputErr = safeHandler.BuildStream()
	safeHandler.Start()

	go func() {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			outputErr.Put(err)
		}
		outputStream.Write(stream.NewSimpleDatapack(resp.Body, nil))
	}()

	return

}

func (u *MultipartUploader) buildHttpRequest() (*http.Request, error) {

	req, err := http.NewRequest(http.MethodPost, u.url, u.pr)
	if err != nil {
		return nil, err
	}

	for k, v := range u.headers {
		req.Header.Add(k, v)
	}
	req.Header.Add("Content-Type", u.mw.FormDataContentType())

	for i := range u.cookies {
		req.AddCookie(u.cookies[i])
	}

	return req, nil

}

func (u *MultipartUploader) datapackHandleFn(rc io.ReadCloser, extra interface{}) error {

	if rc == nil {
		return nil
	}

	if extra == nil {
		return errors.New("extra should be a string represents filename, but got nil")
	}

	var filename string

	if u.getFilename != nil {
		filename = u.getFilename(extra)
	} else {
		ok := false
		filename, ok = extra.(string)
		if !ok || len(strings.TrimSpace(filename)) == 0 {
			return fmt.Errorf("extra should be a string represents filename, but got type = %s, value = %v",
				reflect.TypeOf(extra).Name(), extra)
		}
	}

	formFile, err := u.mw.CreateFormFile(u.fieldName, extra.(string))
	if err != nil {
		return err
	}

	_, err = io.Copy(formFile, rc)
	if err != nil {
		return err
	}

	return u.mw.Close()

}

func (u *MultipartUploader) finalizeFn() {
	u.mw.Close()
	u.pw.Close()
}
