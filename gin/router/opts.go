package router

import "github.com/gin-gonic/gin"

var (
	defaultRouteOpt    RouteOpt
	defaultRouteConfig *routeCfg
)

var (
	defaultParseReqErrCB = func(gctx *gin.Context, err error) {
		gctx.JSON(400, gin.H{"error": err.Error()})
	}

	defaultInternalErrCallback = func(gctx *gin.Context, err error) {
		gctx.JSON(500, gin.H{"error": err.Error()})
	}

	defaultSuccessCallback = func(gctx *gin.Context, resp any) {
		gctx.JSON(200, gin.H{"data": resp})
	}
)

type RouteOpt interface {
	apply(*routeCfg)
}

type routeOpt struct {
	fn func(*routeCfg)
}

func (opt *routeOpt) apply(cfg *routeCfg) {
	opt.fn(cfg)
}

type routeCfg struct {
	shouldBindQuery  bool
	shouldBindJSON   bool
	shouldBindXML    bool
	shouldBindYAML   bool
	shouldBindTOML   bool
	shouldBindHeader bool

	parseReqErrCallback func(gctx *gin.Context, err error)
	internalErrCallback func(gctx *gin.Context, err error)
	successCallback     func(gctx *gin.Context, resp any)
}

func init() {
	defaultRouteOpt = DefaultRouteOpt()
	defaultRouteConfig = &routeCfg{}
	defaultRouteOpt.apply(defaultRouteConfig)
}

func DefaultRouteOpt() RouteOpt {
	return &routeOpt{
		fn: func(cfg *routeCfg) {
			cfg.parseReqErrCallback = defaultParseReqErrCB
			cfg.internalErrCallback = defaultInternalErrCallback
			cfg.successCallback = defaultSuccessCallback
		},
	}
}

func BindQueryRouteOpt() RouteOpt {
	return &routeOpt{
		fn: func(cfg *routeCfg) {
			cfg.shouldBindQuery = true
		},
	}
}

func BindJSONRouteOpt() RouteOpt {
	return &routeOpt{
		fn: func(cfg *routeCfg) {
			cfg.shouldBindJSON = true
		},
	}
}

func BindXMLRouteOpt() RouteOpt {
	return &routeOpt{
		fn: func(cfg *routeCfg) {
			cfg.shouldBindXML = true
		},
	}
}

func BindYAMLRouteOpt() RouteOpt {
	return &routeOpt{
		fn: func(cfg *routeCfg) {
			cfg.shouldBindYAML = true
		},
	}
}

func BindTOMLRouteOpt() RouteOpt {
	return &routeOpt{
		fn: func(cfg *routeCfg) {
			cfg.shouldBindTOML = true
		},
	}
}

func BindHeaderRouteOpt() RouteOpt {
	return &routeOpt{
		fn: func(cfg *routeCfg) {
			cfg.shouldBindHeader = true
		},
	}
}

func ParseReqErrCallbackRouteOpt(callback func(*gin.Context, error)) RouteOpt {
	return &routeOpt{
		fn: func(cfg *routeCfg) {
			cfg.parseReqErrCallback = callback
		},
	}
}

func InternalErrCallbackRouteOpt(callback func(*gin.Context, error)) RouteOpt {
	return &routeOpt{
		fn: func(cfg *routeCfg) {
			cfg.internalErrCallback = callback
		},
	}
}

func SuccessCallbackRouteOpt(callback func(*gin.Context, any)) RouteOpt {
	return &routeOpt{
		fn: func(cfg *routeCfg) {
			cfg.successCallback = callback
		},
	}
}
