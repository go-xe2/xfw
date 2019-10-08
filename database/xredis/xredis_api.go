package xredis

import (
	"encoding/json"
	"github.com/gogf/gf/g/util/gconv"
	"time"
)

// 设置值
func Set(key string, v interface{}, expiration ...time.Duration) (bool, error) {
	redis := Client().DB()
	n  := 0 * time.Second
	if len(expiration) > 0 {
		n = expiration[0]
	}
	data := gconv.String(v)
	err := redis.Set(key, data, n).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 获取键值
func Get(key string, out ...interface{}) (string, error) {
	redis := Client().DB()
	s, err := redis.Get(key).Result()
	if err != nil {
		return "", nil
	}
	if len(out) > 0 {
		err := json.Unmarshal([]byte(s), out[0])
		if err != nil {
			return "", err
		}
		return s, nil
	}
	return s, nil
}

// 设置超时时间
func SetExpire(key string, expiration time.Duration) bool  {
	redis := Client().DB()
	err := redis.Expire(key, expiration).Err()
	return  err == nil
}

func Remove(key string) bool  {
	redis := Client().DB()
	err := redis.Del(key).Err()
	return  err == nil
}



