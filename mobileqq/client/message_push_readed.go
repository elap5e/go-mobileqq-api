package client

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

type MessagePushReadedRequest struct {
	Type       uint8                     `jce:",0" json:",omitempty"`
	C2C        []MessageReadedC2C        `jce:",1" json:",omitempty"`
	Group      []MessageReadedGroup      `jce:",2" json:",omitempty"`
	Discussion []MessageReadedDiscussion `jce:",3" json:",omitempty"`
}

type MessageReadedC2C struct {
	Uin          uint64 `jce:",0" json:",omitempty"`
	LastReadTime uint64 `jce:",1" json:",omitempty"`
}

type MessageReadedGroup struct {
	PeerUin    uint64 `jce:",0" json:",omitempty"`
	Type       uint64 `jce:",1" json:",omitempty"`
	MemberSeq  uint64 `jce:",2" json:",omitempty"`
	MessageSeq uint64 `jce:",3" json:",omitempty"`
}

type MessageReadedDiscussion struct {
	PeerUin    uint64 `jce:",0" json:",omitempty"`
	Type       uint64 `jce:",1" json:",omitempty"`
	MemberSeq  uint64 `jce:",2" json:",omitempty"`
	MessageSeq uint64 `jce:",3" json:",omitempty"`
}

func (c *Client) handleMessagePushReaded(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error) {
	msg := uni.Message{}
	req := MessagePushReadedRequest{}
	if err := uni.Unmarshal(ctx, s2c.Buffer[4:], &msg, map[string]interface{}{
		"req": &req,
	}); err != nil {
		return nil, err
	}
	c.dumpServerToClientMessage(s2c, &req)
	// TODO: handle

	return nil, nil
}
