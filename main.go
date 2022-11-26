package main

import (
	"github.com/weiqiangxu/common-config/logger"
	"github.com/weiqiangxu/protocol/user"
	"github.com/weiqiangxu/user/application"
	"github.com/weiqiangxu/user/config"
	"github.com/weiqiangxu/user/global/router"
	"github.com/weiqiangxu/user/infra"
	"github.com/weiqiangxu/user/infra/transport"
	"github.com/weiqiangxu/user/infra/transport/grpc"
	"github.com/weiqiangxu/user/infra/transport/http"
)

func main() {
	// 这里对配置注入我们使用的nacos client读取数据注入配置
	config.Conf = config.Config{}
	application.Init()
	hs := http.NewServer(http.WithAddress(config.Conf.HttpConfig.ListenHTTP),
		http.WithPrometheus(config.Conf.HttpConfig.Prometheus),
		http.WithProfile(config.Conf.HttpConfig.Profile))

	// 启动Grpc服务将自身微服务实现的Grpc server端口暴露出去
	gs := grpc.NewServer(grpc.Address(config.Conf.UserGrpcConfig.Addr))
	srv := []transport.Server{hs, gs}
	if len(application.App.Event) > 0 {
		srv = append(srv, application.App.Event...)
	}
	app := infra.New(
		infra.Name(config.Conf.Application.DataId),
		infra.Version(config.Conf.Application.DataId),
		infra.Server(srv...),
	)
	// 注册自己实现的RPC服务
	user.RegisterLoginServer(gs, application.App.AdminService.UserGrpcService)
	router.Init(hs.Server())
	router.RegisterPrometheus()
	if err := app.Run(); err != nil {
		logger.Fatal(err)
	}
}
