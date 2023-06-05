package utils

import "os"

func GetCurrentPath() string {
	path, _ := os.Getwd()
	return path
}
