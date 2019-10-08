package xmgo

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/os/gcfg"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	c *gcfg.Config
	mgoDb Ixmgo
	testId = bson.NewObjectId()
)

func init() {
	g.Config().SetFileName("config.json")
	c := g.Config()
	r, err := InitXmgo(c)
	if err != nil {
		fmt.Printf("mgo init fail:%v\n", err)
	}
	mgoDb = r
}

// 删除集合
func ExampleDrop() {
	err := Drop("test")
	fmt.Printf("drop err:%v", err)

	// Output:

}

func ExampleInsert() {
	err := Insert("test", map[string]interface{}{
		"_id": testId,
		"name":"张三",
		"age":32,
		"birthday": time.Now(),
		"audit": true,
	})
	fmt.Printf("insert err:%v\n", err)

	// Output:

}

func ExampleFindOne() {
	one, err := FindOne("test", nil)
	fmt.Printf("findOne result:%v err:%v\n", one, err)

	// Output:

}

func ExampleUpdateById() {
	err := UpdateById("test", testId, bson.M{
		"name": "李四",
		"age": 18,
	})
	fmt.Printf("updateById result err:%v", err)

	// Output:

}

func ExampleFindById() {
	one, err := FindById("test",testId)
	fmt.Printf("findById result:%v err:%v\n", one, err)

	// Output:

}

func ExampleQuery() {
	query, err := Query("test", nil)
	if err != nil {
		fmt.Printf("query err:%v\n", err)
	}
	var rows []map[string]interface{}
	err = query.Select(map[string]interface{}{
		"name": 1,
		"age": 1,
	}).All(&rows)
	fmt.Printf("query rows:%v, err:%v\n", rows, err)

	// Output:

}