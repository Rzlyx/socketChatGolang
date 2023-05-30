package redis

func AddMsg(msg string, id string) error {
	res := rdb.Do("rpush", id, msg)
	return res.Err()
}

func GetMsg(id string) (error, []string) {
	res := rdb.LRange(id, 0, -1)
	return res.Err(), res.Val()
}
