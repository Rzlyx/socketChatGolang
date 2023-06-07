package redis

func AddMsg(msg string, id string) error {
	res := rdb.Do("rpush", id, msg)
	return res.Err()
}

func GetMsg(id string) ([]string, error) {
	res := rdb.LRange(id, 0, -1)
	return res.Val(), res.Err()
}

func DeleteMsg(id string) error {
	res := rdb.Del(id)
	return res.Err()
}