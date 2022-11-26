package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/weiqiangxu/common-config/metrics"
	"github.com/weiqiangxu/user/application"
	"github.com/weiqiangxu/user/config"
)

func Init(r *gin.Engine) {
	monitorHandle := metrics.RequestMonitor()
	r.Use(monitorHandle)
	game := r.Group("/user")
	{
		game.GET("/list", application.App.FrontService.UserHttp.GetUserList)
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
