package client

import (
	"encoding/json"
	"reflect"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

func (c *Client) dumpClientToServerMessage(
	c2s *codec.ClientToServerMessage,
	msg interface{},
) {
	typ := reflect.TypeOf(msg)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	dump, _ := json.Marshal(&msg)
	log.Debug().
		Uint32("@seq", c2s.Seq).
		Str("method", c2s.ServiceMethod).
		Str("uin", c2s.Username).
		Msg("<<< [dump] message:" + typ.String() + ":" + string(dump))
}

func (c *Client) dumpServerToClientMessage(
	s2c *codec.ServerToClientMessage,
	msg interface{},
) {
	typ := reflect.TypeOf(msg)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	dump, _ := json.Marshal(&msg)
	log.Debug().
		Uint32("@seq", s2c.Seq).
		Str("method", s2c.ServiceMethod).
		Str("uin", s2c.Username).
		Msg(">>> [dump] message:" + typ.String() + ":" + string(dump))
}
