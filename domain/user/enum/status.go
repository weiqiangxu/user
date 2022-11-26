package enum

type UserDataStatus int

const (
	LogicDelete UserDataStatus = 1 // 逻辑删除掉的数据
	Normal      UserDataStatus = 2 // 正常的数据
)
