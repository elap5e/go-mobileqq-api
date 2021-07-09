package client

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

type OnlinePushRequest struct {
	Type   uint32 `jce:",1" json:",omitempty"`
	Seq    uint64 `jce:",3" json:",omitempty"`
	Buffer []byte `jce:",2" json:",omitempty"`
}

type OnlinePushResponse struct {
	Type   uint32 `jce:",1" json:",omitempty"`
	Seq    uint64 `jce:",2" json:",omitempty"`
	Buffer []byte `jce:",3" json:",omitempty"`
}

func (c *Client) handleOnlinePushRequest(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error) {
	msg := uni.Message{}
	req := OnlinePushRequest{}
	if err := uni.Unmarshal(ctx, s2c.Buffer, &msg, map[string]interface{}{
		"req": &req,
	}); err != nil {
		return nil, err
	}
	c.dumpServerToClientMessage(s2c, &req)
	// TODO: handle

	return nil, nil
}
