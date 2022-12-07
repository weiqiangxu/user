package pprof_tool

import (
	netPprof "net/http/pprof"

	"github.com/gin-gonic/gin"
)

const (
	// DefaultPrefix url prefix of pprof
	DefaultPrefix = "/debug/pprof"
)

func getPrefix(prefixOptions ...string) string {
	prefix := DefaultPrefix
	if len(prefixOptions) > 0 {
		prefix = prefixOptions[0]
	}
	return prefix
}

// Register the standard Handler from the net/http/pprof package with
// the provided gin.Engine. prefixOptions is optional. If not prefixOptions,
// the default path prefix is used, otherwise first prefixOptions will be path prefix.
func Register(r *gin.Engine, prefixOptions ...string) {
	RouteRegister(&(r.RouterGroup), prefixOptions...)
}

// RouteRegister the standard Handler Func
func RouteRegister(route *gin.RouterGroup, prefixOptions ...string) {
	prefix := getPrefix(prefixOptions...)
	prefixRouter := route.Group(prefix)
	{
		prefixRouter.GET("/", gin.WrapF(netPprof.Index))
		prefixRouter.GET("/cmdline", gin.WrapF(netPprof.Cmdline))
		prefixRouter.GET("/profile", gin.WrapF(netPprof.Profile))
		prefixRouter.POST("/symbol", gin.WrapF(netPprof.Symbol))
		prefixRouter.GET("/symbol", gin.WrapF(netPprof.Symbol))
		prefixRouter.GET("/trace", gin.WrapF(netPprof.Trace))
		prefixRouter.GET("/allocs", gin.WrapH(netPprof.Handler("allocs")))
		prefixRouter.GET("/block", gin.WrapH(netPprof.Handler("block")))
		prefixRouter.GET("/goroutine", gin.WrapH(netPprof.Handler("goroutine")))
		prefixRouter.GET("/heap", gin.WrapH(netPprof.Handler("heap")))
		prefixRouter.GET("/mutex", gin.WrapH(netPprof.Handler("mutex")))
		prefixRouter.GET("/threadcreate", gin.WrapH(netPprof.Handler("threadcreate")))
	}
}
