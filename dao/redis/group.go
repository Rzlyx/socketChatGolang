package redis

import (
	"encoding/json"
	"fmt"
)

type UserGroup struct {
	UserId    int64 `json:"user_id"`
	GroupID   int64 `json:"group_id"`
	Type      int   `json:"type"`
	IsSlience bool  `json:"is_slience"`
}

func AddUserGroup(msg string, key string) error {
	res := rdb.Do("rpush", key, msg)
	return res.Err()
}

func GetUserGroup(key string) ([]string, error) {
	res := rdb.LRange(key, 0, -1)
	return res.Val(), res.Err()
}

func DeleteUserGroup(key string) error {
	res := rdb.Del(key)
	return res.Err()
}

func UpdateUserGroup(key string, newValue string) error {
	res := rdb.Do("SET", key, newValue)
	return res.Err()
}

func TurnStringFromNode(node *UserGroup) (string, error) {
	data, err := json.Marshal(*node)
	if err != nil {
		fmt.Println("[TurnStringFromNode], Marshal err is ", err.Error())
		return "", err
	}
	value := string(data)
	return value, nil
}

func TurnNodeFromNode(value string) (*UserGroup, error) {
	var node UserGroup
	err := json.Unmarshal([]byte(value), &node)
	if err != nil {
		fmt.Println("[TurnNodeFromNode], Unmarshal err is ", err.Error())
		return &node, err
	}
	return &node, nil
}
