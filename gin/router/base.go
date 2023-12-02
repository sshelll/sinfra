package router

import "github.com/gin-gonic/gin"

type RouteHandler[Q, P any] func(gctx *gin.Context, req *Q) (resp *P, err error)
