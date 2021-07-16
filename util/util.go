package util

import (
	"bufio"
	"fmt"
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

func STBytesTobytes(v [16]byte) (ret []byte) {
	ret = make([]byte, 16)
	copy(ret, v[:])
	return
}

func BytesToSTBytes(v []byte) (ret [16]byte) {
	copy(ret[:], v)
	return
}

func HashToBraceString(p []byte) string {
	u := make([]byte, 16)
	copy(u, p)
	return fmt.Sprintf("{%X-%X-%X-%X-%X}", u[0:4], u[4:6], u[6:8], u[8:10], u[10:16])
}

func ParseExtToPictureType(ext string) uint32 {
	switch ext {
	case ".jpeg", ".jpg":
		return 1000
	case ".png":
		return 1001
	case ".webp":
		return 1002
	case ".sharpp":
		return 1004
	case ".bmp":
		return 1005
	case ".gif":
		return 2000
	case ".apng":
		return 2001
	}
	return 0
}
