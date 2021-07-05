package client

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

func (c *Client) handlePushOnlineSIDExpired(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error) {
	// go func() {
	// 	time.Sleep(3000 * time.Millisecond)
	// 	_, _ = c.AuthGetSessionTicketsWithoutPassword(
	// 		ctx,
	// 		NewAuthGetSessionTicketsWithoutPasswordRequest(s2c.Username, false),
	// 	)
	// }()
	return &codec.ClientToServerMessage{
		Username:      s2c.Username,
		Seq:           s2c.Seq,
		ServiceMethod: ServiceMethodPushOnlineSIDExpired,
		Buffer:        nil,
		Simple:        true,
	}, nil
}
