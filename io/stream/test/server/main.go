package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	register(r)
	r.Run(":8080")
}

func register(r *gin.Engine) {
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world")
	})
	r.GET("/download", handleDownload)
	r.POST("/upload", handleUpload)
}

func handleDownload(c *gin.Context) {
	filename := c.Query("filename")
	filepath := "./blob/" + filename
	if _, err := os.Stat(filepath); err != nil {
		internalErr(c, err)
		return
	}
	c.File(filepath)
}

func handleUpload(c *gin.Context) {
	mr, err := c.Request.MultipartReader()
	if err != nil {
		internalErr(c, err)
		return
	}

	part, err := mr.NextPart()
	if err != nil {
		internalErr(c, err)
		return
	}

	filename := part.FileName()

	file, err := os.Create("./blob/" + filename)
	if err != nil {
		internalErr(c, err)
		return
	}

	log.Println("start receiving file")
	_, err = io.Copy(file, part)
	if err != nil {
		internalErr(c, err)
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", filename))
}

func internalErr(c *gin.Context, err error) {
	log.Println("internal error:", err.Error())
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
}
