package xmgo

import (
	"github.com/gogf/gf/g/encoding/gparser"
	"github.com/gogf/gf/g/os/gfile"
	xfw "gmny-server/xfw/os"
	"gmny-server/xfw/util/xmap"
	"gmny-server/xfw/xerrors"
	"gopkg.in/mgo.v2"
	"time"
)

const xmgoConfigName = "xmgo_config_name"


type xmgoConfig struct {
	Hosts		[]xmgoConfigHost	`json:"hosts"`
	Timeout		int64				`json:"timeout"`
	PoolLimit	int 				`json:"poolLimit"`
	DB 			string				`json:"db"`
}

// 获取mongo数据库唯一配置
func Config(name ...string) *xmgoConfig {
	cfgName := xmgoConfigName
	if len(name) > 0 {
		cfgName = name[0]
	}
	if c := xfw.GetInstance(cfgName); c != nil {
		return c.(*xmgoConfig)
	}

	c := &xmgoConfig{
		Hosts: []xmgoConfigHost{{
			Host:		"127.0.0.1",
			Port:		27017,
		}},
		Timeout: 	0,
		PoolLimit:	5,
		DB: 		"",
	}
	xfw.SetInstance(cfgName, c)
	return c
}

// 复制配置
func (c *xmgoConfig) Assign(src interface{}) {
	var hosts []xmgoConfigHost
	if cfg, ok := src.(*xmgoConfig); ok {
		c.Hosts 		= append(hosts, cfg.Hosts...)
		c.Timeout 		= cfg.Timeout
		c.PoolLimit 	= cfg.PoolLimit
		c.DB 			= cfg.DB
	} else if cfg, ok := src.(xmgoConfig); ok {
		c.Hosts 		= append(hosts, cfg.Hosts...)
		c.Timeout 		= cfg.Timeout
		c.PoolLimit 	= cfg.PoolLimit
		c.DB 			= cfg.DB
	}
}

func (c *xmgoConfig) ToMgoOptions() *mgo.DialInfo {
	var hosts []string
	for _, cfg := range c.Hosts {
		hosts = append(hosts, cfg.String())
	}
	return &mgo.DialInfo{
		Addrs: 		hosts,
		Direct: 	false,
		Timeout: 	time.Duration(c.Timeout) * time.Second,
		PoolLimit: 	c.PoolLimit,
	}
}

// 从map解析
func (c *xmgoConfig) ParseMap(cfg map[string]interface{}) error {
	c.Timeout	= xmap.GetInt64(cfg, "timeout")
	c.PoolLimit = xmap.GetIntDef(cfg, "poolLimit", 5)
	c.DB		= xmap.GetStr(cfg, "db")
	hosts, ok := cfg["hosts"].([]interface{})
	if !ok {
		return xerrors.New("hosts项配置不正确")
	}
	for i, item := range hosts {
		if cfg, ok := item.(map[string]interface{}); !ok {
			return xerrors.New("hosts第[", i,"]项不正确")
		} else {
			cfgItem := xmgoConfigHost{}
			cfgItem.ParseMap(cfg)
			c.Hosts = append(c.Hosts, cfgItem)
		}
	}
	return nil
}

// 从字符串解析
func (c *xmgoConfig) ParseStr(cfg string) error {
	if parser, err := gparser.LoadContent([]byte(cfg)); err != nil {
		return err
	} else {
		cfgMap := parser.ToMap()
		return c.ParseMap(cfgMap)
	}
}

// 从文件加载配置
func (c *xmgoConfig) LoadFromFile(fileName string) error {
	if !gfile.IsFile(fileName) {
		return xerrors.New("配置文件:", fileName, "路径无效")
	}
	parser, err := gparser.Load(fileName)
	if err != nil {
		return err
	}
	cfgMap := parser.ToMap()
	return c.ParseMap(cfgMap)
}

// 解析配置项
// @cfg 支持map[string]interface{}, 文件路径， json,xml, yaml, toml格式
func (c *xmgoConfig) Parse(cfg interface{}) error {
	switch value := cfg.(type) {
	case string:
		if gfile.IsFile(value) {
			return c.LoadFromFile(value)
		}
		return c.ParseStr(value)
	case map[string]interface{}:
		return c.ParseMap(value)
	case *xmgoConfig:
		c.Assign(value)
		return nil
	case xmgoConfig:
		c.Assign(value)
		return nil
	default:
		return xerrors.New("不支持的配置数据")
	}
}
