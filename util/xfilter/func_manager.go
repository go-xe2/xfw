package xfilter

import (
	xfw "github.com/go-xe2/xfw/os"
	"github.com/go-xe2/xfw/xerrors"
)

type DefValueFunc func(params ...interface{}) interface{}

type IFuncManager interface {
	Register(name string, fn DefValueFunc) bool
	Unregister(fnName string) bool
	RegisterFormat(name string, fn DefValueFunc) bool
	UnregisterFormat(name string) bool
	HasFunc(fnName string) bool
	HasFormatter(fnName string) bool
	Call(fnName string, params ...interface{}) (interface{}, error)
	Format(fnName string, params ...interface{}) (interface{}, error)
}

type funcManager struct {
	fnItems 	map[string]DefValueFunc
	formats		map[string]DefValueFunc
}

const x_FUNC_MANAGER_NAME = "funcManager"

func FuncManager(name ...string) IFuncManager {
	fnName := x_FUNC_MANAGER_NAME
	if len(name) > 0 {
		fnName = name[0]
	}
	if c := xfw.GetInstance(fnName); c != nil {
		return c.(IFuncManager)
	}

	c := &funcManager{
		fnItems: make(map[string]DefValueFunc),
		formats: make(map[string]DefValueFunc),
	}
	xfw.SetInstance(fnName, c)
	return c
}


func (f *funcManager) Register(name string, fn DefValueFunc) bool {
	if _, ok := f.fnItems[name]; !ok {
		f.fnItems[name] = fn
		return true
	}
	return false
}

func (f *funcManager) Unregister(fn string) bool {
	if _, ok := f.fnItems[fn]; ok {
		delete(f.fnItems, fn)
		return true
	}
	return false
}

// 注册格式化方法
func (f *funcManager) RegisterFormat(name string, fn DefValueFunc) bool {
	if _, ok := f.formats[name]; !ok {
		f.formats[name] = fn
		return true
	}
	return false
}

// 注册格式化方法
func (f *funcManager) UnregisterFormat(name string) bool {
	if _, ok := f.fnItems[name]; ok {
		delete(f.formats, name)
		return true
	}
	return false
}

func (f *funcManager) Call(fnName string, params ...interface{}) (interface{}, error) {
	if fn, ok := f.fnItems[fnName]; ok {
		return fn(params...), nil
	}
	return nil, xerrors.New(fnName, "不存在")
}
func (f *funcManager) Format(fnName string, params ...interface{}) (interface{}, error) {
	if fn, ok := f.formats[fnName]; ok {
		return fn(params...), nil
	}
	return nil, xerrors.New(fnName, "不存在")
}

func (f *funcManager) HasFunc(fnName string) bool {
	if _, ok := f.fnItems[fnName]; ok {
		return true
	}
	return false
}

func (f *funcManager) HasFormatter(fnName string) bool {
	if _, ok := f.formats[fnName]; ok {
		return true
	}
	return false
}

