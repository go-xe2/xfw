package xmgo

import (
	"github.com/go-xe2/xfw/xerrors"
	"github.com/gogf/gf/g/os/gcfg"
)

func InitXmgo(cfg *gcfg.Config) (Ixmgo, error) {
	mongoCfg := cfg.GetString("mongo")
	if mongoCfg == "" {
		return nil, xerrors.New("未配置mongo数据库连接参数")
	}
	xmgo := Xmgo()
	xmgo.Config(mongoCfg)
	err := xmgo.Open()
	if err != nil {
		return nil, err
	}
	return xmgo, nil
}

