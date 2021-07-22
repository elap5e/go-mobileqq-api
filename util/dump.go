package util

import (
	"encoding/binary"
	"encoding/json"
	"net"
	"reflect"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/rs/zerolog"
)

func DumpClientToServerMessage(
	c2s *codec.ClientToServerMessage,
	msg interface{},
) {
	if log.GetLevel() > zerolog.DebugLevel {
		return
	}
	value := reflect.ValueOf(msg)
	for value.Kind() == reflect.Interface ||
		(value.Kind() == reflect.Ptr && !value.IsNil()) {
		value = value.Elem()
	}
	typ := value.Type()
	dump, _ := json.Marshal(&msg)
	log.Debug().
		Uint32("@seq", c2s.Seq).
		Str("method", c2s.ServiceMethod).
		Str("uin", c2s.Username).
		Msg("<<< [dump] message:" + typ.String() + ":" + string(dump))
}

func DumpServerToClientMessage(
	s2c *codec.ServerToClientMessage,
	msg interface{},
) {
	if log.GetLevel() > zerolog.DebugLevel {
		return
	}
	value := reflect.ValueOf(msg)
	for value.Kind() == reflect.Interface ||
		(value.Kind() == reflect.Ptr && !value.IsNil()) {
		value = value.Elem()
	}
	typ := value.Type()
	dump, _ := json.Marshal(&msg)
	log.Debug().
		Uint32("@seq", s2c.Seq).
		Str("method", s2c.ServiceMethod).
		Str("uin", s2c.Username).
		Msg(">>> [dump] message:" + typ.String() + ":" + string(dump))
}

func DumpServerIP(u uint32) {
	ip := net.IP{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(ip, u)
	log.Debug().
		Msg(">>> [dump] server ip:" + ip.String())
}
