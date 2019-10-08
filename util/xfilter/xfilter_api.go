// 表单过滤库
package xfilter

import (
	"encoding/json"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/text/gstr"
	"github.com/gogf/gf/g/util/gconv"
	"github.com/gogf/gf/g/util/gvalid"
	"reflect"
	"strings"
)

const XFilterTagName = "filter"
const XOrmTagName = "orm"

// 检查是否包含字段
func hasField(fieldName string, selects ...map[string]interface{}) bool {
	if len(selects) == 0 {
		return true
	}
	fields := selects[0]
	if n, ok := fields[fieldName]; ok {
		return gconv.Bool(n)
	} else {
		return false
	}
}

// 从map中获取过滤字段列表
func GetFilterFromMap(mp map[string]string, selects ...map[string]interface{}) []*FilterTag {
	var result []*FilterTag
	for fdName, str := range mp {
		if !hasField(fdName, selects...) {
			continue
		}
		item := NewTag(str)
		item.SetFieldName(FieldName(fdName))
		result = append(result, item)
	}
	return result
}

// 从字符串数组中获取过滤规则列表
func GetFilterFromArray(arr []string, selects ...map[string]interface{}) []*FilterTag {
	var result []*FilterTag
	for _, str := range arr {
		item := NewTag(str)
		if !hasField(string(item.fieldName), selects...) {
			continue
		}
		result = append(result, item)
	}
	return result
}

// 从json字符串获取过滤规则
func GetFilterFromJson(js string, selects ...map[string]interface{}) []*FilterTag {
	var v map[string]interface{}
	var result []*FilterTag
	err := json.Unmarshal([]byte(js), &v)
	if err != nil {
		glog.Error("GetFilterFromJson parse json error:", err)
		return []*FilterTag{}
	}
	for fdName, item := range v {
		if f, ok := item.(map[string]interface{}); !ok {
			continue
		} else {
			if !hasField(fdName, selects...) {
				continue
			}
			ftItem := NewTag()
			ftItem.FromJson(f)
			ftItem.SetFieldName(FieldName(fdName))
			result = append(result, ftItem)
		}
	}
	return result
}

// 从struct中获取过滤规则
func GetFilterFromStruct(obj interface{}, selects ...map[string]interface{}) []*FilterTag {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		glog.Error("GetFilterFromStruct fail, obj is not struct ptr, obj kind is:", reflect.TypeOf(obj).Kind())
		return []*FilterTag{}
	}
	var result []*FilterTag
	elem := reflect.TypeOf(obj).Elem()
	for i := 0; i < elem.NumField(); i++ {
		fieldName := elem.Field(i).Name
		if !hasField(fieldName, selects...) {
			continue
		}
		tag := elem.Field(i).Tag.Get(XFilterTagName)
		if tag != "" {
			item := NewTag(tag)
			if item.fieldName == "" {
				ormTag := elem.Field(i).Tag.Get(XOrmTagName)
				if ormTag != "" {
					item.fieldName = FieldName(ormTag)
				} else {
					item.fieldName = FieldName(fieldName)
				}
			}
			if item.dataType == "" {
				switch elem.Field(i).Type.Kind() {
				case reflect.String:
					item.dataType = "string"
					break
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					item.dataType = "int"
					break
				case reflect.Float32, reflect.Float64:
					item.dataType = "float"
					break
				case reflect.Bool:
					item.dataType = "bool"
					break
				}
			}
			result = append(result, item)
		}
	}
	return result
}

// 解析默认值
func parseDefaultValue(value interface{}) interface{}  {
	if expr, ok := value.(string); ok && expr != "" {
		if strings.HasPrefix(expr, ":") {
			fnName := gstr.SubStr(expr, 1)
			if FuncManager().HasFunc(fnName) {
				v, err := FuncManager().Call(fnName)
				if err == nil {
					return v
				}
			}
		}
	}
	return value
}

// 转换成规则所需要的数据类型
func convertToRuleDataType(dateType FieldDataType, value interface{}) interface{} {
	switch dateType {
	case "string":
		return gconv.String(value)
	case "int":
		return gconv.Int64(value)
	case "float":
		return gconv.Float64(value)
	case "bool":
		return gconv.Bool(value)
	case "time":
		return gconv.Time(value)
	default:
		return gconv.String(value)
	}
}

// 根据规则列表过滤map
// onlyParams 是否只返回src中存在的字段
func doFilter(src map[string]interface{}, filters []*FilterTag, onlyParams ...bool) (map[string]interface{}, *gvalid.Error)  {
	var result = make(map[string]interface{})
	var validRuleMap []string
	var only = false
	if len(onlyParams) > 0 {
		only = onlyParams[0]
	}
	for _, filter := range filters {
		fieldName := string(filter.fieldName)
		paramName := filter.paramName
		if paramName == "" {
			paramName = fieldName
		}
		if filter.rules != "" {
			validRuleMap = append(validRuleMap, filter.ToRuleTag())
		}
		// 忽略参数
		if paramName == "-" {
			v := parseDefaultValue(filter.defValue)
			if v != nil {
				result[fieldName] = v
			}
			continue
		}
		fnMgr := FuncManager()
		if v, ok := src[paramName]; ok && v != nil {
			if fnMgr.HasFormatter(filter.formatter) {
				if v, err := fnMgr.Format(filter.formatter, convertToRuleDataType(filter.dataType, v)); err == nil {
					result[fieldName] = v
				}
			} else {
				result[fieldName] = convertToRuleDataType(filter.dataType, v)
			}
		} else {
			if !only && filter.defValue != nil {
				v := parseDefaultValue(filter.defValue)
				if v == nil {
					continue
				}
				if fnMgr.HasFormatter(filter.formatter) {
					if v1, err := fnMgr.Format(filter.formatter, v); err != nil {
						result[fieldName] = v1
					}
				} else {
					result[fieldName] = convertToRuleDataType(filter.dataType, v)
				}
			}
		}
	}

	// 验证数据
	err := gvalid.CheckMap(result, validRuleMap)
	if err != nil {
		return result, err
	}
	return result, nil
}

// 使用map规则过滤
func FilterMap(src map[string]interface{}, rules map[string]string, selects ...map[string]interface{}) (map[string]interface{}, *gvalid.Error) {
	filters := GetFilterFromMap(rules, selects...)
	return doFilter(src, filters)
}

// 使用规则数组过滤
func Filter(src map[string]interface{}, rules []string, selects ...map[string]interface{}) (map[string]interface{}, *gvalid.Error)  {
	filters := GetFilterFromArray(rules, selects...)
	return doFilter(src, filters)
}

// 使用规则数组过滤,并且只返回输入参数中存在的字段
func FilterOnly(src map[string]interface{}, rules []string, selects ...map[string]interface{}) (map[string]interface{}, *gvalid.Error)  {
	filters := GetFilterFromArray(rules, selects...)
	return doFilter(src, filters, true)
}


// 使用struct定义过滤
func FilterStruct(src map[string]interface{}, rules interface{}, selects ...map[string]interface{}) (map[string]interface{}, *gvalid.Error)  {
	filters := GetFilterFromStruct(rules, selects...)
	return doFilter(src, filters)
}

// 使用struct定义过滤并且只返回src中存在的字段
func FilterStructOnly(src map[string]interface{}, rules interface{}, selects ...map[string]interface{}) (map[string]interface{}, *gvalid.Error)  {
	filters := GetFilterFromStruct(rules, selects...)
	return doFilter(src, filters, true)
}

// 使用json验证规则过滤
func FilterJson(src map[string]interface{}, jsonRule string, selects ...map[string]interface{}) (map[string]interface{}, *gvalid.Error)  {
	filters := GetFilterFromJson(jsonRule, selects...)
	return doFilter(src, filters)
}

// 使用json验证规则过滤,前且只返回src中存在的字段
func FilterJsonOnly(src map[string]interface{}, jsonRule string, selects ...map[string]interface{}) (map[string]interface{}, *gvalid.Error)  {
	filters := GetFilterFromJson(jsonRule, selects...)
	return doFilter(src, filters, true)
}


