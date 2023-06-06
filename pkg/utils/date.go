package utils

import "time"

func GetNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetTime(year int, month int, day int) string {
	return time.Now().AddDate(year, month, day).Format("2006-01-02 15:04:05")
}
