package http

import (
	"net/http"

	"github.com/weiqiangxu/user/application/front_service/dtos"

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
	info, _ := m.userDomainSrv.GetUserInfo(10)
	dto := &dtos.UserDto{
		Name: info.Name,
	}
	c.JSON(http.StatusOK, dto)
}
