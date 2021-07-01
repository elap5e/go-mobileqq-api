package rpc

import (
	"context"
)

func (c *Client) handlePushOnlineSIDTicketExpired(
	ctx context.Context,
	s2c *ServerToClientMessage,
) (*ClientToServerMessage, error) {
	return &ClientToServerMessage{
		Username:      s2c.Username,
		Seq:           s2c.Seq,
		ServiceMethod: ServiceMethodPushOnlineSIDTicketExpired,
		Buffer:        nil,
		Simple:        true,
	}, nil
}
