package xfw

import "github.com/gogf/gf/g/container/gmap"

var (
	instanceMap = gmap.NewStrAnyMap()
)

// 保存实例
func SetInstance(name string, i interface{}) {
	instanceMap.Set(name, i)
}

// 获取实例
func GetInstance(name string) interface{} {
	return instanceMap.Get(name)
}
