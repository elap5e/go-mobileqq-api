package rpc

import (
	"context"
	"log"
	"strconv"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/pb"
)

func (c *Client) handlePushOnlineGroupMessage(
	ctx context.Context,
	s2c *ServerToClientMessage,
) (*ClientToServerMessage, error) {
	push := pb.OnlinePushMessage{}
	if err := proto.Unmarshal(s2c.Buffer, &push); err != nil {
		return nil, err
	}
	c.dumpServerToClientMessage(s2c, &push)
	msg := push.GetMessage()
	data, err := c.marshalMessage(msg)
	if err != nil {
		return nil, err
	}
	peerUin := msg.GetMessageHead().GetGroupInfo().GetGroupCode()
	fromUin := msg.GetMessageHead().GetFromUin()
	if s2c.Username != strconv.FormatInt(int64(fromUin), 10) {
		c.setSyncSeq(peerUin, msg.GetMessageHead().GetMessageSeq())
		msg := pb.Message{}
		if err := mark.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		seq := c.getNextSyncSeq(peerUin)
		if c.cfg.Debug {
			log.Printf(
				"<<< [dump] peer:%d seq:%d from:%s to:%d markdown:\n%s",
				peerUin, seq, s2c.Username, fromUin, string(data),
			)
		}
		_, _ = c.MessageSendMessage(
			ctx, s2c.Username, NewMessageSendMessageRequest(
				&pb.RoutingHead{Group: &pb.Group{Code: peerUin}},
				msg.GetContentHead(),
				msg.GetMessageBody(),
				seq,
				c.syncCookie,
			),
		)
	}
	return nil, nil
}
