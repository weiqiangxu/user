package application

import (
	"time"

	redisApi "github.com/weiqiangxu/common-config/cache"
	"github.com/weiqiangxu/common-config/logger"
	"github.com/weiqiangxu/user/config"
	"github.com/weiqiangxu/user/domain/user"
	"github.com/weiqiangxu/user/infra/transport"

	adminGrpc "github.com/weiqiangxu/user/application/admin_service/grpc"
	"github.com/weiqiangxu/user/application/event"
	frontHttp "github.com/weiqiangxu/user/application/front_service/http"
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
	// 在这里将领域驱动对象注入
	redis := redisApi.NewRedisApi(config.Conf.WikiRedisDb)
	userDomain := user.NewUserService(user.WithRedis(redis))
	frontSrv := &frontService{}
	frontSrv.UserHttp = frontHttp.NewUserAppHttpService(
		frontHttp.WithUserDomainService(userDomain),
	)
	adminSrv := &adminService{}
	adminSrv.UserGrpcService = adminGrpc.NewUserAppGrpcService()

	// 将生成好的rpc客户端注入service那么service就可以使用这个rpc拉取数据
	//UserLogicGrpc, err := grpc.DialInsecure(context.Background(), grpc.WithEndpoint(config.Conf.UserGrpcConfig.Addr))
	//if err != nil {
	//	logger.Fatal(err)
	//}
	//userLoginGrpcClient := user.NewLoginClient(UserLogicGrpc)

	// 注入定时任务
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
