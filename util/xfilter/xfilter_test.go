package xfilter

import (
	"fmt"
	"time"
)

var (
	filterTag *FilterTag
)


var (
	rules = map[string]string{
		"name": "&string!required#登录名不能为空|登录名长度在:min到:max之间",
		"password": "&string!required|length:6,8|same:password2#密码不能为空|密码长度在:min到:max之间|两次密码输入不一致",
		"password2": "&string!required|length:6,8#",
		"email": "&string!required|email#邮箱不能为空|邮箱格式不正确",
	}

	arrRules = []string {
		"^name&string!required#登录名不能为空|登录名长度在:min到:max之间",
		"^password&string!required|length:6,8|same:password2#密码不能为空|密码长度在:min到:max之间|两次密码输入不一致",
		"^password2&string!required|length:6,8#",
		"^email&string!required|email#邮箱不能为空|邮箱格式不正确",
	}

	jsonRules = `
		{
			"name": {
				"dataType": "string",
				"rules":	"required",
				"validMsg":	"登录名不能为空|登录名长度在:min到:max之间"
			},
			"password": {
				"dataType": "string",
				"rules": "required|length:6,8|same:password2",
				"validMsg": "密码不能为空|密码长度在:min到:max之间|两次密码输入不一"
			},
			"password2": {
				"dataType": "string",
				"rules": "required|length:6,8"
			},
			"email": {
				"dataType": "string",
				"rules": "required|email",
				"validMsg": "邮箱不能为空|邮箱格式不正确"
			}
		}
	`

	testData = map[string]interface{} {
		"name": "user1",
		"password": "123456",
		"password2": "123456",
		"email": "ytx@dfdfd",
		"more1": "drop field1",
		"more2": "drop field2",
	}

)

type TestFilterStruct struct {
	Name		string		`filter:"^name!required#登录名不能为空|登录名长度在:min到:max之间"`
	Password 	string		`filter:"^password!required|length:6,8|same:password2#密码不能为空|密码长度在:min到:max之间|两次密码输入不一"`
	Password2 	time.Time	`filter:"^password!required|length:6,8"`
	Email		bool		`filter:"^email!required|email#邮箱不能为空|邮箱格式不正确"`
}



func ExampleNewTag() {
	filterTag = NewTag("^user_id&string~admin")
	fmt.Printf("new tag:%v", filterTag)

	// Output:
	// new tag:^user_id&string~admin

}

func ExampleFilterTag_Parse() {
	filterTag.Parse("^name&string~张三$author!required")
	fmt.Printf("parse:%v", filterTag)

	// Output:
	// parse:^name&string$author~张三!required

}


func ExampleFilterTag_String() {
	s := filterTag.String()
	fmt.Printf("string:%v", s)

	// Output:
	// string:^name&string$author~张三!required

}

func ExampleFilterTag_ToMap() {
	fmt.Printf("map:%v", filterTag.ToMap())

	// Output:
	// map:map[dataType:string defValue:张三 fieldName:name formatter:author rules:required]
}

func ExampleFilterTag_ToJson() {
	fmt.Printf("json:%v", filterTag.ToJson())

	// Output:
	//json:{
	//	"dataType": "string",
	//	"defValue": "张三",
	//	"fieldName": "name",
	//	"formatter": "author",
	//	"rules": "required"
	//}
}

func ExampleGetFilterFromMap() {
	rules := map[string]string{
		"name": 	"&string$author~!required",
		"age":		"&int~1",
		"birthday":	"&time$datetime~2019-01-01!required",
		"isAudit":	"&bool~false",
	}
	filters := GetFilterFromMap(rules)
	for i, item := range filters {
		fmt.Printf(fmt.Sprintf("%d:%s\n", i, item.ToJson()))
	}

	// Output:
	//0:{
	//	"dataType": "string",
	//	"defValue": "",
	//	"fieldName": "",
	//	"formatter": "author",
	//	"rules": "required"
	//}
	//1:{
	//	"dataType": "int",
	//	"defValue": "1",
	//	"fieldName": "",
	//	"formatter": "",
	//	"rules": ""
	//}
	//2:{
	//	"dataType": "time",
	//	"defValue": "2019-01-01",
	//	"fieldName": "",
	//	"formatter": "datetime",
	//	"rules": "required"
	//}
	//3:{
	//	"dataType": "bool",
	//	"defValue": "false",
	//	"fieldName": "",
	//	"formatter": "",
	//	"rules": ""
	//}

}

func ExampleGetFilterFromArray() {
	rules := []string{
		"^name&string$author~!required",
		"^age&int~1",
		"^birthday&time$datetime~2019-01-01!required",
		"^isAudit&bool~false",
	}
	filters := GetFilterFromArray(rules)
	for i, item := range filters {
		fmt.Printf(fmt.Sprintf("%d:%s\n", i, item.ToJson()))
	}

	// Output:
	//0:{
	//	"dataType": "string",
	//	"defValue": "",
	//	"fieldName": "name",
	//	"formatter": "author",
	//	"rules": "required"
	//}
	//1:{
	//	"dataType": "int",
	//	"defValue": "1",
	//	"fieldName": "age",
	//	"formatter": "",
	//	"rules": ""
	//}
	//2:{
	//	"dataType": "time",
	//	"defValue": "2019-01-01",
	//	"fieldName": "birthday",
	//	"formatter": "datetime",
	//	"rules": "required"
	//}
	//3:{
	//	"dataType": "bool",
	//	"defValue": "false",
	//	"fieldName": "isAudit",
	//	"formatter": "",
	//	"rules": ""
	//}
}


func ExampleGetFilterFromJson() {
	filters := GetFilterFromJson(jsonRules)
	if filters == nil {
		fmt.Printf("parse fail.")
	}
	for i, item := range filters {
		fmt.Printf(fmt.Sprintf("%d:%s\n", i, item.ToJson()))
	}


	// Output:
	//0:{
	//	"dataType": "bool",
	//	"defValue": "",
	//	"fieldName": "isAudit",
	//	"formatter": "",
	//	"rules": ""
	//}
	//1:{
	//	"dataType": "string",
	//	"defValue": "",
	//	"fieldName": "name",
	//	"formatter": "author",
	//	"rules": "required"
	//}
	//2:{
	//	"dataType": "int",
	//	"defValue": "",
	//	"fieldName": "age",
	//	"formatter": "",
	//	"rules": ""
	//}
	//3:{
	//	"dataType": "time",
	//	"defValue": "2019-01-01",
	//	"fieldName": "birthday",
	//	"formatter": "",
	//	"rules": "required"
	//}

}

func ExampleGetFilterFromStruct() {
	filters := GetFilterFromStruct(&TestFilterStruct{})
	if filters == nil {
		fmt.Printf("parse fail.")
	}
	for i, item := range filters {
		fmt.Printf(fmt.Sprintf("%d:%s\n", i, item.ToJson()))
	}


	// Output:
}



func ExampleFilterMap() {

	data, errs := FilterMap(testData, rules)
	fmt.Printf("filter1 result:%v, errors:%v\n", data, errs)

	data1, errs1 := FilterMap(testData, rules, map[string]interface{}{
		"name": 1,
		"email": 1,
	})
	fmt.Printf("filter2 result:%v, errors:%v\n", data1, errs1)

	// Output:

}

func ExampleFilter() {
	data, errs := Filter(testData, arrRules)
	fmt.Printf("filter3 result:%v, errors:%v\n", data, errs)

	data1, errs1 := Filter(testData, arrRules, map[string]interface{}{
		"name": 1,
		"email": 1,
	})
	fmt.Printf("filter4 result:%v, errors:%v\n", data1, errs1)

	// Output:
}

func ExampleFilterJson() {
	data, errs := FilterJson(testData, jsonRules)
	fmt.Printf("filter5 result:%v, errors:%v\n", data, errs)

	data1, errs1 := FilterJson(testData, jsonRules, map[string]interface{}{
		"name": 1,
		"email": 1,
	})
	fmt.Printf("filter6 result:%v, errors:%v\n", data1, errs1)

	// Output:
}

func ExampleFilterStruct() {
	data, errs := FilterStruct(testData, &TestFilterStruct{})
	fmt.Printf("filter7 result:%v, errors:%v\n", data, errs)

	data1, errs1 := FilterStruct(testData, &TestFilterStruct{}, map[string]interface{}{
		"name": 1,
		"email": 1,
	})
	fmt.Printf("filter8 result:%v, errors:%v\n", data1, errs1)

	// Output:
}