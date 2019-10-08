package xmgo

import (
	xfw "github.com/go-xe2/xfw/os"
	"github.com/go-xe2/xfw/xerrors"
	"gopkg.in/mgo.v2"
)

type xmgo struct {
	config 			*xmgoConfig
	session 		*mgo.Session
}

const mapXmgoName = "xmgo_client"

// 获取mongo客户端单一实例
func Xmgo(name ...string) Ixmgo {
	xmgoName := mapXmgoName
	if len(name) > 0 {
		xmgoName = name[0]
	}

	if c := xfw.GetInstance(xmgoName); c != nil {
		return c.(Ixmgo)
	}

	c := &xmgo{
		config:Config(),
		session:nil,
	}

	xfw.SetInstance(xmgoName, c)
	return c
}

// 打开数据库连接
func (x *xmgo) Open(cfg ...interface{}) error {
	if len(cfg) > 0 {
		x.config.Parse(cfg[0])
	}
	options := x.config.ToMgoOptions()
	var err error
	x.session, err = mgo.DialWithInfo(options)
	if err != nil {
		return err
	}
	x.session.SetMode(mgo.Monotonic, true)
	if err := x.session.Ping(); err != nil {
		return xerrors.New("连接mongo数据库失败：", err)
	}
	return nil
}

// 关闭数据库连接
func (x *xmgo) Close() error {
	if x.session != nil {
		x.session.Close()
	}
	return nil
}

// 配置数据库连接
func (x *xmgo) Config(cfg interface{}) error {
	return x.config.Parse(cfg)
}

// 获取配置
func (x *xmgo) GetConfig() interface{} {
	return x.config
}

// 获取数据库操作会话
func (x *xmgo) DB(db ...string) *mgo.Database {
	if len(db) > 0 {
		return x.session.DB(db[0])
	}
	return x.session.DB(x.config.DB)
}
