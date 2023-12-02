package standard

import "github.com/gin-gonic/gin"

type RouteOpts struct {
	shouldBindQuery  bool
	shouldBindJSON   bool
	shouldBindXML    bool
	shouldBindYAML   bool
	shouldBindTOML   bool
	shouldBindHeader bool

	parseReqErrCallback func(gctx *gin.Context, err error)
	internalErrCallback func(gctx *gin.Context, err error)
	successCallback     func(gctx *gin.Context, resp interface{})
}

func DefaultRouteOpts() *RouteOpts {
	return &RouteOpts{
		shouldBindQuery: true,
		shouldBindJSON:  true,
	}
}

func OnlyParseJSONRouteOpts() *RouteOpts {
	return &RouteOpts{
		shouldBindJSON: true,
	}
}

func OnlyParseQueryRouteOpts() *RouteOpts {
	return &RouteOpts{
		shouldBindQuery: true,
	}
}

func (opts *RouteOpts) ParseQuery(b bool) *RouteOpts {
	opts.shouldBindQuery = b
	return opts
}

func (opts *RouteOpts) ParseJSON(b bool) *RouteOpts {
	opts.shouldBindJSON = b
	return opts
}

func (opts *RouteOpts) ParseXML(b bool) *RouteOpts {
	opts.shouldBindXML = b
	return opts
}

func (opts *RouteOpts) ParseYAML(b bool) *RouteOpts {
	opts.shouldBindYAML = b
	return opts
}

func (opts *RouteOpts) ParseTOML(b bool) *RouteOpts {
	opts.shouldBindTOML = b
	return opts
}

func (opts *RouteOpts) ParseHeader(b bool) *RouteOpts {
	opts.shouldBindHeader = b
	return opts
}

func (opts *RouteOpts) WithParseReqErrCallback(callback func(gctx *gin.Context, err error)) *RouteOpts {
	opts.parseReqErrCallback = callback
	return opts
}

func (opts *RouteOpts) WithInternalErrCallback(callback func(gctx *gin.Context, err error)) *RouteOpts {
	opts.internalErrCallback = callback
	return opts
}

func (opts *RouteOpts) WithSuccessCallback(callback func(gctx *gin.Context, resp interface{})) *RouteOpts {
	opts.successCallback = callback
	return opts
}
