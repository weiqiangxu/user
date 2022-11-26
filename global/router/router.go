package router

import (
	commonMetrics "code.skyhorn.net/backend/common/metrics"
	ginMiddleware "code.skyhorn.net/backend/infra/gin-middleware"
	"code.skyhorn.net/backend/wiki-service/application"
	"code.skyhorn.net/backend/wiki-service/config"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func Init(r *gin.Engine) {
	jwtHandle := ginMiddleware.JWT([]byte(config.Conf.JwtConfig.Secret))
	monitorHandle := commonMetrics.RequestMonitor()
	r.Use(monitorHandle)
	game := r.Group("/user")
	{
		game.GET("/list", jwtHandle, application.App.FrontService.UserHttp.GetUserList)
	}
}

// RegisterPrometheus register prometheus
func RegisterPrometheus() {
	if !config.Conf.HttpConfig.Prometheus {
		return
	}
	prometheus.MustRegister(commonMetrics.RequestLatencyHistogram)
	prometheus.MustRegister(commonMetrics.RequestGauge)
}
