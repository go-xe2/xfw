package xredis

import (
	"fmt"
	"github.com/go-xe2/xfw/encoding/xjson"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/os/gcfg"
	"time"
)

var (
	c *gcfg.Config
	xredis IXredisClient
)
func init() {
	g.Config().SetFileName("config.json")
	c := g.Config()
	r, err := InitXredis(c)
	if err != nil {
		fmt.Printf("redis init fail:%v", err)
	}
	xredis = r
}

func ExampleXredisClient_GetConfig() {
	c := xredis.GetConfig()
	fmt.Printf("config:%v", xjson.ToString(c))

	// Output:
	// config:{"host":"127.0.0.1","port":6379,"password":"","poolSize":5,"db":0}

}

func ExampleSetAndGet() {
	var b bool
	var err error
	var result string

	fmt.Print("begin set value:\n")

	b, err = Set("str", "this is string")
	fmt.Printf("result:%v, err:%v\n", b, err)

	b, err = Set("bool", true)
	fmt.Printf("result:%v, err:%v\n", b, err)

	b, err = Set("int", 5566)
	fmt.Printf("result:%v, err:%v\n", b, err)

	b, err = Set("float", 5532.3)
	fmt.Printf("result:%v, err:%v\n", b, err)

	b, err = Set("map", map[string]interface{}{
		"name": "user1",
		"value": 323,
		"age": 34,
		"createdate": time.Now(),
	})
	fmt.Printf("result:%v, err:%v\n", b, err)

	fmt.Printf("begin get value:\n")

	result, err = Get("str")
	fmt.Printf("result:%v, err:%v\n", result, err)

	result, err = Get("bool")
	fmt.Printf("result:%v, err:%v\n", result, err)

	result, err = Get("int")
	fmt.Printf("result:%v, err:%v\n", result, err)

	result, err = Get("float")
	fmt.Printf("result:%v, err:%v\n", result, err)

	result, err = Get("map")
	fmt.Printf("result:%v, err:%v\n", result, err)

	var obj map[string]interface{}
	_, err = Get("map", &obj)
	fmt.Printf("get map with map[string]interface{} result:%v, err:%v\n", obj, err)


	// Output:


}
