package rpc

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/pb"
	"google.golang.org/protobuf/proto"
)

func (c *Client) handlePushOnlineMessage(
	ctx context.Context,
	s2c *ServerToClientMessage,
) (*ClientToServerMessage, error) {
	msg := pb.OnlinePushMessage{}
	if err := proto.Unmarshal(s2c.Buffer, &msg); err != nil {
		return nil, err
	}
	c.dumpServerToClientMessage(s2c, &msg)
	_, _ = c.marshalMessage(msg.GetMessage())
	return nil, nil
}
