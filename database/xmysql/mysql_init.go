package xmysql

import (
	"encoding/json"
	"github.com/gogf/gf/g/errors/gerror"
	"github.com/gogf/gf/g/os/gcfg"
	"github.com/gogf/gf/g/os/gfile"
)

// 初始化mysql数据库
func InitMySql(c *gcfg.Config) (IMysqlClient, error)  {
	mysql := Mysql()
	mysqlConfigFile := c.GetString("mysql")

	if mysqlConfigFile == "" {
		return nil, gerror.NewText("mysql数据库未配置")
	}
	if gfile.IsFile(mysqlConfigFile) {
		if err := mysql.ConfigFromFile(mysqlConfigFile); err != nil {
			return nil, err
		}
		return mysql, nil
	}
	var cmap map[string]interface{}
	err := json.Unmarshal([]byte(mysqlConfigFile), &cmap)
	if err != nil {
		return nil, err
	}
	if err := mysql.Config(cmap); err != nil {
		return nil, err
	}
	if err := mysql.Open(); err != nil {
		return nil, err
	}
	return mysql, nil
}

