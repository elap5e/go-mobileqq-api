package client

import (
	"context"
	"encoding/hex"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

func (c *Client) handleOnlinePushSIDTicketExpired(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error) {
	log.Warn().Msg("<-> [todo] OnlinePushSIDTicketExpired, user SID ticket needs to be update")
	log.Debug().Msg(">>> [dump]\n" + hex.Dump(s2c.Buffer))
	return &codec.ClientToServerMessage{
		Username:      s2c.Username,
		Seq:           s2c.Seq,
		ServiceMethod: ServiceMethodOnlinePushSIDTicketExpired,
		Buffer:        nil,
		Simple:        true,
	}, nil
}
