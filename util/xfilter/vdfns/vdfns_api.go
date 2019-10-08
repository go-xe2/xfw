package vdfns

import (
	"github.com/go-xe2/xfw/util/xfilter"
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/util/gconv"
	"time"
)


// 获取时间戳
func Timestamp(params ...interface{}) interface{} {
	return time.Now().Unix()
}

func DateString(params ...interface{}) interface{} {
	return gtime.Now().Format("Y-m-d H:i:s")
}

// 字符串传时间戳
func Str2timestamp(params ...interface{}) interface{} {
	if len(params) > 0 {
		if t := gtime.ParseTimeFromContent(gconv.String(params[0])); t != nil {
			return t.Unix()
		}
	}
	return time.Now().Unix()
}

// 时间戳转字符串
func Timestamp2str(params ...interface{}) interface{} {
	if len(params) > 0 {
		if t := gtime.NewFromTimeStamp(gconv.Int64(params[0])); t != nil {
			return t.Format("Y-m-d H:i:s")
		}
	}
	return ""
}

func init() {
	fnMgr := xfilter.FuncManager()
	fnMgr.Register("timestamp", Timestamp)
	fnMgr.Register("datestring", DateString)
	fnMgr.RegisterFormat("str2timestamp", Str2timestamp)
	fnMgr.RegisterFormat("timestamp2str", Timestamp2str)
}
