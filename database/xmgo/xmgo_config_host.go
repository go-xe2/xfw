package xmgo

import (
	"fmt"
	"gmny-server/xfw/util/xmap"
)

type xmgoConfigHost struct {
	Host	string		`json:"host"`
	Port	int			`json:"port"`
}



// 复制配置项
func (c *xmgoConfigHost) Assign(src interface{}) {
	if cfg, ok := src.(*xmgoConfigHost); ok {
		c.Host = cfg.Host
		c.Port = cfg.Port
	} else if cfg1 , ok := src.(xmgoConfigHost); ok {
		c.Host  = cfg1.Host
		c.Port	= cfg1.Port
	}
}

// 转换成字符串
func (c *xmgoConfigHost) String() string {
	return fmt.Sprintf("%v:%v", c.Host, c.Port)
}

// 从map解析
func (c *xmgoConfigHost) ParseMap(cfg map[string]interface{}) {
	c.Host = xmap.GetStrDef(cfg, "host", "127.0.0.1")
	c.Port = xmap.GetIntDef(cfg, "port", 27017)
}

// 转换成map
func (c *xmgoConfigHost) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	result["host"] = c.Host
	result["port"] = c.Port
	return result
}
