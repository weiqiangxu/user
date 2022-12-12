package router

import (
	"github.com/gin-gonic/gin"
	"github.com/weiqiangxu/user/application"
	"github.com/weiqiangxu/user/global/enum"
)

func RequestTracing() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span := application.App.Tracer.StartSpan(ctx.FullPath())
		ctx.Set(enum.TraceSpanName, span.Context())
		ctx.Next()
		span.Finish()
	}
}
