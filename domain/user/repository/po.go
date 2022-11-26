package repository

import "code.skyhorn.net/backend/wiki-service/domain/user/entity"

// UserPo 仓储层模型，此时的updateTime是MySQL自带的字段模型对领域层而言是不需要的是无意义的
// 比如有一个字段模型是 is_deleted 是表示数据已经被删掉了，那么这个字段对于MySQL是可见的
// 对于领域层而言这个字段永远都不存在因为领域层面对的所有数据都是未删除的有效数据
type UserPo struct {
	Name       string
	UpdateTime int64
	UpdateUser int64
	Deleted    bool
}

// do 仓储层模型转领域层模型
// 此处的依赖而言是仓储层是可见领域层模型的，而领域层对持久化层模型不可见
// 就像我的仓储层可以是Elastic可以是MySQL可以是Redis但是领域层是不关注的
// 这样的好处就是我可以随时替换一个领域的数据存储方式，这个数据也可以是http可以是rpc可以是MySQL
// 分层解耦也是有利于协同开发
func (po *UserPo) do() *entity.UserDo {
	return &entity.UserDo{
		Name: po.Name,
	}
}
