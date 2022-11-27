package application

import (
	"context"
	"time"

	redisApi "github.com/weiqiangxu/common-config/cache"
	"github.com/weiqiangxu/common-config/logger"
	"github.com/weiqiangxu/net/transport"
	"github.com/weiqiangxu/net/transport/grpc"
	"github.com/weiqiangxu/protocol/order"
	adminGrpc "github.com/weiqiangxu/user/application/admin_service/grpc"
	"github.com/weiqiangxu/user/application/event"
	frontHttp "github.com/weiqiangxu/user/application/front_service/http"
	"github.com/weiqiangxu/user/config"
	"github.com/weiqiangxu/user/domain/user"
)

var App app

type app struct {
	FrontService *frontService
	AdminService *adminService
	Event        []transport.Server
}

type frontService struct {
	UserHttp *frontHttp.UserAppHttpService
}

type adminService struct {
	UserGrpcService *adminGrpc.UserAppGrpcService
}

func Init() {
	// connect order rpc server to create order grpc client
	OrderGrpc, err := grpc.DialInsecure(context.Background(), grpc.WithEndpoint(config.Conf.OrderGrpcConfig.Addr), grpc.WithTracing(true))
	if err != nil {
		logger.Fatal(err)
	}
	orderGrpcClient := order.NewOrderClient(OrderGrpc)
	// inject rpc client && redis into domain service
	redis := redisApi.NewRedisApi(config.Conf.WikiRedisDb)
	userDomain := user.NewUserService(user.WithRedis(redis))
	frontSrv := &frontService{}
	frontSrv.UserHttp = frontHttp.NewUserAppHttpService(
		frontHttp.WithUserDomainService(userDomain),
		frontHttp.WithOrderRpcClient(orderGrpcClient),
	)
	adminSrv := &adminService{}
	adminSrv.UserGrpcService = adminGrpc.NewUserAppGrpcService()
	// inject cron event of match
	matchEvent := event.NewMatchEvent(
		event.WithTicker(time.NewTicker(time.Second*30)),
		event.WithMatchCronAction(func() error {
			logger.Info("start pull data from baidu")
			return nil
		}),
	)
	App = app{}
	App.FrontService = frontSrv
	App.AdminService = adminSrv
	App.Event = []transport.Server{matchEvent}
}
