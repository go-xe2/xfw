package xredis

import "gopkg.in/redis.v4"

type IXredisClient interface {
	// 打开连接
	Open(cfg ...interface{}) error
	// 关闭连接
	Close() error
	// redis数据库客户端
	DB()  *redis.Client
	// 配置
	Config(cfg interface{}) error
	// 获取配置
	GetConfig() interface{}
}
