package util

import (
	"strconv"
	"time"
)

func GetServerCurrentTime(v int64) int64 {
	return time.Now().UnixNano()/1e6 + v
}

func GetServerCurrentTimeOffset(v int64) int64 {
	return v - time.Now().UnixNano()/1e6
}

func CheckUinAccount(s string) bool {
	uin, err := strconv.Atoi(s)
	if err != nil || uin < 10000 || uin > 4000000000 {
		return false
	}
	return true
}
