package xredis

import (
	"github.com/go-xe2/xfw/xerrors"
	"github.com/gogf/gf/g/os/gcfg"
)

// redis 数据库初始化
func InitXredis (c *gcfg.Config) (IXredisClient, error) {
	s := c.GetString("redis")
	if s == "" {
		return nil, xerrors.New("redis数据库未配置")
	}
	clt := Client()
	err := clt.Open(s)
	if err != nil {
		return nil, err
	}
	return clt, nil
}
