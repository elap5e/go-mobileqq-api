package rpc

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
)

func (c *Client) dumpClientToServerMessage(c2s *codec.ClientToServerMessage, msg interface{}) {
	if c.cfg.LogLevel&LogLevelTrace != 0 {
		typ := reflect.TypeOf(msg)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		dump, _ := json.MarshalIndent(&msg, "", "  ")
		log.Printf(
			"<<< [dump] seq:0x%08x uin:%s method:%s message:%s:\n%s",
			c2s.Seq, c2s.Username, c2s.ServiceMethod, typ.String(), dump,
		)
	}
}

func (c *Client) dumpServerToClientMessage(s2c *codec.ServerToClientMessage, msg interface{}) {
	if c.cfg.LogLevel&LogLevelTrace != 0 {
		typ := reflect.TypeOf(msg)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		dump, _ := json.MarshalIndent(&msg, "", "  ")
		log.Printf(
			">>> [dump] seq:0x%08x uin:%s method:%s message:%s:\n%s",
			s2c.Seq, s2c.Username, s2c.ServiceMethod, typ.String(), dump,
		)
	}
}

func (c *Client) marshalMessage(msg *pb.Message) ([]byte, error) {
	data, err := mark.Marshal(msg)
	log.Printf(
		">>> [dump] time:%d type:%d peer:%d seq:%d uid:%d from:%d to:%d markdown:\n%s",
		msg.GetMessageHead().GetMessageTime(),
		msg.GetMessageHead().GetMessageType(),
		msg.GetMessageHead().GetGroupInfo().GetGroupCode(),
		msg.GetMessageHead().GetMessageSeq(),
		msg.GetMessageHead().GetMessageUid(),
		msg.GetMessageHead().GetFromUin(),
		msg.GetMessageHead().GetToUin(),
		string(data),
	)
	return data, err
}
