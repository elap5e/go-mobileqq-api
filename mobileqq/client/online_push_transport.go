package client

import (
	"context"
	"encoding/hex"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
	"github.com/elap5e/go-mobileqq-api/util"
	"google.golang.org/protobuf/proto"
)

func (c *Client) handleOnlinePushTransport(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error) {
	push := pb.OnlinePush_Transport{}
	if err := proto.Unmarshal(s2c.Buffer, &push); err != nil {
		log.Debug().Msg(">>> [dump]\n" + hex.Dump(s2c.Buffer))
		return nil, err
	}
	util.DumpServerToClientMessage(s2c, &push)

	return NewOnlinePushMessageResponse(ctx, s2c.Username, []MessageDelete{{
		FromUin:     int64(push.GetFromUin()),
		MessageTime: int64(push.GetMessageTime()),
		MessageSeq:  int32(push.GetMessageSeq()),
	}}, Uint32IPType(push.GetServerIp()), int32(s2c.Seq))
}
