package http

import (
	"code.skyhorn.net/backend/wiki-service/domain/user"
	"github.com/gin-gonic/gin"
	"github.com/weiqiangxu/common-config"
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
	common.ResponseSuccess(c, info)
}
