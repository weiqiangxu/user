package http

import (
	"github.com/weiqiangxu/user/application/front_service/dtos"
	"net/http"
	"strings"
	"time"

	"github.com/weiqiangxu/protocol/order"

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

type UserAppHttpService struct {
	userDomainSrv  user.DomainInterface
	orderRpcClient order.OrderClient
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
		ch <- true
	}()
	<-ch
	info, _ := m.userDomainSrv.GetUserInfo(10)
	dto := &dtos.UserDto{
		Name: info.Name,
	}
	c.JSON(http.StatusOK, dto)
}
