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
	MessageTime   int64  `jce:",1" json:",omitempty"`
	MessageSeq    uint16 `jce:",2" json:",omitempty"`
	MessageCookie []byte `jce:",3" json:",omitempty"`
}

func NewOnlinePushMessageResponse(
	ctx context.Context,
	username string,
	infos []MessageDeleteInfo,
	serverIP int32,
	seq uint32,
) (*codec.ClientToServerMessage, error) {
	if len(infos) == 0 {
		return nil, fmt.Errorf("zero length")
	}

	uin, err := strconv.ParseInt(username, 10, 64)
	if err != nil {
		return nil, err
	}
	resp := OnlinePushMessageResponse{
		Uin:      uint64(uin),
		Infos:    infos,
		ServerIP: serverIP,
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

func (c *Client) handleOnlinePushMessage(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error) {
	push := pb.OnlinePush{}
	if err := proto.Unmarshal(s2c.Buffer, &push); err != nil {
		return nil, err
	}
	c.dumpServerToClientMessage(s2c, &push)

	msg := push.GetMessage()
	data, err := mark.Marshal(msg)
	if err != nil {
		return nil, err
	}

	chatID := msg.GetMessageHead().GetGroupInfo().GetGroupCode()
	peerID := uint64(0)
	fromID := msg.GetMessageHead().GetFromUin()
	if msg.GetMessageHead().GetC2CCmd() != 0 {
		chatID = msg.GetMessageHead().GetC2CTempMessageHead().GetGroupCode()
		peerID = msg.GetMessageHead().GetToUin()
	}
	chatName := strconv.FormatUint(chatID, 10)
	peerName := strconv.FormatUint(peerID, 10)
	fromName := strconv.FormatUint(fromID, 10)
	seq := msg.GetMessageHead().GetMessageSeq()
	text := string(data)

	syncMessage(msg)
	log.PrintMessage(
		time.Unix(int64(msg.GetMessageHead().GetMessageTime()), 0),
		chatName, peerName, fromName, chatID, peerID, uint64(fromID), seq, text,
	)

	if s2c.Username != strconv.FormatInt(int64(fromID), 10) {
		id := fmt.Sprintf("@%d_%d", chatID, peerID)
		routingHead := &pb.RoutingHead{}
		if msg.GetMessageHead().GetC2CCmd() == 0 {
			c.setMessageSeq(id, msg.GetMessageHead().GetMessageSeq())
			routingHead = &pb.RoutingHead{Group: &pb.Group{Code: chatID}}
		} else if chatID == 0 {
			routingHead = &pb.RoutingHead{C2C: &pb.C2C{ToUin: peerID}}
		} else {
			routingHead = &pb.RoutingHead{
				GroupTemp: &pb.GroupTemp{Uin: chatID, ToUin: peerID},
			}
		}
		seq = c.getNextMessageSeq(id)

		msg := pb.Message{}
		if err := mark.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		resp, err := c.MessageSendMessage(
			ctx, s2c.Username, NewMessageSendMessageRequest(
				routingHead,
				msg.GetContentHead(),
				msg.GetMessageBody(),
				seq,
				c.syncCookie,
			),
		)
		if err != nil {
			return nil, err
		}

		data, err := mark.Marshal(&msg)
		if err != nil {
			return nil, err
		}
		fromID, _ = strconv.ParseUint(s2c.Username, 10, 64)
		fromName = s2c.Username
		text = string(data)
		log.PrintMessage(
			time.Unix(resp.GetSendTime(), 0),
			chatName, peerName, fromName, chatID, peerID, uint64(fromID), seq, text,
		)
	}

	return NewOnlinePushMessageResponse(ctx, s2c.Username, []MessageDeleteInfo{{
		FromUin:     fromID,
		MessageTime: msg.GetMessageHead().GetMessageTime(),
		MessageSeq:  uint16(seq),
	}}, push.GetServerIp(), s2c.Seq)
}
