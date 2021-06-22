package util

import (
	"time"
)

func GetServerCurrentTime(v int64) int64 {
	return time.Now().UnixNano()/1e6 + v
}

func GetServerCurrentTimeOffset(v int64) int64 {
	return v - time.Now().UnixNano()/1e6
}
