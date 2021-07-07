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
	Uin      uint64              `jce:",0" json:",omitempty"`
	Infos    []MessageDeleteInfo `jce:",1" json:",omitempty"`
	ServerIP int32               `jce:",2" json:",omitempty"`
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
	data, err := mark.Marshal(msg)
	if err != nil {
		return nil, err
	}

	chatID := msg.GetMessageHead().GetGroupInfo().GetGroupCode()
	peerID := uint64(0)
	fromID := int(msg.GetMessageHead().GetFromUin())
	chatName := strconv.Itoa(int(chatID))
	peerName := strconv.Itoa(int(peerID))
	fromName := strconv.Itoa(int(fromID))
	seq := msg.GetMessageHead().GetMessageSeq()
	text := string(data)

	log.Debug().
		Str("chat", fmt.Sprintf("%d:%d", chatID, peerID)).
		Int("from", fromID).
		Uint32("seq", seq).
		Uint32("time", msg.GetMessageHead().GetMessageTime()).
		Uint32("type", msg.GetMessageHead().GetMessageType()).
		Uint64("uid", msg.GetMessageHead().GetMessageUid()).
		Msg("--> [recv]")

	infoList := []MessageDeleteInfo{}
	if s2c.Username != strconv.FormatInt(int64(fromID), 10) {
		log.PrintMessage(
			time.Unix(int64(msg.GetMessageHead().GetMessageTime()), 0),
			chatName, peerName, fromName, chatID, peerID, uint64(fromID), seq, text,
		)

		infoList = append(infoList, MessageDeleteInfo{
			FromUin:     uint64(fromID),
			MessageTime: uint64(msg.GetMessageHead().GetMessageTime()),
			MessageSeq:  uint16(seq),
		})

		c.setSyncSeq(chatID, seq)
		msg := pb.Message{}
		if err := mark.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		seq = c.getNextSyncSeq(chatID)
		resp, err := c.MessageSendMessage(
			ctx, s2c.Username, NewMessageSendMessageRequest(
				&pb.RoutingHead{Group: &pb.Group{Code: chatID}},
				msg.GetContentHead(),
				msg.GetMessageBody(),
				seq,
				c.syncCookie,
			),
		)
		if err != nil {
			return nil, err
		}
		fromName = s2c.Username
		fromID, _ = strconv.Atoi(s2c.Username)
		data, err := mark.Marshal(&msg)
		if err != nil {
			return nil, err
		}
		text = string(data)
		log.PrintMessage(
			time.Unix(int64(resp.GetSendTime()), 0),
			chatName, peerName, fromName, chatID, peerID, uint64(fromID), seq, text,
		)
	}
	return NewOnlinePushMessageResponse(ctx, s2c.Username, infoList, req.GetServerIp(), s2c.Seq)
}

func NewOnlinePushMessageResponse(
	ctx context.Context,
	username string,
	infos []MessageDeleteInfo,
	ip int32,
	seq uint32,
) (*codec.ClientToServerMessage, error) {
	if len(infos) > 0 {
		uin, err := strconv.ParseInt(username, 10, 64)
		if err != nil {
			return nil, err
		}
		resp := OnlinePushMessageResponse{
			Uin:      uint64(uin),
			Infos:    infos,
			ServerIP: ip,
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
