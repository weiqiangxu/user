package main

import (
	"github.com/weiqiangxu/common-config/format"
	"github.com/weiqiangxu/common-config/logger"
	"github.com/weiqiangxu/net"
	"github.com/weiqiangxu/net/transport"
	"github.com/weiqiangxu/net/transport/grpc"
	"github.com/weiqiangxu/net/transport/http"
	"github.com/weiqiangxu/protocol/user"
	"github.com/weiqiangxu/user/application"
	"github.com/weiqiangxu/user/config"
	"github.com/weiqiangxu/user/global/router"
)

func main() {
	// inject config from nacos
	config.Conf = config.Config{
		Application: config.AppInfo{},
		HttpConfig: format.HttpConfig{
			ListenHTTP: ":8080",
			Profile:    false,
			Verbose:    false,
			Tracing:    false,
			Prometheus: false,
		},
		UserGrpcConfig: format.GrpcConfig{
			Addr: ":8989",
		},
		OrderGrpcConfig: format.GrpcConfig{},
		LogConfig:       format.LogConfig{},
		WikiMongoDb:     format.MongoConfig{},
		WikiRedisDb:     format.RedisConfig{},
		JwtConfig: config.JwtConfig{
			Secret:  "",
			Timeout: 0,
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
	// register user grpc server
	grpcServer := grpc.NewServer(grpc.Address(config.Conf.UserGrpcConfig.Addr))
	user.RegisterLoginServer(grpcServer, application.App.AdminService.UserGrpcService)
	// register http server && rpc server to gin engine and run
	serverList := []transport.Server{httpServer, grpcServer}
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
