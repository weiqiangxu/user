package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/weiqiangxu/common-config/metrics"
	"github.com/weiqiangxu/user/application"
	"github.com/weiqiangxu/user/config"
	"github.com/weiqiangxu/user/global/pprof_tool"
)

func Init(r *gin.Engine) {
	monitorHandle := metrics.RequestMonitor()
	r.Use(monitorHandle)
	r.Use(RequestTracing())
	pprof_tool.Register(r) // register pprof to gin
	game := r.Group("/user")
	{
		game.GET("/list", application.App.FrontService.UserHttp.GetUserList)
		game.GET("/info", application.App.FrontService.UserHttp.GetUserInfo)
		game.GET("/detail", application.App.FrontService.UserHttp.GetUserDetail)
	}
}

// RegisterPrometheus register prometheus
func RegisterPrometheus() {
	if !config.Conf.HttpConfig.Prometheus {
		return
	}
	prometheus.MustRegister(metrics.RequestLatencyHistogram)
	prometheus.MustRegister(metrics.RequestGauge)
}
