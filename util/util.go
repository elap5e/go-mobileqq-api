package util

import (
	"bufio"
	"strconv"
	"strings"
	"time"
)

func GetServerCurrentTime(v int64) int64 {
	return time.Now().UnixNano()/1e6 + v
}

func GetServerCurrentTimeOffset(v int64) int64 {
	return v - time.Now().UnixNano()/1e6
}

func CheckUsername(username string) bool {
	uin, err := strconv.Atoi(username)
	if err != nil || uin < 10000 || uin > 4000000000 {
		return false
	}
	return true
}

func ReadLine(rd *bufio.Reader) (string, error) {
	str, err := rd.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(str, "\n"), nil
}
