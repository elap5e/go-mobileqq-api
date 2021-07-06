package client

import (
	"encoding/json"
	"reflect"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
)

func (c *Client) dumpClientToServerMessage(c2s *codec.ClientToServerMessage, msg interface{}) {
	typ := reflect.TypeOf(msg)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	dump, _ := json.Marshal(&msg)
	log.Debug().
		Str("method", c2s.ServiceMethod).
		Uint32("seq", c2s.Seq).
		Str("uin", c2s.Username).
		Msg("<<< [dump] message:" + typ.String() + ":" + string(dump))
}

func (c *Client) dumpServerToClientMessage(s2c *codec.ServerToClientMessage, msg interface{}) {
	typ := reflect.TypeOf(msg)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	dump, _ := json.Marshal(&msg)
	log.Debug().
		Str("method", s2c.ServiceMethod).
		Uint32("seq", s2c.Seq).
		Str("uin", s2c.Username).
		Msg(">>> [dump] message:" + typ.String() + ":" + string(dump))
}

func (c *Client) marshalMessage(msg *pb.Message) ([]byte, error) {
	data, err := mark.Marshal(msg)
	log.Info().
		Str("@mark", string(data)).
		Uint64("@peer", msg.GetMessageHead().GetGroupInfo().GetGroupCode()).
		Uint32("@seq", msg.GetMessageHead().GetMessageSeq()).
		Uint32("@time", msg.GetMessageHead().GetMessageTime()).
		Uint64("from", msg.GetMessageHead().GetFromUin()).
		Uint64("to", msg.GetMessageHead().GetToUin()).
		Uint32("type", msg.GetMessageHead().GetMessageType()).
		Uint64("uid", msg.GetMessageHead().GetMessageUid()).
		Msg("--> [recv] message")
	return data, err
}
