package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

type OnlinePushRequest struct {
	Uin             int64             `jce:",0" json:",omitempty"`
	MessageTime     int64             `jce:",1" json:",omitempty"`
	MessageInfos    []*MessageInfo    `jce:",2" json:",omitempty"`
	ServerIP        uint32            `jce:",3" json:",omitempty"`
	SyncCookie      []byte            `jce:",4" json:",omitempty"`
	UinPairMessages []*UinPairMessage `jce:",5" json:",omitempty"`
	Previews        map[string][]byte `jce:",6" json:",omitempty"`
	UserActive      int32             `jce:",7" json:",omitempty"`
	GeneralFlag     int32             `jce:",12" json:",omitempty"`
}

type UinPairMessage struct {
	LastReadTime     int64          `jce:",1" json:",omitempty"`
	PeerUin          int64          `jce:",2" json:",omitempty"`
	MessageCompleted int64          `jce:",3" json:",omitempty"`
	MessageInfos     []*MessageInfo `jce:",4" json:",omitempty"`
}

type OnlinePushResponse struct {
	Type   int32  `jce:",1" json:",omitempty"`
	Seq    int64  `jce:",2" json:",omitempty"`
	Buffer []byte `jce:",3" json:",omitempty"`
}

func NewOnlinePushMessageResponse(
	ctx context.Context,
	username string,
	infos []MessageDeleteInfo,
	serverIP uint32,
	seq int32,
) (*codec.ClientToServerMessage, error) {
	if len(infos) == 0 {
		return nil, fmt.Errorf("zero length")
	}

	uin, err := strconv.ParseInt(username, 10, 64)
	if err != nil {
		return nil, err
	}
	resp := OnlinePushMessageResponse{
		Uin:      uin,
		Infos:    infos,
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
	dumpServerToClientMessage(s2c, &push)

	uin, _ := strconv.ParseUint(s2c.Username, 10, 64)
	infos := []MessageDeleteInfo{}
	for _, msg := range push.MessageInfos {
		switch msg.MessageType {
		case 0x0210:
			body, err := c.decodeMessageType0210(uin, msg.MessageBytes)
			if err != nil {
				return nil, err
			} else if body != nil {
				dumpServerToClientMessage(s2c, &body)
			}
		case 0x02DC:
			body, err := c.decodeMessageType02DC(uin, msg.MessageBytes)
			if err != nil {
				return nil, err
			} else if body != nil {
				dumpServerToClientMessage(s2c, &body)
			}
		}
		infos = append(infos, MessageDeleteInfo{
			FromUin:       msg.FromUin,
			MessageTime:   msg.MessageTime,
			MessageSeq:    msg.MessageSeq,
			MessageCookie: msg.MessageCookies,
		})
	}

	return NewOnlinePushMessageResponse(ctx, s2c.Username, infos, push.ServerIP, int32(s2c.Seq))
}
