package http

import (
	"github.com/weiqiangxu/common-config/logger"
	"net/http"
	"strings"
	"time"

	"github.com/weiqiangxu/user/application/front_service/dtos"

	"github.com/weiqiangxu/protocol/order"
	pbUser "github.com/weiqiangxu/protocol/user"

	"github.com/gin-gonic/gin"
	"github.com/weiqiangxu/user/domain/user"
)

type UserAppHttpOption func(service *UserAppHttpService)

func WithUserDomainService(t user.DomainInterface) UserAppHttpOption {
	return func(service *UserAppHttpService) {
		service.userDomainSrv = t
	}
}

func WithOrderRpcClient(t order.OrderClient) UserAppHttpOption {
	return func(service *UserAppHttpService) {
		service.orderRpcClient = t
	}
}

func WithUserRpcClient(t pbUser.LoginClient) UserAppHttpOption {
	return func(service *UserAppHttpService) {
		service.userRpcClient = t
	}
}

type UserAppHttpService struct {
	userDomainSrv  user.DomainInterface
	orderRpcClient order.OrderClient
	userRpcClient  pbUser.LoginClient
}

func NewUserAppHttpService(options ...UserAppHttpOption) *UserAppHttpService {
	srv := &UserAppHttpService{}
	for _, o := range options {
		o(srv)
	}
	return srv
}

// GetUserList get user list
func (m *UserAppHttpService) GetUserList(c *gin.Context) {
	// 如果没有下面这一段执行的太快了在list查看不到
	ch := make(chan bool)
	go func() {
		var stringSlice []string
		for i := 0; i < 20; i++ {
			// pprof 显示在这里占用2MB的内存开销
			repeat := strings.Repeat("hello,world", 50000)
			stringSlice = append(stringSlice, repeat)
			time.Sleep(time.Millisecond * 500)
		}
		logger.Info(len(stringSlice))
		ch <- true
	}()
	<-ch
	info, _ := m.userDomainSrv.GetUserInfo(10)
	dto := &dtos.UserDto{
		Name: info.Name,
	}
	c.JSON(http.StatusOK, dto)
}

// GetUserInfo get user info
func (m *UserAppHttpService) GetUserInfo(c *gin.Context) {
	response, err := m.userRpcClient.GetUserInfo(c.Request.Context(), &pbUser.GetUserInfoRequest{
		UniqueId: "1",
		NameMain: "2",
		NameSub:  "3",
	})
	if err != nil {
		c.JSON(http.StatusOK, err)
		return
	}
	c.JSON(http.StatusOK, response.UserInfo)
}
