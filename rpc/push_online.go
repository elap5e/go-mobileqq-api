package rpc

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/encoding/uni"
)

type PushOnlineRequest struct {
	Type   uint32 `jce:",1" json:",omitempty"`
	Seq    uint64 `jce:",3" json:",omitempty"`
	Buffer []byte `jce:",2" json:",omitempty"`
}

type PushOnlineResponse struct {
	Type   uint32 `jce:",1" json:",omitempty"`
	Seq    uint64 `jce:",2" json:",omitempty"`
	Buffer []byte `jce:",3" json:",omitempty"`
}

func (c *Client) handlePushOnlineRequest(
	ctx context.Context,
	s2c *ServerToClientMessage,
) (*ClientToServerMessage, error) {
	msg := uni.Message{}
	req := PushOnlineRequest{}
	if err := uni.Unmarshal(ctx, s2c.Buffer, &msg, map[string]interface{}{
		"req": &req,
	}); err != nil {
		return nil, err
	}
	return nil, nil
}
