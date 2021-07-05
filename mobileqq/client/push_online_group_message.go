package client

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
)

func (c *Client) handlePushOnlineGroupMessage(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error) {
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
		log.Info().
			Str("@mark", string(data)).
			Str("from", s2c.Username).
			Uint64("peer", peerUin).
			Uint32("seq", seq).
			Uint64("to", fromUin).
			Int64("time", time.Now().Unix()).
			Msg("<== [send] message")
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
