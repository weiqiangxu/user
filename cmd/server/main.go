package main

import (
	"github.com/weiqiangxu/common-config/format"
	"github.com/weiqiangxu/common-config/logger"
	"github.com/weiqiangxu/net"
	"github.com/weiqiangxu/net/transport"
	"github.com/weiqiangxu/net/transport/grpc"
	"github.com/weiqiangxu/protocol/user"
	"github.com/weiqiangxu/user/application"
	"github.com/weiqiangxu/user/config"
	"github.com/weiqiangxu/user/global/router"
)

func main() {
	// 配置依赖注入
	config.Conf = config.Config{
		Application:    config.AppInfo{Name: "server", Version: "v0.0.1"},
		UserGrpcConfig: format.GrpcConfig{Addr: ":9191"},
	}
	// mongodb && redis 等服务依赖
	application.Init()
	router.RegisterPrometheus()
	// 注册grpc服务
	grpcServer := grpc.NewServer(grpc.Address(config.Conf.UserGrpcConfig.Addr))
	user.RegisterLoginServer(grpcServer, application.App.AdminService.UserGrpcService)
	serverList := []transport.Server{grpcServer}
	if len(application.App.Event) > 0 {
		serverList = append(serverList, application.App.Event...)
	}
	// 将grpc && http 服务注入应用
	app := net.New(
		net.Name(config.Conf.Application.Name),
		net.Version(config.Conf.Application.Version),
		net.Server(serverList...),
	)
	if err := app.Run(); err != nil {
		logger.Fatal(err)
	}
}
