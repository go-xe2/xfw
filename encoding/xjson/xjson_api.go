package xjson

import "encoding/json"

// 将数据转换成json字符串
// @data 要转换成json字符串的数据
// @options 格式化输出参数(可选)，参数1：prefix, 参数2: intent
func ToString(data interface{}, options ...string) string {
	var j []byte
	var err error
	if err != nil {
		return "{}"
	}
	if len(options) == 2 {
		j, err = json.MarshalIndent(data, options[0], options[1])
	} else {
		j, err = json.Marshal(data)
	}
	return string(j)
}

// 解析json字符串
func Parse (s string) interface{} {
	var v interface{}
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		return nil
	}
	return v
}
