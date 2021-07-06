package client

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

func (c *Client) handleOnlinePushSIDTicketExpired(
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
	log.Warn().Msg("<-> [todo] OnlinePushSIDTicketExpired, user SID ticket needs to be update")
	return &codec.ClientToServerMessage{
		Username:      s2c.Username,
		Seq:           s2c.Seq,
		ServiceMethod: ServiceMethodOnlinePushSIDTicketExpired,
		Buffer:        nil,
		Simple:        true,
	}, nil
}
