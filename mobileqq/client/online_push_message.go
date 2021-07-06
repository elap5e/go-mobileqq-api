package client

import (
	"context"
	"fmt"
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
	Uin   uint64              `jce:",0" json:",omitempty"`
	Infos []MessageDeleteInfo `jce:",1" json:",omitempty"`
	Svrip int32               `jce:",2" json:",omitempty"`
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
	req := pb.OnlinePushMessage{}
	if err := proto.Unmarshal(s2c.Buffer, &req); err != nil {
		return nil, err
	}
	c.dumpServerToClientMessage(s2c, &req)
	msg := req.GetMessage()
	data, err := c.marshalMessage(msg)
	if err != nil {
		return nil, err
	}
	infoList := []MessageDeleteInfo{}
	peerUin := msg.GetMessageHead().GetGroupInfo().GetGroupCode()
	fromUin := msg.GetMessageHead().GetFromUin()
	if s2c.Username != strconv.FormatInt(int64(fromUin), 10) {
		infoList = append(infoList, MessageDeleteInfo{
			FromUin:     fromUin,
			MessageTime: uint64(msg.GetMessageHead().GetMessageTime()),
			MessageSeq:  uint16(msg.GetMessageHead().GetMessageSeq()),
		})

		c.setSyncSeq(peerUin, msg.GetMessageHead().GetMessageSeq())
		msg := pb.Message{}
		if err := mark.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		seq := c.getNextSyncSeq(peerUin)
		log.Info().
			Str("@peer", fmt.Sprintf("%d:%s:%d", peerUin, s2c.Username, fromUin)).
			Uint32("@seq", seq).
			Int64("@time", time.Now().Unix()).
			Str("mark", string(data)).
			Msg("<-- [send] message")
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
	return NewOnlinePushMessageResponse(ctx, s2c.Username, infoList, req.GetSvrip(), s2c.Seq)
}

func NewOnlinePushMessageResponse(
	ctx context.Context,
	username string,
	infos []MessageDeleteInfo,
	svrip int32,
	seq uint32,
) (*codec.ClientToServerMessage, error) {
	if len(infos) > 0 {
		uin, err := strconv.ParseInt(username, 10, 64)
		if err != nil {
			return nil, err
		}
		resp := OnlinePushMessageResponse{
			Uin:   uint64(uin),
			Infos: infos,
			Svrip: svrip,
		}
		buf, err := uni.Marshal(ctx, &uni.Message{
			Version:     0x0003,
			PacketType:  0x00,
			MessageType: 0x00000000,
			RequestID:   seq,
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
			Username:      username,
			Seq:           seq,
			ServiceMethod: ServiceMethodOnlinePushResponse,
			Buffer:        buf,
			Simple:        false,
		}, nil
	}
	return nil, nil
}
