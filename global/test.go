package global

import "github.com/weiqiangxu/user/config"

// SetupDev dev环境配置注入 用于单元测试配置注入
func SetupDev() {
	config.Conf = config.Config{}
}

// SetupTesting 测试集群环境配置注入 用于单元测试配置注入
func SetupTesting() {
}
