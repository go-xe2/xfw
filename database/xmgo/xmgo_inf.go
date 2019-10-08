package xmgo

import (
	"gopkg.in/mgo.v2"
)

type Ixmgo interface {
	// 获取数据库连接
	DB(db ...string) *mgo.Database
	// 打开连接
	Open(cfg ...interface{}) error
	// 关闭数据库连接
	Close() error
	// 配置mongo数据库连接
	Config(cfg interface{}) error
	// 获取数据库连接
	GetConfig() interface{}
}

// mongo collection集合名称
type IxmgoCollection interface {
	Name() string
}
