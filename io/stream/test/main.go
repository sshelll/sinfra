package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/SCU-SJL/sinfra/io/stream"
	"github.com/SCU-SJL/sinfra/io/stream/archive"
	streamhttp "github.com/SCU-SJL/sinfra/io/stream/http"
)

const uploadURL = `http://localhost:8080/upload`

func main() {
	testDownloadToZipToUpload()
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

	downloadReader := downloadFile()

	outputStream, outputErr := streamProc(
		stream.NewClosedIOStream(stream.NewSimpleDatapack(ctx, downloadReader)),
		stream.NewClosedErrorPasser(),
	)

	datapack, _ := outputStream.Read()
	rc := datapack.ReadCloser()

	// start uploading
	req, err := http.NewRequest(http.MethodPost, uploadURL, rc)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", multiparter.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	bs, _ := ioutil.ReadAll(resp.Body)
	log.Println("upload resp:", string(bs))

	log.Println("output err =", outputErr.Get())

}

func testDownload() {
	downloaded := downloadFile()
	f, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}
	io.Copy(f, downloaded)
	downloaded.Close()
}

func downloadFile() io.ReadCloser {
	filename := "test.txt"
	req, err := http.NewRequest(http.MethodGet, genDownloadURL(filename), nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("do http req failed: %v\n", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("do http req not ok, resp = %v\n", string(readHttpResp(resp)))
	}
	return resp.Body
}

func readHttpResp(resp *http.Response) []byte {
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return bs
}

func genDownloadURL(filename string) string {
	return fmt.Sprintf("http://localhost:8080/download?filename=%s", filename)
}
