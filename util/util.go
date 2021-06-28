package util

import (
	"bufio"
	"strconv"
	"strings"
	"time"
)

var (
	serverTimeInterval = uint32(0x00000000) // TODO: add atomic
)

func SetServerTime(v uint32) {
	serverTimeInterval = v - uint32(time.Now().Unix())
}

func GetServerTime() uint32 {
	return uint32(time.Now().Unix()) + serverTimeInterval
}

func CheckUin(uin uint64) bool {
	if uin < 10000 || uin > 4000000000 {
		return false
	}
	return true
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
