package xredis

import (
	"encoding/json"
	"fmt"
	xfw "github.com/go-xe2/xfw/os"
	"github.com/go-xe2/xfw/util/xmap"
	"github.com/go-xe2/xfw/xerrors"
	"github.com/gogf/gf/g/encoding/gjson"
	"github.com/gogf/gf/g/encoding/gparser"
	"github.com/gogf/gf/g/os/gfile"
	"gopkg.in/redis.v4"
	"io/ioutil"
)

type xredisConfig struct {
	Host 		string		`json:"host"`
	Port		int 		`json:"port"`
	Password	string		`json:"password"`
	PoolSize	int 		`json:"poolSize"`
	DB 			int			`json:"db"`
}

const mapXredisConfigName = "xredis_config_name"

// 获取redis配置单一实例
func Config(name ...string) *xredisConfig {
	cfgName := mapXredisConfigName
	if len(name) > 0 {
		cfgName = name[0]
	}
	if c := xfw.GetInstance(cfgName); c != nil {
		return c.(*xredisConfig)
	}
	c := &xredisConfig{
		Host: 		"127.0.0.1",
		Port: 		6379,
		Password:   "",
		PoolSize: 	5,
		DB: 		0,
	}
	xfw.SetInstance(cfgName, c)
	return c
}

func (c *xredisConfig) Assign(src interface{}) {
	if cfg, ok := src.(*xredisConfig); ok {
		c.Host 		= cfg.Host
		c.Port 		= cfg.Port
		c.PoolSize 	= cfg.PoolSize
		c.DB		= cfg.DB
		c.Password	= cfg.Password
	} else if cfg, ok := src.(xredisConfig); ok {
		c.Host 		= cfg.Host
		c.Port 		= cfg.Port
		c.PoolSize 	= cfg.PoolSize
		c.DB		= cfg.DB
		c.Password	= cfg.Password
	}
}

// 从map解析
func (c *xredisConfig) ParseMap(cfg map[string]interface{}) error {
	c.Host 		= xmap.GetStrDef(cfg, "host", "127.0.0.1")
	c.Port 		= xmap.GetIntDef(cfg, "port", 6379)
	c.PoolSize 	= xmap.GetInt(cfg, "poolSize")
	c.DB 		= xmap.GetInt(cfg, "db")
	c.Password	= xmap.GetStr(cfg, "password")
	return nil
}

// 从文件解析, 支持json, xml, YAML and TOML
func (c *xredisConfig) LoadFromFile(fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return xerrors.New("配置文件[", fileName, "]读取错误")
	}
	parser, err := gparser.LoadContent(data)
	if err != nil {
		return xerrors.New("配置文件[", fileName, "]格式不支持")
	}
	return c.ParseMap(parser.ToMap())
}

// 从字符串解析
func (c *xredisConfig) ParseStr(cfg string) error {
	if gfile.IsFile(cfg) {
		return c.LoadFromFile(cfg)
	}
	var v xredisConfig
	err := json.Unmarshal([]byte(cfg), &v)
	if err != nil {
		return xerrors.New("不支持的配置格式:", err)
	}
	c.Assign(v)
	return nil
}

// 解析配置参数
// @s 可传入配置文件路径，json字符串或map[string]interface
func (c *xredisConfig) Parse(s interface{}) error {
	if cfg, ok := s.(map[string]interface{}); ok {
		return c.ParseMap(cfg)
	} else if cfg, ok := s.(string); ok {
		return c.ParseStr(cfg)
	} else {
		return xerrors.New("不支持的配置格式")
	}
}

// 配置转换成redis配置
func (c *xredisConfig) ToRedisOptions() *redis.Options {
	return &redis.Options{
		Addr: 		fmt.Sprintf("%s:%d", c.Host, c.Port),
		DB:   		c.DB,
		Password: 	c.Password,
		PoolSize: 	c.PoolSize,
	}
}

func (c *xredisConfig) String() string {
	b, err := gjson.Encode(c)
	if err != nil {
		return ""
	}
	return string(b)
}
