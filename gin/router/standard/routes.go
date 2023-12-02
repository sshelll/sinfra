package standard

import (
	"github.com/gin-gonic/gin"
	"github.com/sshelll/sinfra/gin/router"
)

func Handle[Q, P any](r *gin.RouterGroup, method, relativePath string, handler router.RouteHandler[Q, P], routeOpts ...*RouteOpts) {
	r.Handle(method, relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}

func POST[Q, P any](r *gin.RouterGroup, relativePath string, handler router.RouteHandler[Q, P], routeOpts ...*RouteOpts) {
	r.POST(relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}

func GET[Q, P any](r *gin.RouterGroup, relativePath string, handler router.RouteHandler[Q, P], routeOpts ...*RouteOpts) {
	r.GET(relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}

func DELETE[Q, P any](r *gin.RouterGroup, relativePath string, handler router.RouteHandler[Q, P], routeOpts ...*RouteOpts) {
	r.DELETE(relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}

func PATCH[Q, P any](r *gin.RouterGroup, relativePath string, handler router.RouteHandler[Q, P], routeOpts ...*RouteOpts) {
	r.PATCH(relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}

func PUT[Q, P any](r *gin.RouterGroup, relativePath string, handler router.RouteHandler[Q, P], routeOpts ...*RouteOpts) {
	r.PUT(relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}

func OPTIONS[Q, P any](r *gin.RouterGroup, relativePath string, handler router.RouteHandler[Q, P], routeOpts ...*RouteOpts) {
	r.OPTIONS(relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}

func HEAD[Q, P any](r *gin.RouterGroup, relativePath string, handler router.RouteHandler[Q, P], routeOpts ...*RouteOpts) {
	r.HEAD(relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}

func Any[Q, P any](r *gin.RouterGroup, relativePath string, handler router.RouteHandler[Q, P], routeOpts ...*RouteOpts) {
	r.Any(relativePath, BuildGinHandlerFunc(handler, routeOpts...))
}
