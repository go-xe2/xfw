package xmysql

import (
	"github.com/go-xe2/xfw/xerrors"
	"github.com/go-xe2/xorm"
)

// 检查数据库中是否存在记录
func Exists (table interface{}, where ...interface{}) (bool, error) {
	db := Mysql().DB()
	result, err := db.Table(table).Where(where...).Count()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// 普通查询
func Select(table string, fields []string, where []interface{}, order string) ([]xorm.Data, error) {
	db := Mysql().DB()
	r, err := db.Table(table).Fields(fields...).Where(where...).Order(order).Get()
	if err != nil {
		return nil, xerrors.New("Select error:", db.LastSql())
	}
	return r, nil
}

// 分页查询
func SelectPage(table string, fields []string, where []interface{}, order string, pageIndex, pageSize int) (rows []xorm.Data, count int64, pageCount int, err error) {
	db := Mysql().DB()
	var mCount int64
	if pageSize == 0 {
		pageSize = 10
	}
	if pageIndex < 0 {
		pageIndex = 0
	}
	mCount, err = db.Table(table).Where(where...).Count()
	if err != nil {
		return nil, 0, 0, xerrors.New("Select error:", db.LastSql())
	}
	db.Reset()
	r, err := db.Table(table).Fields(fields...).Where(where...).Order(order).Offset(pageIndex * pageSize).Limit(pageSize).Get()
	if err != nil {
		return nil, 0, 0, xerrors.New("Select error:", db.LastSql())
	}
	dev := mCount % int64(pageSize)
	pageCount = int(mCount / int64(pageSize))
	if dev > 0 {
		pageCount += 1
	}
	return r, mCount, pageCount, nil
}


// 获取指定条数的记录
func SelectTop (table string, fields string, where []interface{}, order string, count int) ([]xorm.Data, error)  {
	db := Mysql().DB()
	r, err := db.Table(table).Fields(fields).Where(where...).Order(order).Limit(count).Get()
	if err != nil {
		return nil, xerrors.New("Select error:", db.LastSql())
	}
	return r, nil
}

