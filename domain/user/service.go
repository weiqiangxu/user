package user

import (
	"github.com/gogf/gf/util/gconv"
	redisApi "github.com/weiqiangxu/common-config/cache"
	"github.com/weiqiangxu/user/domain/user/entity"
	"github.com/weiqiangxu/user/domain/user/repository"
)

type DomainInterface interface {
	GetUserInfo(id int) (*entity.UserDo, error)
}

type MainLogicOption func(service *MainLogic)

func WithRepository(r repository.Repository) MainLogicOption {
	return func(srv *MainLogic) {
		srv.rep = r
	}
}

func WithRedis(r redisApi.RedisInterface) MainLogicOption {
	return func(srv *MainLogic) {
		srv.redis = r
	}
}

type MainLogic struct {
	rep   repository.Repository
	redis redisApi.RedisInterface
}

func NewUserService(options ...MainLogicOption) DomainInterface {
	srv := &MainLogic{}
	for _, o := range options {
		o(srv)
	}
	return srv
}

func (m *MainLogic) GetUserInfo(id int) (*entity.UserDo, error) {
	z := gconv.Map(`{"name":"jack","age":18}`)
	return &entity.UserDo{
		Name: gconv.String(z),
	}, nil
}
