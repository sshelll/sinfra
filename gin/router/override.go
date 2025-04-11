package router

import "github.com/gin-gonic/gin"

func Handle[Q, P any](r *gin.RouterGroup, method, relativePath string, handler RouteHandler[Q, P], routeOpts ...RouteOpt) {
	r.Handle(method, relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}

func POST[Q, P any](r *gin.RouterGroup, relativePath string, handler RouteHandler[Q, P], routeOpts ...RouteOpt) {
	r.POST(relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}

func GET[Q, P any](r *gin.RouterGroup, relativePath string, handler RouteHandler[Q, P], routeOpts ...RouteOpt) {
	r.GET(relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}

func DELETE[Q, P any](r *gin.RouterGroup, relativePath string, handler RouteHandler[Q, P], routeOpts ...RouteOpt) {
	r.DELETE(relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}

func PATCH[Q, P any](r *gin.RouterGroup, relativePath string, handler RouteHandler[Q, P], routeOpts ...RouteOpt) {
	r.PATCH(relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}

func PUT[Q, P any](r *gin.RouterGroup, relativePath string, handler RouteHandler[Q, P], routeOpts ...RouteOpt) {
	r.PUT(relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}

func OPTIONS[Q, P any](r *gin.RouterGroup, relativePath string, handler RouteHandler[Q, P], routeOpts ...RouteOpt) {
	r.OPTIONS(relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}

func HEAD[Q, P any](r *gin.RouterGroup, relativePath string, handler RouteHandler[Q, P], routeOpts ...RouteOpt) {
	r.HEAD(relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}

func Any[Q, P any](r *gin.RouterGroup, relativePath string, handler RouteHandler[Q, P], routeOpts ...RouteOpt) {
	r.Any(relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}
