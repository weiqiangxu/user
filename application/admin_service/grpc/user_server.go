package grpc

import (
	"context"

	"github.com/pkg/errors"
	redisApi "github.com/weiqiangxu/common-config/cache"
	"github.com/weiqiangxu/protocol/user"
)

type UserAppGrpcOption func(service *UserAppGrpcService)

func WithRedisApi(m *redisApi.RedisApi) UserAppGrpcOption {
	return func(srv *UserAppGrpcService) {
		srv.redis = m
	}
}

type UserAppGrpcService struct {
	*user.UnimplementedLoginServer
	redis *redisApi.RedisApi
}

func NewUserAppGrpcService(options ...UserAppGrpcOption) *UserAppGrpcService {
	srv := &UserAppGrpcService{}
	for _, o := range options {
		o(srv)
	}
	return srv
}

func (srv *UserAppGrpcService) GetUserInfo(ctx context.Context, request *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	return &user.GetUserInfoResponse{
		ErrorCode: user.ERROR_CODE_SuccessCode,
		UserInfo: &user.UserInfo{
			Name: "jack",
			Icon: "jack.icon",
			Age:  28,
		},
	}, nil
}

func (srv *UserAppGrpcService) DeleteUser(ctx context.Context, request *user.DeleteUserRequest) (*user.CommonReply, error) {
	return nil, nil
}
