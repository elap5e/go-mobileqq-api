package mobileqq

import (
	"github.com/spf13/viper"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/pb"
	"github.com/elap5e/go-mobileqq-api/rpc"
)

func (c *Client) MessageSendMessage(username string, text string) error {
	msg := pb.Message{}
	if err := mark.Unmarshal([]byte(text), &msg); err != nil {
		return err
	}
	seq := c.rpc.GetNextSyncSeq(0)
	_, _ = c.rpc.MessageSendMessage(
		c.ctx, username, rpc.NewMessageSendMessageRequest(
			&pb.RoutingHead{C2C: &pb.C2C{Uin: viper.GetUint64("targets.0.uin")}},
			msg.GetContentHead(),
			msg.GetMessageBody(),
			seq,
			nil,
		),
	)
	return nil
}

func (c *Client) MessageGetMessage(username string) error {
	_, _ = c.rpc.MessageGetMessage(
		c.ctx, username, rpc.NewMessageGetMessageRequest(
			0x00000000, nil,
		),
	)
	return nil
}
