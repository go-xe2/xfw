package xmysql

import (
	_ "github.com/go-sql-driver/mysql"
	xfw "github.com/go-xe2/xfw/os"
	"github.com/go-xe2/xorm"
	"github.com/gogf/gf/g/util/gconv"
)


const mapMysqlClientName = "default_mysql_client"

// mysql数据库连接客户端
type mysqlClient struct {
	config 	*mysqlConfig
	engin   *xorm.Engin
}

// mysql单一实例
func Mysql(name ...string) IMysqlClient {
	clientName := mapMysqlClientName
	if len(name) > 0 {
		clientName = gconv.String(name[0])
	}
	if c := xfw.GetInstance(clientName); c != nil {
		return c.(IMysqlClient)
	}
	c := &mysqlClient{
		config: Config(),
		engin: nil,
	}
	xfw.SetInstance(clientName, c)
	return c
}

func (s *mysqlClient) FullTableName(tableName string) string {
	return s.config.Prefix + tableName
}

// 通过map配置
func (s *mysqlClient) Config(c map[string]interface{}) error {
	err := s.config.ParseMap(c)
	return err
}

// 获取数据库配置
func (s *mysqlClient) GetConfig(configType ...MysqlConfigType) interface{} {
	return s.config.GetRoseConfig(configType...)
}

// 从文件中加载配置
func (s *mysqlClient) ConfigFromFile(fileName string) error {
	return s.config.LoadFromFile(fileName)
}

// 打开数据库连接
func (s *mysqlClient) Open(config ...interface{}) error {
	if len(config) > 0 {
		if c, ok := config[0].(map[string]interface{}); ok {
			s.config.ParseMap(c)
		} else if c, ok := config[0].(*mysqlConfig); ok {
			s.config.Assign(c)
		} else if c, ok := config[0].(mysqlConfig); ok {
			s.config.Assign(&c)
		}
	}
	var err error
	c := s.config.GetRoseConfig()
	if s.engin, err = xorm.Open(c); err != nil {
		return err
	}
	err = s.engin.Ping()
	return err
}

func (s *mysqlClient) Close() error {
	return nil
}

// 获取数据库操作接口，
// 是为了复用db对象, 可以在最后使用 db.LastSql() 获取最后执行的sql
// 如果不复用 db, 而是直接使用 DB(), 则会新建一个orm对象, 每一次都是全新的对象
// 所以复用 db, 一定要在当前会话周期内
func (s *mysqlClient) DB() xorm.IOrm {
	if s.engin == nil {
		panic("数据库连接未打开")
	}
	db := s.engin.NewOrm()
	db.Reset()
	return db
}




