package xmysql

import (
	"fmt"
	"github.com/gogf/gf/g/encoding/gparser"
	"github.com/gogf/gf/g/errors/gerror"
	"github.com/gogf/gf/g/util/gconv"
	"github.com/go-xe2/xorm"
	xfw "gmny-server/xfw/os"
	"gmny-server/xfw/util/xmap"
	"gmny-server/xfw/xerrors"
	"io/ioutil"
)

const mapMysqlClientConfigName = "default_mysql_config"

type mysqlConfigItem struct {
	Protocol 	string		`json:"protocol"`
	Charset	 	string		`json:"charset"`
	Host 	 	string		`json:"host"`
	Port	 	int			`json:"port"`
	User   	 	string		`json:"user"`
	Password 	string		`json:"password"`
	MaxOpenCons int			`json:"maxOpenCons"`
	MaxIdleCons int			`json:"maxIdleCons"`
}


type mysqlClusterConfig struct {
	Master  []mysqlConfigItem	`json:"master"`
	Slave   []mysqlConfigItem	`json:"slave"`
}

type mysqlConfig struct {
	Driver  string  			`json:"driver"` // 驱动
	DB 		string  			`json:"db"`		// 数据库
	Prefix	string				`json:"prefix"`	// 数据库前掇
	Config  string				`json:"config"`	// 使用的配置
	Single 	mysqlConfigItem		`json:"single"`
	Cluster mysqlClusterConfig	`json:"cluster"`
}

// 获取mysql数据库配置
func Config(name ...string) *mysqlConfig {
	cfgName := mapMysqlClientConfigName
	if len(name) > 0 {
		cfgName = gconv.String(name[0])
	}
	if c := xfw.GetInstance(cfgName); c != nil {
		return c.(*mysqlConfig)
	}
	c := &mysqlConfig{
		Driver: "mysql",
		DB: "mnyun",
		Prefix: "mny_",
		Config: "cluster",
		Single: mysqlConfigItem{
			Protocol:	"tcp",
			Charset:	"utf8",
			Host:		"127.0.0.1",
			Port:		3306,
			User:   	"root",
			Password:	"",
			MaxOpenCons: 0,
			MaxIdleCons: 3,
		},
		Cluster: mysqlClusterConfig{
			Master:[]mysqlConfigItem{{
				Protocol:	"tcp",
				Charset:	"utf8",
				Host:		"127.0.0.1",
				Port:		3306,
				User:   	"root",
				Password:	"",
				MaxOpenCons: 0,
				MaxIdleCons: 3,
			}},
			Slave:[]mysqlConfigItem{{
				Protocol:	"tcp",
				Charset:	"utf8",
				Host:		"127.0.0.1",
				Port:		3306,
				User:   	"root",
				Password:	"",
				MaxOpenCons: 0,
				MaxIdleCons: 3,
			}},
		},
	}
	xfw.SetInstance(cfgName, c)
	return c
}


// 转换为dsn
func (c *mysqlConfigItem) ToDsn(dbName string) string {
	return fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=%s&&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Protocol,
		c.Host,
		c.Port,
		dbName,
		c.Charset,
	)
}

func (c *mysqlConfigItem) ParseMap(cmap map[string]interface{}) error {
	c.User 			= xmap.GetStrDef(cmap, "user", "root")
	c.Password 		= xmap.GetStrDef(cmap, "password", "")
	c.Protocol 		= xmap.GetStrDef(cmap, "protocol", "tcp")
	c.Host 			= xmap.GetStrDef(cmap, "host", "127.0.0.1")
	c.Port 			= xmap.GetIntDef(cmap, "port", 3306)
	c.Charset 		= xmap.GetStrDef(cmap, "charset", "utf8")
	c.MaxIdleCons 	= xmap.GetIntDef(cmap, "maxOpenCons", 3)
	c.MaxOpenCons 	= xmap.GetIntDef(cmap, "maxOpenCons", 0)
	return nil
}

// 转换为rose单项配置项
func (c *mysqlConfigItem) ToRoseConfig(dbName string, driver, prefix string) *xorm.Config {
	return &xorm.Config{
		Dsn: 			 c.ToDsn(dbName),
		SetMaxOpenConns: c.MaxOpenCons,
		SetMaxIdleConns: c.MaxIdleCons,
		Driver: 		 driver,
		Prefix: 		 prefix,
	}
}

// 配置项转换成rose集群配置
func (c *mysqlClusterConfig) ToRoseConfig(dbName string, driver, prefix string) *xorm.ConfigCluster {
	var masters []xorm.Config
	var slaves  []xorm.Config
	for _, item := range c.Master {
		configItem := item.ToRoseConfig(dbName, driver, prefix)
		masters = append(masters, *configItem)
	}
	for _, item := range c.Slave {
		configItem := item.ToRoseConfig(dbName, driver, prefix)
		slaves = append(slaves, *configItem)
	}
	return &xorm.ConfigCluster{
		Master:  masters,
		Slave:   slaves,
		Driver:  driver,
		Prefix:  prefix,
	}
}


// 解析配置项
func (c *mysqlClusterConfig) ParseMap (m map[string]interface{}) error {
	masterCfgItems, ok := m["master"].([]interface{})
	if !ok {
		return gerror.NewText("master配置项不正确")
	}
	var masterItems []mysqlConfigItem
	for i, c := range masterCfgItems {
		item := mysqlConfigItem{}
		itemCfg, ok := c.(map[string]interface{})
		if !ok {
			return xerrors.New("cluster.master中第%d项配置不正确", i)
		}
		if err := item.ParseMap(itemCfg); err == nil {
			masterItems = append(masterItems, item)
		}
	}
	c.Master = masterItems
	var slaveItems []mysqlConfigItem
	if slaveConfigItems, ok := m["slave"].([]interface{}); !ok {
		return nil
	} else {
		for i, c := range slaveConfigItems {
			item := mysqlConfigItem{}
			itemCfg, ok := c.(map[string]interface{})
			if !ok {
				return xerrors.New("cluster.slave中第%d项配置不正确", i)
			}
			if err := item.ParseMap(itemCfg); err == nil {
				slaveItems = append(slaveItems, item)
			}
		}
	}
	c.Slave = slaveItems
	return nil
}


// 获取rose配置
func (c *mysqlConfig) GetRoseConfig(configType ...MysqlConfigType) interface{} {
	if len(configType) > 0 {
		if configType[0] == MYSQL_CONFIG_TYPE_SINGLE {
			return c.Single.ToRoseConfig(c.DB, c.Driver, c.Prefix)
		}
		return c.Cluster.ToRoseConfig(c.DB, c.Driver, c.Prefix)
	}
	if c.Config == "cluster" {
		return c.Cluster.ToRoseConfig(c.DB, c.Driver, c.Prefix)
	} else {
		return c.Single.ToRoseConfig(c.DB, c.Driver, c.Prefix)
	}
}

func (c *mysqlConfig) Assign(src interface{}) {
	//Driver  string  // 驱动
	//DB 		string  // 数据库
	//Prefix	string	// 数据库前掇
	//Config  string	// 使用的配置
	if cfg, ok := src.(*mysqlConfig); ok {
		c.Driver 	= cfg.Driver
		c.DB 		= cfg.DB
		c.Prefix 	= cfg.Prefix
		c.Config 	= cfg.Config
		c.Single 	= cfg.Single
		c.Cluster 	= cfg.Cluster
	} else if cfg, ok := src.(mysqlConfig); ok {
		c.Driver 	= cfg.Driver
		c.DB 		= cfg.DB
		c.Prefix 	= cfg.Prefix
		c.Config 	= cfg.Config
		c.Single 	= cfg.Single
		c.Cluster 	= cfg.Cluster
	}
}

// 从文件读取配置
// 支持 JSON, XML, YAML and TOML格式
func (c *mysqlConfig) LoadFromFile (fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return gerror.NewText(fmt.Sprintf("读取数据库配置出错:%v", err))
	}
	parser, err := gparser.LoadContent(data)
	if err != nil {
		return err
	}
	return c.ParseMap(parser.ToMap())
}

// 从map解析配置项
func (c *mysqlConfig) ParseMap (cmap map[string]interface{}) error {
	var config = new(mysqlConfig)
	var ok bool
	if config.Driver, ok = cmap["driver"].(string); !ok {
		return gerror.NewText("缺少driver项")
	}
	if config.DB, ok = cmap["db"].(string); !ok {
		return gerror.NewText("缺少db配置项")
	}
	if config.Prefix, ok = cmap["prefix"].(string); !ok {
		config.Prefix = ""
	}
	if config.Config, ok = cmap["config"].(string); !ok {
		config.Config = "single"
	}

	if singleCfg, ok := cmap["single"].(map[string]interface{}); !ok && config.Config == "single" {
		return gerror.NewText("mysql未配置single节点")
	} else {
		if err := config.Single.ParseMap(singleCfg); err != nil {
			return err
		}
	}

	if clusterCfg, ok := cmap["cluster"].(map[string]interface{}); !ok && config.Config == "cluster" {
		return gerror.NewText("mysql未配置cluster节点")
	} else {
		if err := config.Cluster.ParseMap(clusterCfg); err != nil {
			return err
		}

	}
	c.Assign(config)
	return nil
}


