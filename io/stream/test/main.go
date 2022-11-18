package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/SCU-SJL/sinfra/io/stream"
	"github.com/SCU-SJL/sinfra/io/stream/archive"
	streamhttp "github.com/SCU-SJL/sinfra/io/stream/http"
)

const uploadURL = `http://localhost:8080/upload`

func main() {
	testBatchDownloadToZipToUpload()
}

func testBatchDownloadToZipToUpload() {

	// data src
	downloader := NewFileDownloader([]string{"foo.txt", "bar.txt"})
	iowriter := stream.NewSafeIOStreamWriter(downloader)

	// stream procs
	zipper := archive.NewZipper()
	zipper.SetDownstreamCtx(
		context.WithValue(context.Background(), "filename", "foobar.zip"),
	)
	zipper.SetCtxKeyOfFileName("filename")

	multiparter := streamhttp.NewMultiparter()
	multiparter.SetCtxKeyOfFileName("filename")

	proc := stream.BuildProcChain(zipper.Proc, multiparter.Proc)

	// start proc
	iostream, errp := iowriter.Start()
	outStream, outErr := proc(iostream, errp)

	for {
		data, closed := outStream.Read()
		if closed {
			break
		}
		rc := data.ReadCloser()
		uploadFile(rc, multiparter.FormDataContentType())
		rc.Close()
	}

	for {
		err, done := outErr.Check()
		if done {
			break
		}
		log.Println("out err:", err.Error())
	}

}

func testDownloadToZipToUpload() {

	ctx := context.Background()
	ctx = context.WithValue(ctx, "zip_filename", "zipped.txt")
	ctx = context.WithValue(ctx, "mp_filename", "mp.zip")

	zipper := archive.NewZipper()
	zipper.SetDownstreamCtx(ctx)
	zipper.SetCtxKeyOfFileName("zip_filename")

	multiparter := streamhttp.NewMultiparter()
	multiparter.SetCtxKeyOfFileName("mp_filename")

	streamProc := stream.BuildProcChain(zipper.Proc, multiparter.Proc)

	downloadReader := downloadFile("test.txt")

	outputStream, outputErr := streamProc(
		stream.NewClosedIOStream(stream.NewSimpleDatapack(ctx, downloadReader)),
		stream.NewClosedErrorPasser(),
	)

	datapack, _ := outputStream.Read()
	rc := datapack.ReadCloser()

	uploadFile(rc, multiparter.FormDataContentType())
	rc.Close()

	log.Println("output err =", outputErr.Get())

}

func downloadFile(filename string) io.ReadCloser {
	url := fmt.Sprintf("http://localhost:8080/download?filename=%s", filename)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("do http req failed: %v\n", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("do http req not ok, resp = %v\n", readHttpResp(resp))
	}
	return resp.Body
}

func uploadFile(r io.Reader, contentType string) {
	req, err := http.NewRequest(http.MethodPost, uploadURL, r)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", contentType)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	log.Println("upload resp:", readHttpResp(resp))
}

func readHttpResp(resp *http.Response) string {
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

type FileDownloader struct {
	idx      int
	fileList []string
}

func NewFileDownloader(fileList []string) *FileDownloader {
	return &FileDownloader{
		idx:      0,
		fileList: fileList,
	}
}

func (f *FileDownloader) Next() (datapack stream.Datapack, hasNext bool, err error) {
	filename := f.fileList[f.idx]
	ctx := context.Background()
	filepath := fmt.Sprintf("dir%d/%s", f.idx, filename)
	ctx = context.WithValue(ctx, "filename", filepath)
	rc := downloadFile(f.fileList[f.idx])
	f.idx++
	datapack = stream.NewSimpleDatapack(ctx, rc)
	hasNext = f.idx < len(f.fileList)
	return
}
