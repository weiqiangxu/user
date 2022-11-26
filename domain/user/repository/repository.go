package repository

import (
	"context"

	"code.skyhorn.net/backend/wiki-service/domain/user/entity"

	"go.mongodb.org/mongo-driver/mongo"
)

// Repository 仓储层
type Repository interface {
	GetUserInfo(ctx context.Context, id int) (*entity.UserDo, error)
}

type UserRepositoryOption func(*userRepository)

func WithMongo(m *mongo.Database) UserRepositoryOption {
	return func(r *userRepository) {
		r.mongoDb = m
	}
}

func NewMatchRepository(os ...UserRepositoryOption) Repository {
	rep := &userRepository{}
	for _, o := range os {
		o(rep)
	}
	return rep
}

type userRepository struct {
	mongoDb *mongo.Database
}

// GetUserInfo 获取用户个人信息
func (r *userRepository) GetUserInfo(ctx context.Context, id int) (*entity.UserDo, error) {
	po := getDataFromMySQL(ctx, id)
	return po.do(), nil
}

// getDataFromMySQL 仓储层数据来源可以是redis可以是MySQL但是依赖必须是自己的私有属性
// 并且这个私有属性也是外部注入的
func getDataFromMySQL(ctx context.Context, id int) *UserPo {
	return &UserPo{}
}
