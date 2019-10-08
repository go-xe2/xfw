package xmgo

import (
	"gmny-server/xfw/xerrors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)


// 获取集合
func C(c interface{}) (*mgo.Collection, error) {
	mgodb := Xmgo().DB()
	switch value := c.(type) {
	case string:
		return mgodb.C(value), nil
	case IxmgoCollection:
		return mgodb.C(value.Name()), nil
	default:
		return nil, xerrors.New("集合不支持string和实现IxmgoCollection的struct")
	}
}

// 获取一条数据
func FindOne(c interface{}, where map[string]interface{}, selects ...interface{}) (interface{}, error) {
	cl, err := C(c)
	if err != nil {
		return nil, err
	}
	query := cl.Find(where)
	if len(selects) > 0 {
		query = query.Select(selects[0])
	}
	if _, ok := c.(IxmgoCollection); ok {
		if err := query.One(c); err != nil {
			return nil, err
		}
		return c, nil
	}

	// c is string
	var doc interface{}
	if err := query.One(&doc); err != nil {
		return nil, err
	}
	return doc, nil
}

// 根据id查找记录
func FindById(c interface{}, id interface{}, selects ...map[string]interface{}) (interface{}, error) {
	cl , err := C(c)
	if err != nil {
		return nil, err
	}
	switch value := id.(type) {
	case string:
		if !bson.IsObjectIdHex(value) {
			return nil, xerrors.New("FindById失败,不是有效的ObjectId:", id)
		}
		id = bson.ObjectIdHex(value)
		break
	case bson.ObjectId:
		break
	default:
		return nil, xerrors.New("FindById失败，不是有效的ObjectId:", id)
	}
	query := cl.FindId(id)
	if len(selects) > 0 {
		query = query.Select(selects[0])
	}
	if _, ok := c.(IxmgoCollection); ok {
		if err := query.One(c); err != nil {
			return nil, err
		}
		return c, nil
	}

	// c is string
	var doc interface{}
	if err = query.One(&doc); err != nil {
		return nil, err
	}
	return doc, nil
}

// 查询数据
func Query(c interface{}, where map[string]interface{}) (*mgo.Query, error) {
	cl, err := C(c)
	if err != nil {
		return nil, err
	}
	return cl.Find(where), nil
}

// 根据ID更新记录
func UpdateById(c interface{}, id interface{}, data interface{}) error {
	cl, err := C(c)
	if err != nil {
		return err
	}
	switch value := id.(type) {
	case string:
		if !bson.IsObjectIdHex(value) {
			return xerrors.New("UpdateById失败,不是有效的ObjectId:", id)
		}
		id = bson.ObjectIdHex(value)
		break
	case bson.ObjectId:
		id = value
		break
	default:
		return xerrors.New("UpdateById失败， 不是有效的ObjectId:", id)
	}
	return cl.UpdateId(id, bson.M{
		"$set": data,
	})
}

// 更新一条记录
func Update(c interface{}, where map[string]interface{}, data interface{}) error  {
	cl, err := C(c)
	if err != nil {
		return err
	}
	return cl.Update(where, bson.M{
		"$set": data,
	})
}

// 更新所有记录
func UpdateAll(c interface{}, where map[string]interface{}, data interface{}) (*mgo.ChangeInfo, error)  {
	cl, err := C(c)
	if err != nil {
		return nil, err
	}
	return cl.UpdateAll(where, bson.M{
		"$set": data,
	})
}

// 添加记录
func Insert(c interface{}, docs map[string]interface{}) error {
	cl, err := C(c)
	if err != nil {
		return err
	}
	return cl.Insert(docs)
}

// 删除集合
func Drop(c interface{}) error {
	cl, err := C(c)
	if err != nil {
		return err
	}
	return cl.DropCollection()
}




