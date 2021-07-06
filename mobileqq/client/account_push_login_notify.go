package client

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

func (c *Client) handleAccountPushLoginNotify(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error) {
	return nil, nil
}
