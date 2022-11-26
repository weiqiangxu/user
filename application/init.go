package application

import (
	"code.skyhorn.net/backend/wiki-service/config"
	"code.skyhorn.net/backend/wiki-service/domain/user"
	redisapi "github.com/weiqiangxu/common-config/cache"
	"time"

	"code.skyhorn.net/backend/infra/gms/transport"
	"code.skyhorn.net/backend/infra/logger"
	admin_grpc "code.skyhorn.net/backend/wiki-service/application/admin_service/grpc"
	"code.skyhorn.net/backend/wiki-service/application/event"
	front_http "code.skyhorn.net/backend/wiki-service/application/front_service/http"
)

var App app

type app struct {
	FrontService *frontService
	AdminService *adminService
	Event        []transport.Server
}

type frontService struct {
	UserHttp *front_http.UserAppHttpService
}

type adminService struct {
	UserGrpcService *admin_grpc.UserAppGrpcService
}

func Init() {
	// 在这里将领域驱动对象注入
	redis := redisapi.NewRedisApi(config.Conf.WikiRedisDb)
	userDomain := user.NewUserService(user.WithRedis(redis))
	frontSrv := &frontService{}
	frontSrv.UserHttp = front_http.NewUserAppHttpService(
		front_http.WithUserDomainService(userDomain),
	)
	adminSrv := &adminService{}
	adminSrv.UserGrpcService = admin_grpc.NewUserAppGrpcService()

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
