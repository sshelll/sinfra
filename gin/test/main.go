package main

import (
	"github.com/gin-gonic/gin"
	infraGin "github.com/sshelll/sinfra/gin"
)

func main() {
	r := gin.Default()
	r.Any("/ping", handlePing())
	r.POST("/test", handleTest())
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func handlePing() gin.HandlerFunc {
	type req struct {
		Msg *string `json:"msg"`
	}
	type resp struct {
		Msg *string `json:"msg,omitempty"`
	}
	return infraGin.MakeGinHandlerFunc(
		func(gctx *gin.Context, req *req) (rsp *resp, err error) {
			return &resp{Msg: req.Msg}, nil
		},
	)
}

func handleTest() gin.HandlerFunc {
	type req struct {
		Name *string `json:"name,omitempty"`
		Age  *int    `json:"age,omitempty"`
		Addr *string `json:"addr,omitempty"`
	}
	type resp struct {
		Name *string `json:"name,omitempty"`
		Age  *int    `json:"age,omitempty"`
		Addr *string `json:"addr,omitempty"`
	}
	return infraGin.MakeGinHandlerFunc(
		func(gctx *gin.Context, req *req) (rsp *resp, err error) {
			return &resp{
				Name: req.Name,
				Age:  req.Age,
				Addr: req.Addr,
			}, nil
		},
	)
}
