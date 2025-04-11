package router

import (
	"io"

	"github.com/gin-gonic/gin"
)

type RouteHandler[Q, P any] func(gctx *gin.Context, req *Q) (resp *P, err error)

func BuildGinHandlerFunc[Q, P any](handler RouteHandler[Q, P], routeOpts ...RouteOpt) gin.HandlerFunc {
	cfg := defaultRouteConfig
	if len(routeOpts) > 0 {
		cfg = &routeCfg{}
		defaultRouteOpt.apply(cfg)
		for _, opt := range routeOpts {
			opt.apply(cfg)
		}
	}

	return func(gctx *gin.Context) {
		req := new(Q)
		if err := parseRequest(gctx, req, cfg); err != nil {
			cfg.parseReqErrCallback(gctx, err)
			return
		}

		resp, err := handler(gctx, req)
		if err != nil {
			cfg.internalErrCallback(gctx, err)
			return
		}

		cfg.successCallback(gctx, resp)
	}
}

func parseRequest[Q any](gctx *gin.Context, req *Q, cfg *routeCfg) error {
	if req == nil {
		panic("req is nil")
	}

	if cfg.shouldBindJSON {
		if err := gctx.ShouldBindJSON(req); err != nil && err != io.EOF {
			return err
		}
	}

	if cfg.shouldBindQuery {
		if err := gctx.ShouldBindQuery(req); err != nil && err != io.EOF {
			return err
		}
	}

	if cfg.shouldBindXML {
		if err := gctx.ShouldBindXML(req); err != nil && err != io.EOF {
			return err
		}
	}

	if cfg.shouldBindYAML {
		if err := gctx.ShouldBindYAML(req); err != nil && err != io.EOF {
			return err
		}
	}

	if cfg.shouldBindTOML {
		if err := gctx.ShouldBindTOML(req); err != nil && err != io.EOF {
			return err
		}
	}

	if cfg.shouldBindHeader {
		if err := gctx.ShouldBindHeader(req); err != nil && err != io.EOF {
			return err
		}
	}

	return nil
}
