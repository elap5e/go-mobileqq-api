package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/util"
)

type OnlinePushRequest struct {
	Uin             int64             `jce:",0" json:"uin,omitempty"`
	Time            int64             `jce:",1" json:"time,omitempty"`
	Messages        []*Message        `jce:",2" json:"messages,omitempty"`
	ServerIP        Uint32IPType      `jce:",3" json:"server_ip,omitempty"`
	SyncCookie      []byte            `jce:",4" json:"sync_cookie,omitempty"`
	UinPairMessages []*UinPairMessage `jce:",5" json:"uin_pair_messages,omitempty"`
	Previews        map[string][]byte `jce:",6" json:"previews,omitempty"`
	UserActive      int32             `jce:",7" json:"user_active,omitempty"`
	GeneralFlag     int32             `jce:",12" json:"general_flag,omitempty"`
}

type OnlinePushResponse struct {
	Type   int32  `jce:",1" json:"type,omitempty"`
	Seq    int64  `jce:",2" json:"seq,omitempty"`
	Buffer []byte `jce:",3" json:"buffer,omitempty"`
}

func NewOnlinePushMessageResponse(
	ctx context.Context,
	username string,
	items []MessageDelete,
	serverIP Uint32IPType,
	seq int32,
) (*codec.ClientToServerMessage, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("zero length")
	}

	uin, err := strconv.ParseInt(username, 10, 64)
	if err != nil {
		return nil, err
	}
	resp := OnlinePushMessageResponse{
		Uin:      uin,
		Items:    items,
		ServerIP: serverIP,
	}
	buf, err := uni.Marshal(ctx, &uni.Message{
		Version:     0x0003,
		PacketType:  0x00,
		MessageType: 0x00000000,
		RequestID:   seq,
		ServantName: "OnlinePush",
		FuncName:    "SvcRespPushMsg",
		Buffer:      []byte{},
		Timeout:     0x00000000,
		Context:     map[string]string{},
		Status:      map[string]string{},
	}, map[string]interface{}{
		"resp": resp,
	})
	if err != nil {
		return nil, err
	}
	return &codec.ClientToServerMessage{
		Username:      username,
		Seq:           uint32(seq),
		ServiceMethod: ServiceMethodOnlinePushResponse,
		Buffer:        buf,
		Simple:        false,
	}, nil
}

func (c *Client) handleOnlinePushRequest(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error) {
	msg := uni.Message{}
	push := OnlinePushRequest{}
	if err := uni.Unmarshal(ctx, s2c.Buffer, &msg, map[string]interface{}{
		"req": &push,
	}); err != nil {
		return nil, err
	}
	util.DumpServerToClientMessage(s2c, &push)

	uin, _ := strconv.ParseUint(s2c.Username, 10, 64)
	items := []MessageDelete{}
	for _, msg := range push.Messages {
		switch msg.MessageType {
		case 0x0210:
			body, err := c.decodeMessageType0210Jce(uin, msg.MessageBytes)
			if err != nil {
				return nil, err
			} else if body != nil {
				util.DumpServerToClientMessage(s2c, &body)
			}
		case 0x02DC:
			body, err := c.decodeMessageType02DC(uin, msg.MessageBytes)
			if err != nil {
				return nil, err
			} else if body != nil {
				util.DumpServerToClientMessage(s2c, &body)
			}
		}
		items = append(items, MessageDelete{
			FromUin:       msg.FromUin,
			MessageTime:   msg.MessageTime,
			MessageSeq:    msg.MessageSeq,
			MessageCookie: msg.MessageCookie,
		})
	}

	return NewOnlinePushMessageResponse(ctx, s2c.Username, items, push.ServerIP, int32(s2c.Seq))
}
