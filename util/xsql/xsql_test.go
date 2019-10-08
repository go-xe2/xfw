package xsql

import "testing"


var ModelFields = map[string]interface{} {
	"user_id": 		"user_id",
	"login_name": 	"login_name",
	"nick_name": 	"nick_name",
	"pwd": 			"-",
	"enc": 			"-",
	"sex": 			[]interface{}{ "S",`case sex when 0 then "未知" when 1 then "男" when 2 then "女" end`, "sex_name" },
	"mobile": 		"mobile",
	"qq": 			"qq",
	"province": 	"province",
	"city": 		"city",
	"county": 		"county",
	"town": 		"town",
	"province_id": 	"-",
	"city_id": 		"-",
	"county_id": 	"-",
	"town_id": 		"-",
	"address": 		"address",
	"head": 		"head",
	"cr_date": 		[]interface{}{ "ML,UL,GL,MD", "from_unixtime(cr_date)" },
	"up_date": 		[]interface{}{ "ML,UL,GL,MD", "from_unixtime(up_date)" },
	"last_login": 	[]interface{}{ "ML,MD", "from_unixtime(last_login)" },
	"last_ip": 		[]interface{}{ "ML,MD" },
	"status": 		[]interface{}{ "ML,MD", `case status when 0 then "未审核" when 1 then "审核失败" when 2 then "审核通过" when 3 then "锁定" end`, "status_name" },
	"fav_count": 	"fav_count",
	"visit_count": 	"visit_count",
}

func TestBuildSqlField(t *testing.T) {
	s1 := BuildSqlField(ModelFields, "MD")
	t.Log("MD:", s1)

	s2 := BuildSqlField(ModelFields, "ML")
	t.Log("ML", s2)

	s3 := BuildSqlField(ModelFields, "UD")
	t.Log("UD", s3)

	s4 := BuildSqlField(ModelFields, "UL")
	t.Log("UL", s4)

	s5 := BuildSqlField(ModelFields, "GD")
	t.Log("GD", s5)

	s6 := BuildSqlField(ModelFields, "GL")
	t.Log("GL", s6)

	s7 := BuildSqlField(ModelFields, "MD", "u")
	t.Log("MD as u:", s7)
}

func TestBuildSqlFieldBySelect(t *testing.T) {
	s1 := BuildSqlFieldBySelect(ModelFields, "user_id,pwd, enc,sex,up_date")
	t.Log("select fields:", s1)
}

func TestFieldNameValid(t *testing.T) {
	b := CheckFieldNameValid(ModelFields, "sex")
	t.Log("b:", b)

	b1 := CheckFieldNameValid(ModelFields, "b1")
	t.Log("b1:", b1)
}

func TestBuildSqlAllFieldArray(t *testing.T) {
	arr := BuildSqlAllFieldArray(ModelFields)
	t.Log("all fields arr:", arr)
}

func TestBuildSqlAllField(t *testing.T) {
	s1 := BuildSqlAllField(ModelFields)
	t.Log("all fields:", s1)
}
