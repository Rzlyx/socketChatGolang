package redis

import "fmt"

type UserCache struct {
	//userID             int64 `json:"user_id"`
	PrivateChat_black *[]int64 `json:"privateChat_Black"`
	PrivateChat_white *[]int64 `json:"privateChat_White"`
	PrivateChat_gray  *[]int64 `json:"privateChat_Gray"`
}

func Add(key string, val string) error {
	res := rdb.Do("set", key, val)
	return res.Err()
}

func Get(key string) (string, error) {
	res := rdb.Do("get", key)
	val := fmt.Sprintf("%v", res.Val())
	return val, res.Err()
}

func Delete(key string, val string) error {
	res := rdb.Do("del", key)
	return res.Err()
}
