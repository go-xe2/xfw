package xfilter

import (
	"bytes"
	"github.com/go-xe2/xfw/encoding/xjson"
	"github.com/go-xe2/xfw/util/xmap"
	"github.com/gogf/gf/g/util/gconv"
	"reflect"
)

const (
	// 字段名前掇
	x_TAG_FIELD_NAME_TAG 	= '^'
	// 数据类型前掇
	x_TAG_DATA_TYPE_TAG		= '&'
	// 数据格式化类名前掇
	x_TAG_FORMATTER_TAG		= '$'
	// 默认值前掇
	x_TAG_DEFAULT_VALUE_TAG	= '~'
	// 验证规则前掇
	x_TAG_RULES_TAG			= '!'
	x_TAG_VALID_MSG_TAG		= '#'
	x_TAG_PARAM_NAME_TAG	= '?'
)

var x_TAG_DATA_TYPES = map[string]reflect.Kind {
	"string":		reflect.String,
	"int":			reflect.Int64,
	"float":		reflect.Float64,
	"bool":			reflect.Bool,
	"time":			reflect.Uint64,
}

type FieldName string
type FieldDataType string

// 字段过滤器使用的tag
type FilterTag struct {
	fieldName	FieldName		// 数据库字段名称
	paramName	string			// 参数名称
	// 数据类型
	dataType 	FieldDataType
	formatter   string			// 格式化工具
	// 默认值
	defValue 	interface{}
	rules 		string			// 验证规则
	validMsg	string			// 验证失败消息
}

// 创建过滤器
func NewTag(structTag ...string) *FilterTag {
	filter := &FilterTag{
		fieldName: 		"",
		paramName:		"",
		dataType: 		"string",
		formatter: 		"",
		defValue: 		nil,
		rules:			"",
		validMsg:		"",
	}
	if len(structTag) > 0 {
		filter.Parse(structTag[0])
	}
	return filter
}

func (f *FilterTag) FromJson(v map[string]interface{}) {
	f.fieldName = FieldName(xmap.GetStr(v, "fieldName"))
	f.dataType	= FieldDataType(xmap.GetStr(v, "dataType"))
	f.formatter	= xmap.GetStr(v, "formatter")
	if data, ok := v["defValue"]; ok {
		switch f.dataType {
		case "string":
			f.defValue = gconv.String(data)
			break
		case "int":
			f.defValue = gconv.Int64(data)
			break
		case "float":
			f.defValue = gconv.Float64(data)
			break
		case "time":
			f.defValue = gconv.Time(data)
			break
		case "bool":
			f.defValue = gconv.Bool(data)
			break
		default:
			f.defValue = gconv.String(data)
			f.dataType = "string"
		}
	}
	f.paramName	= xmap.GetStr(v, "paramName")
	f.rules		= xmap.GetStr(v, "rules")
	f.validMsg	= xmap.GetStr(v, "validMsg")
}

func (f *FilterTag) Parse(structTag string) {
	if structTag == "" {
		return
	}
	var valBuf bytes.Buffer
	var lastChar byte = 0
	var process = func() {
		if lastChar == 0 {
			return
		}
		switch lastChar {
		case x_TAG_FIELD_NAME_TAG:
			f.fieldName = FieldName(valBuf.String())
			valBuf.Reset()
			break
		case x_TAG_DATA_TYPE_TAG:
			f.dataType = FieldDataType(valBuf.String())
			valBuf.Reset()
			break
		case x_TAG_FORMATTER_TAG:
			f.formatter = valBuf.String()
			valBuf.Reset()
			break
		case x_TAG_DEFAULT_VALUE_TAG:
			f.defValue = valBuf.String()
			valBuf.Reset()
			break
		case x_TAG_RULES_TAG:
			f.rules = valBuf.String()
			valBuf.Reset()
			break
		case x_TAG_VALID_MSG_TAG:
			f.validMsg = valBuf.String()
			valBuf.Reset()
			break
		case x_TAG_PARAM_NAME_TAG:
			f.paramName = valBuf.String()
			valBuf.Reset()
			break
		}
	}
	for i := 0; i < len(structTag); i++ {
		c := structTag[i]
		switch c {
		case x_TAG_FIELD_NAME_TAG:
			process()
			lastChar = x_TAG_FIELD_NAME_TAG
			break
		case x_TAG_DATA_TYPE_TAG:
			process()
			lastChar = x_TAG_DATA_TYPE_TAG
			break
		case x_TAG_FORMATTER_TAG:
			process()
			lastChar = x_TAG_FORMATTER_TAG
			break
		case x_TAG_DEFAULT_VALUE_TAG:
			process()
			lastChar = x_TAG_DEFAULT_VALUE_TAG
			break
		case x_TAG_RULES_TAG:
			process()
			lastChar = x_TAG_RULES_TAG
			break
		case x_TAG_VALID_MSG_TAG:
			process()
			lastChar = x_TAG_VALID_MSG_TAG
			break
		case x_TAG_PARAM_NAME_TAG:
			process()
			lastChar = x_TAG_PARAM_NAME_TAG
			break
		default:
			valBuf.WriteByte(c)
		}
	}
	process()
}

// 转换成tag字符串
// @字段名 #数据类型 $格式化工具名 ~默认值 !验证规则
// @fieldName#dataType$formatter~defValue!rules
func (f *FilterTag) String() string {
	var buf bytes.Buffer
	if f.fieldName != "" {
		buf.WriteString(string(x_TAG_FIELD_NAME_TAG))
		buf.WriteString(string(f.fieldName))
	}
	if f.dataType != "" {
		buf.WriteString(string(x_TAG_DATA_TYPE_TAG))
		buf.WriteString(string(f.dataType))
	}
	if f.formatter != "" {
		buf.WriteString(string(x_TAG_FORMATTER_TAG))
		buf.WriteString(f.formatter)
	}
	if f.defValue != nil {
		buf.WriteString(string(x_TAG_DEFAULT_VALUE_TAG))
		buf.WriteString(gconv.String(f.defValue))
	}
	if f.rules != "" {
		buf.WriteString(string(x_TAG_RULES_TAG))
		buf.WriteString(f.rules)
	}
	if f.validMsg != "" {
		buf.WriteString(string(x_TAG_VALID_MSG_TAG))
		buf.WriteString(f.validMsg)
	}
	if f.paramName != "" {
		buf.WriteString(string(x_TAG_PARAM_NAME_TAG))
		buf.WriteString(f.paramName)
	}
	return buf.String()
}

// 转换成验证规则字符串
func (f *FilterTag) ToRuleTag() string {
	var buf bytes.Buffer
	buf.WriteString(string(f.fieldName))
	buf.WriteString("@")
	buf.WriteString(f.rules)
	buf.WriteString("#")
	buf.WriteString(f.validMsg)
	return buf.String()
}

func (f *FilterTag) GetFieldName() FieldName {
	return f.fieldName
}

func (f *FilterTag) GetDataType() FieldDataType {
	return f.dataType
}

func (f *FilterTag) GetFormatter() string {
	return f.formatter
}

func (f *FilterTag) GetDefValue() interface{} {
	return f.defValue
}

func (f *FilterTag) GetRules() string {
	return f.rules
}

func (f *FilterTag) GetValidMsg() string {
	return f.validMsg
}

func (f *FilterTag) SetFieldName(fieldName FieldName) {
	f.fieldName = fieldName
}

func (f *FilterTag) SetDataType(dataType FieldDataType) {
	f.dataType = dataType
}

func (f *FilterTag) SetFormatter(formatter string) {
	f.formatter = formatter
}

func (f *FilterTag) SetDefValue(v interface{}) {
	f.defValue = v
}

func (f *FilterTag) SetRules(rules string) {
	f.rules = rules
}

func (f *FilterTag) SetValidMsg(msg string) {
	f.validMsg = msg
}

func (f *FilterTag) SetParamName(paramName string) {
	f.paramName = paramName
}

func (f *FilterTag) GetParamName() string {
	return f.paramName
}

func (f *FilterTag) ToMap() map[string]interface{} {
	results := make(map[string]interface{})
	results["fieldName"] 	= f.fieldName
	results["paramName"]	= f.paramName
	results["dataType"] 	= f.dataType
	results["formatter"] 	= f.formatter
	results["defValue"] 	= f.defValue
	results["rules"] 		= f.rules
	results["validMsg"] 	= f.validMsg
	return results
}

func (f *FilterTag) ToJson() string {
	return xjson.ToString(f.ToMap(), "", "\t")
}

