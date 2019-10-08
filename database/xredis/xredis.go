package xredis

import (
	xfw "gmny-server/xfw/os"
	"gopkg.in/redis.v4"
)

const mapXredisClientName = "xredis_client"

type xredisClient struct {
	config 		*xredisConfig
	client 		*redis.Client
}

func Client(name ...string) IXredisClient {
	clientName := mapXredisClientName
	if len(name) > 0 {
		clientName = name[0]
	}
	if c := xfw.GetInstance(clientName); c != nil {
		return c.(IXredisClient)
	}
	
	c := &xredisClient{
		config: Config(),
		client: nil,
	}
	xfw.SetInstance(clientName, c)
	return c
}


func (c *xredisClient) Open(cfg ...interface{}) error {
	if len(cfg) > 0 {
		err := c.config.Parse(cfg[0])
		if err != nil {
			return err
		}
	}
	options := c.config.ToRedisOptions()
	c.client = redis.NewClient(options)
	c.client.WrapProcess(func(old func(cmd redis.Cmder) error) func(cmd redis.Cmder) error {
		return func(cmd redis.Cmder) error {
			//glog.Info("redis start process:", cmd, "\n")
			err := old(cmd)
			//glog.Info("redis finished process:", cmd, "\n")
			return err
		}
	})
	return c.client.Ping().Err()
}

// 关闭连接
func (c *xredisClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// redis数据库客户端
func (c *xredisClient) DB()  *redis.Client {
	return c.client
}

// 配置
func (c *xredisClient) Config(cfg interface{}) error {
	return c.config.Parse(cfg)
}

// 获取配置
func (c *xredisClient) GetConfig() interface{} {
	return c.config
}


