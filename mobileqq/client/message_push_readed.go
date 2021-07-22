package client

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/util"
)

type MessagePushReadedRequest struct {
	Type    uint8                  `jce:",0" json:"type,omitempty"`
	C2C     []MessageReadedC2C     `jce:",1" json:"c2c,omitempty"`
	Group   []MessageReadedGroup   `jce:",2" json:"group,omitempty"`
	Discuss []MessageReadedDiscuss `jce:",3" json:"discuss,omitempty"`
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
	util.DumpServerToClientMessage(s2c, &req)
	// TODO: handle

	return nil, nil
}
