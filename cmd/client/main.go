package main

import (
	"github.com/weiqiangxu/common-config/format"
	"github.com/weiqiangxu/common-config/logger"
	"github.com/weiqiangxu/net"
	"github.com/weiqiangxu/net/transport"
	"github.com/weiqiangxu/net/transport/http"
	"github.com/weiqiangxu/user/application"
	"github.com/weiqiangxu/user/config"
	"github.com/weiqiangxu/user/global/router"
)

func main() {
	// inject config from nacos
	config.Conf = config.Config{
		Application:     config.AppInfo{Name: "admin", Version: "v0.0.2"},
		HttpConfig:      format.HttpConfig{ListenHTTP: ":8989", Prometheus: true},
		UserGrpcConfig:  format.GrpcConfig{Addr: ":9191"},
		OrderGrpcConfig: format.GrpcConfig{},
		LogConfig:       format.LogConfig{},
		WikiMongoDb:     format.MongoConfig{},
		WikiRedisDb:     format.RedisConfig{},
		JwtConfig: config.JwtConfig{
			Secret:  "",
			Timeout: 0,
		},
		JaegerConfig: config.JaegerConfig{
			Addr: "http://127.0.0.1:14268/api/traces",
		},
	}
	application.Init()
	// register http server && grpc server
	httpServer := http.NewServer(http.WithAddress(config.Conf.HttpConfig.ListenHTTP),
		http.WithPrometheus(config.Conf.HttpConfig.Prometheus),
		http.WithProfile(config.Conf.HttpConfig.Profile))
	// mount routing and middleware to http server
	router.Init(httpServer.Server())
	router.RegisterPrometheus()
	// register http server && rpc server to gin engine and run
	serverList := []transport.Server{httpServer}
	if len(application.App.Event) > 0 {
		serverList = append(serverList, application.App.Event...)
	}
	app := net.New(
		net.Name(config.Conf.Application.Name),
		net.Version(config.Conf.Application.Version),
		net.Server(serverList...),
	)
	if err := app.Run(); err != nil {
		logger.Fatal(err)
	}
}
