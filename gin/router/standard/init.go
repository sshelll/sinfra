package standard

import "github.com/gin-gonic/gin"

// Regularly, all the apis on the same server should return the same format
// and deal with the errors in the same way, so you only need to set these callbacks once.
//
// If you really want to set the callback for a single api,
// please use RouteOpts when you call standard.POST / GET / xxx to register handlers.

var (
	parseReqErrCallback = func(gctx *gin.Context, err error) {
		gctx.JSON(400, gin.H{"error": err.Error()})
	}

	internalErrCallback = func(gctx *gin.Context, err error) {
		gctx.JSON(500, gin.H{"error": err.Error()})
	}

	successCallback = func(gctx *gin.Context, resp interface{}) {
		gctx.JSON(200, gin.H{"data": resp})
	}
)

func SetParseReqErrCallback(handler func(gctx *gin.Context, err error)) {
	parseReqErrCallback = handler
}

func SetInternalErrCallback(handler func(gctx *gin.Context, err error)) {
	internalErrCallback = handler
}

func SetSuccessCallback(handler func(gctx *gin.Context, resp interface{})) {
	successCallback = handler
}
