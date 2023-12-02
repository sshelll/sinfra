package standard

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sshelll/sinfra/gin/router"
)

func BuildGinHandlerFunc[Q, P any](handler router.RouteHandler[Q, P], routeOpts ...*RouteOpts) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		var (
			parseReqErrCB = parseReqErrCallback
			internalErrCB = internalErrCallback
			successCB     = successCallback
		)

		if len(routeOpts) > 0 {
			opts := routeOpts[0]
			if opts.parseReqErrCallback != nil {
				parseReqErrCB = opts.parseReqErrCallback
			}
			if opts.internalErrCallback != nil {
				internalErrCB = opts.internalErrCallback
			}
			if opts.successCallback != nil {
				successCB = opts.successCallback
			}
		}

		req := new(Q)
		if err := ParseRequest(gctx, req, routeOpts...); err != nil {
			parseReqErrCB(gctx, err)
			return
		}

		resp, err := handler(gctx, req)
		if err != nil {
			internalErrCB(gctx, err)
			return
		}

		successCB(gctx, resp)
	}
}

func ParseRequest[Q any](gctx *gin.Context, req *Q, routeOpts ...*RouteOpts) error {
	if req == nil {
		panic("req is nil")
	}

	opts := &RouteOpts{
		shouldBindQuery: true,
		shouldBindJSON:  true,
	}
	if len(routeOpts) > 0 {
		opts = routeOpts[0]
	}

	if opts.shouldBindJSON {
		if err := gctx.ShouldBindJSON(req); err != nil && err != io.EOF {
			return err
		}
	}

	if opts.shouldBindQuery {
		if err := gctx.ShouldBindQuery(req); err != nil && err != io.EOF {
			return err
		}
	}

	if opts.shouldBindXML {
		if err := gctx.ShouldBindXML(req); err != nil && err != io.EOF {
			return err
		}
	}

	if opts.shouldBindYAML {
		if err := gctx.ShouldBindYAML(req); err != nil && err != io.EOF {
			return err
		}
	}

	if opts.shouldBindTOML {
		if err := gctx.ShouldBindTOML(req); err != nil && err != io.EOF {
			return err
		}
	}

	if opts.shouldBindHeader {
		if err := gctx.ShouldBindHeader(req); err != nil && err != io.EOF {
			return err
		}
	}

	return nil
}
