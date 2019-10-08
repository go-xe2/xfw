package xmysql

import (
	"fmt"
	"github.com/go-xe2/xfw/encoding/xjson"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/os/gcfg"
)

var (
	c		*gcfg.Config
	mysql 	IMysqlClient
)

func init() {
	g.Config().SetFileName("config.json")
	c := g.Config()
	m, err := InitMySql(c)
	if err != nil {
		fmt.Printf("init mysql error:%v\n", err)
	}
	mysql = m
}

func ExampleMysqlClient_GetConfig() {
	c1 := mysql.GetConfig(MYSQL_CONFIG_TYPE_SINGLE)
	fmt.Printf("single config:%v\n", xjson.ToString(c1, "", "\t"))
	c2 := mysql.GetConfig(MYSQL_CONFIG_TYPE_CLUSTER)
	fmt.Printf("cluster config:%v\n", xjson.ToString(c2, "", "\t"))

	// Output:
//	single config:{
//	"Driver": "mysql",
//		"Dsn": "root:123456@tcp(127.0.0.1:3306)/mnyun?charset=utf8\u0026\u0026parseTime=True\u0026loc=Local",
//		"SetMaxOpenConns": 0,
//		"SetMaxIdleConns": 0,
//		"Prefix": "mny_"
//}
//	cluster config:{
//	"Master": [
//	{
//		"Driver": "mysql",
//		"Dsn": "root:123456@tcp(127.0.0.1:3306)/mnyun?charset=utf8\u0026\u0026parseTime=True\u0026loc=Local",
//		"SetMaxOpenConns": 0,
//		"SetMaxIdleConns": 0,
//		"Prefix": "mny_"
//	}
//	],
//	"Slave": [
//	{
//	"Driver": "mysql",
//	"Dsn": "root:123456@tcp(127.0.0.1:3306)/mnyun?charset=utf8\u0026\u0026parseTime=True\u0026loc=Local",
//	"SetMaxOpenConns": 0,
//	"SetMaxIdleConns": 0,
//	"Prefix": "mny_"
//	}
//	],
//	"Driver": "mysql",
//	"Prefix": "mny_"
//}


}

func ExampleSelect() {

	rows1, err1 := Select("admins", []string{"*"}, []interface{}{"user_id", 1}, "")
	fmt.Printf("rows:%v, err:%v\n", xjson.ToString(rows1, "", "\t"), err1)

	rows2, count, pageCount, err2 := SelectPage("admins", []string{"*"}, []interface{}{"user_id", ">", 0}, "", 0, 10)
	fmt.Printf("rows:%v, count:%v, pagecount:%v, err:%v\n", xjson.ToString(rows2, "", "\t"), count, pageCount, err2)


	rows3, err3 := SelectTop("admins", "*", []interface{}{"user_id", ">", 0}, "", 10)
	fmt.Printf("rows:%v, error:%v", xjson.ToString(rows3, "", "\t"), err3)

	// Output:

}

func ExampleExists() {
	b, err := Exists("admins", "user_id", 1)
	fmt.Printf("result:%v err:%v\n", b, err)

	// Output:

}
