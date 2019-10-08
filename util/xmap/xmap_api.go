// map api操作库
package xmap

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/util/gconv"
	"reflect"
	"sort"
	"strings"
	"time"
)

func GetStrDef (m map[string]interface{}, key string, def string) string {
	if s, ok := m[key].(string); !ok {
		return def
	} else {
		return s
	}
}

func GetStr (m map[string]interface{}, key string) string {
	return GetStrDef(m, key, "")
}

func GetIntDef (m map[string]interface{}, key string, def int) int {
	v, ok := m[key]
	if !ok {
		return def
	}
	return gconv.Int(v)
}

func GetInt (m map[string]interface{}, key string) int  {
	return GetIntDef(m, key, 0)
}

func GetInt32Def (m map[string]interface{}, key string, def int32) int32 {
	v, ok := m[key]
	if !ok {
		return def
	}
	return gconv.Int32(v)
}

func GetInt8Def (m map[string]interface{}, key string, def int8) int8 {
	v, ok := m[key]
	if !ok {
		return def
	}
	return gconv.Int8(v)
}

func GetInt8 (m map[string]interface{}, key string) int8 {
	return GetInt8Def(m, key, 0)
}

func GetInt16 (m map[string]interface{}, key string, def int16) int16  {
	v, ok := m[key]
	if !ok {
		return def
	}
	return gconv.Int16(v)
}

func GetInt32 (m map[string]interface{}, key string) int32 {
	return GetInt32Def(m, key, 0)
}

func GetInt64Def (m map[string]interface{}, key string, def int64) int64 {
	v, ok := m[key]
	if !ok {
		return def
	}
	return gconv.Int64(v)
}

func GetInt64 (m map[string]interface{}, key string) int64 {
	return GetInt64Def(m, key, 0)
}

func GetUintDef (m map[string]interface{}, key string, def uint) uint {
	v, ok := m[key]
	if !ok {
		return def
	}
	return gconv.Uint(v)
}

func GetUint8Def (m map[string]interface{}, key string, def uint8) uint8 {
	v, ok := m[key]
	if !ok {
		return def
	}
	return gconv.Uint8(v)
}

func GetUint8 (m map[string]interface{}, key string) uint8  {
	return GetUint8Def(m, key, 0)
}

func GetUint16Def (m map[string]interface{}, key string, def uint16) uint16 {
	v, ok := m[key]
	if !ok {
		return def
	}
	return gconv.Uint16(v)
}

func GetUint16 (m map[string]interface{}, key string) uint16  {
	return GetUint16Def(m, key, 0)
}

func GetUint32Def (m map[string]interface{}, key string, def uint32) uint32  {
	v, ok := m[key]
	if !ok {
		return def
	}
	return gconv.Uint32(v)
}

func GetUint32 (m map[string]interface{}, key string) uint32 {
	return GetUint32Def(m, key, 0)
}

func GetUint64Def (m map[string]interface{}, key string, def uint64) uint64  {
	v, ok := m[key]
	if !ok {
		return def
	}
	return gconv.Uint64(v)
}

func GetUint64 (m map[string]interface{}, key string) uint64 {
	return GetUint64Def(m, key, 0)
}

func GetFloat32Def (m map[string]interface{}, key string, def float32) float32 {
	v, ok := m[key]
	if !ok {
		return def
	}
	return gconv.Float32(v)
}

func GetFloat32 (m map[string]interface{}, key string) float32 {
	return GetFloat32Def(m, key, 0)
}

func GetFloat64Def (m map[string]interface{}, key string, def float64) float64 {
	v, ok := m[key]
	if !ok {
		return def
	}
	return gconv.Float64(v)
}

func GetFloat64 (m map[string]interface{}, key string) float64 {
	return GetFloat64Def(m, key, 0)
}

func GetTime (m map[string]interface{}, key string, def time.Time) time.Time  {
	v, ok := m[key]
	if !ok {
		return def
	}
	return gconv.Time(v, "Y-m-d H:i:s")
}

func GetObject (m map[string]interface{}, key string) interface{} {
	if v, ok := m[key].(interface{}); ok {
		return v
	}
	return nil
}

func GetMap (m map[string]interface{}, key string) map[string]interface{} {
	if v, ok := m[key].(map[string]interface{}); ok {
		return v
	}
	return make(map[string]interface{})
}

func GetArray (m map[string]interface{}, key string) []interface{}  {
	if v, ok := m[key].([]interface{}); ok {
		return v
	}
	return []interface{}{}
}


// map转换成可作为键名的字符串
func ToKeyString(m map[string]interface{}) string {
	if m == nil {
		return ""
	}
	var result bytes.Buffer
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := m[key]
		result.WriteString(key)
		result.WriteString("-")
		result.WriteString(fmt.Sprintf("%v", value))
		result.WriteString("_")
	}
	return result.String()
}

// 拷贝map
// fields 可传字段列表字符串，字符串数组，或map,如果为map可设置默认值
func Clone(m map[string]interface{}, fields ...interface{}) map[string]interface{} {
	var selectFields = make(map[string]interface{})
	if len(fields) > 0 {
		switch items := fields[0].(type) {
		case []string:
			for _, key := range items {
				if v, ok := m[key]; ok {
					selectFields[key] = v
				}
			}
			return selectFields
		case map[string]interface{}:
			for k, v := range items {
				if v1, ok := m[k]; ok {
					selectFields[k] = v1
				} else {
					selectFields[k] = v
				}
			}
			return selectFields
		case string:
			strItems := strings.Split(items, ",")
			for _, key := range strItems {
				key = strings.Trim(key, " ")
				if v , ok := m[key]; ok {
					selectFields[key] = v
				}
			}
			return selectFields
		default:
			glog.Info("xmap.clone unknown fields params type:", reflect.TypeOf(fields[0]))
		}
		return nil
	}
	for k, v := range m {
		selectFields[k] = v
	}
	return selectFields
}