package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weiqiangxu/user/domain/user"
)

type UserAppHttpOption func(service *UserAppHttpService)

func WithUserDomainService(t user.DomainInterface) UserAppHttpOption {
	return func(service *UserAppHttpService) {
		service.userDomainSrv = t
	}
}

type UserAppHttpService struct {
	userDomainSrv user.DomainInterface
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
	c.JSON(http.StatusOK, info)
}
