package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sshelll/sinfra/gin/router/standard"
)

func main() {
	r := gin.Default()
	root := r.Group("/api")

	standard.SetSuccessCallback(func(gctx *gin.Context, resp interface{}) {
		gctx.JSON(200, gin.H{"data": resp})
	})
	standard.POST(root, "/user", handleUser, standard.DefaultRouteOpts().ParseHeader(true))
	standard.POST(root, "/car", handleCar)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

type userReq struct {
	Name *string `json:"name,omitempty" form:"name"`
	Age  *int    `json:"age,omitempty" form:"age"`
	Addr *string `json:"addr,omitempty" form:"addr"`
}

type userResp struct {
	Name *string `json:"name,omitempty"`
	Age  *int    `json:"age,omitempty"`
	Addr *string `json:"addr,omitempty"`
}

func handleUser(gctx *gin.Context, req *userReq) (rsp *userResp, err error) {
	return &userResp{
		Name: req.Name,
		Age:  req.Age,
		Addr: req.Addr,
	}, nil
}

type carReq struct {
	Name  *string `json:"name,omitempty" form:"name"`
	Brand *string `json:"brand,omitempty" form:"brand"`
}

type carResp struct {
	Name  *string `json:"name,omitempty"`
	Brand *string `json:"brand,omitempty"`
}

func handleCar(gctx *gin.Context, req *carReq) (rsp *carResp, err error) {
	return &carResp{
		Name:  req.Name,
		Brand: req.Brand,
	}, nil
}
