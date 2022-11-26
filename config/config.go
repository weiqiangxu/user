package config

import (
	"github.com/weiqiangxu/common-config/format"
)

var Conf Config

type Config struct {
	Application    format.NacosConfig `toml:"application" json:"application"`
	HttpConfig     format.HttpConfig  `toml:"http_config" json:"http_config"`
	UserGrpcConfig format.GrpcConfig  `toml:"grpc_config" json:"grpc_config"`
	LogConfig      format.LogConfig   `toml:"log_config" json:"log_config"`
	WikiMongoDb    format.MongoConfig `toml:"wiki_mongo_db" json:"wiki_mongo_db"`
	WikiRedisDb    format.RedisConfig `toml:"wiki_redis_db" json:"wiki_redis_db"`
	JwtConfig      jwt                `toml:"jwt_config" json:"jwt_config"`
}

type jwt struct {
	Secret  string `toml:"secret"`
	Timeout int64  `toml:"timeout"`
}
