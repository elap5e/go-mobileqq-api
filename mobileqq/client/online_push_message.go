package client

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
)

type OnlinePushMessageResponse struct {
	Uin         uint64              `jce:",0" json:",omitempty"`
	DeleteInfos []MessageDeleteInfo `jce:",1" json:",omitempty"`
	Svrip       int32               `jce:",2" json:",omitempty"`
}

type MessageDeleteInfo struct {
	FromUin       uint64 `jce:",0" json:",omitempty"`
	MessageTime   uint64 `jce:",1" json:",omitempty"`
	MessageSeq    uint16 `jce:",2" json:",omitempty"`
	MessageCookie []byte `jce:",3" json:",omitempty"`
}

func (c *Client) handleOnlinePushMessage(
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
	if s2c.ServiceMethod == ServiceMethodOnlinePushMessageSyncC2C {
		uin, _ := strconv.ParseInt(s2c.Username, 10, 64)
		resp := OnlinePushMessageResponse{
			Uin: uint64(uin),
			DeleteInfos: []MessageDeleteInfo{{
				FromUin:     fromUin,
				MessageTime: uint64(msg.GetMessageHead().GetMessageTime()),
				MessageSeq:  uint16(msg.GetMessageHead().GetMessageSeq()),
			}},
			Svrip: push.GetSvrip(),
		}
		buf, err := uni.Marshal(ctx, &uni.Message{
			Version:     0x0003,
			PacketType:  0x00,
			MessageType: 0x00000000,
			RequestID:   s2c.Seq,
			ServantName: "OnlinePush",
			FuncName:    "SvcRespPushMsg",
			Buffer:      []byte{},
			Timeout:     0x00000000,
			Context:     map[string]string{},
			Status:      map[string]string{},
		}, map[string]interface{}{
			"resp": resp,
		})
		if err != nil {
			return nil, err
		}
		return &codec.ClientToServerMessage{
			Username:      s2c.Username,
			Seq:           s2c.Seq,
			ServiceMethod: ServiceMethodOnlinePushResponse,
			Buffer:        buf,
			Simple:        false,
		}, nil
	}
	return nil, nil
}
