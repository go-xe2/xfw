package xmysql

import "github.com/go-xe2/xorm"


type MysqlConfigType string

// 单数据库配置类型
const MYSQL_CONFIG_TYPE_SINGLE 		MysqlConfigType = "single"

// 数据库集群配置类型
const MYSQL_CONFIG_TYPE_CLUSTER 	MysqlConfigType = "cluster"

// mysql客户端接口
type IMysqlClient interface {
	// 配置数据库连接
	Config(c map[string]interface{}) error
	// 从文件加载数据库配置
	ConfigFromFile(fileName string) error
	// 连接数据库
	Open(cfg ...interface{}) error
	// 关闭数据库连接
	Close() error
	// 获取数据库连接操作会话
	DB() xorm.IOrm
	// 获取数据库连接配置
	GetConfig(configType ...MysqlConfigType) interface{}
	FullTableName(tableName string) string
}
