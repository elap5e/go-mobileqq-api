package rpc

import (
	"context"
)

func (c *Client) handlePushOnlineSIDExpired(
	ctx context.Context,
	s2c *ServerToClientMessage,
) (*ClientToServerMessage, error) {
	// _, _ = c.AuthGetSessionTicketsWithoutPassword(
	// 	ctx,
	// 	NewAuthGetSessionTicketsWithoutPasswordRequest(s2c.Username),
	// )
	return &ClientToServerMessage{
		Username:      s2c.Username,
		Seq:           s2c.Seq,
		ServiceMethod: ServiceMethodPushOnlineSIDExpired,
		Buffer:        nil,
		Simple:        true,
	}, nil
}
