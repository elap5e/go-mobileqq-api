package client

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
)

type OnlinePushMessageResponse struct {
	Uin      int64               `jce:",0" json:",omitempty"`
	Infos    []MessageDeleteInfo `jce:",1" json:",omitempty"`
	ServerIP uint32              `jce:",2" json:",omitempty"`
}

type MessageDeleteInfo struct {
	FromUin       int64  `jce:",0" json:",omitempty"`
	MessageTime   int64  `jce:",1" json:",omitempty"`
	MessageSeq    int32  `jce:",2" json:",omitempty"`
	MessageCookie []byte `jce:",3" json:",omitempty"`
}

func NewOnlinePushMessageResponse(
	ctx context.Context,
	username string,
	infos []MessageDeleteInfo,
	serverIP uint32,
	seq int32,
) (*codec.ClientToServerMessage, error) {
	if len(infos) == 0 {
		return nil, fmt.Errorf("zero length")
	}

	uin, err := strconv.ParseInt(username, 10, 64)
	if err != nil {
		return nil, err
	}
	resp := OnlinePushMessageResponse{
		Uin:      uin,
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
		Seq:           uint32(seq),
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
	dumpServerToClientMessage(s2c, &push)

	msg := push.GetMessage()
	mr := &db.MessageRecord{
		Time:   msg.GetMessageHead().GetMessageTime(),
		Seq:    msg.GetMessageHead().GetMessageSeq(),
		Uid:    int64(msg.GetMessageBody().GetRichText().GetAttribute().GetRandom()) | 1<<56,
		PeerID: msg.GetMessageHead().GetC2CTempMessageHead().GetGroupCode(),
		UserID: msg.GetMessageHead().GetToUin(),
		FromID: msg.GetMessageHead().GetFromUin(),
		Text:   "",
		Type:   msg.GetMessageHead().GetMessageType(),
	}
	if msg.GetMessageHead().GetC2CCmd() == 0 {
		mr.PeerID = msg.GetMessageHead().GetGroupInfo().GetGroupCode()
		mr.UserID = 0
		if mr.Type == 82 {
			c.setMessageSeq(fmt.Sprintf("@%du%d", mr.PeerID, mr.UserID), mr.Seq)
		}
	}
	text, _ := mark.NewMarshaler(mr.PeerID, mr.UserID, mr.FromID).
		Marshal(msg.GetMessageBody().GetRichText().GetElements())
	mr.Text = string(text)

	c.PrintMessageRecord(mr)
	uin, _ := strconv.ParseUint(s2c.Username, 10, 64)
	if c.db != nil {
		err := c.dbInsertMessageRecord(uin, mr)
		if err != nil {
			log.Error().Err(err).Msg(">>> [db  ] dbInsertMessageRecord")
		}
	}

	if uin != uint64(mr.FromID) {
		routingHead := &pb.RoutingHead{}
		if mr.PeerID == 0 {
			routingHead = &pb.RoutingHead{C2C: &pb.C2C{ToUin: mr.UserID}}
		} else if mr.UserID == 0 {
			routingHead = &pb.RoutingHead{Group: &pb.Group{Code: mr.PeerID}}
		} else {
			routingHead = &pb.RoutingHead{
				GroupTemp: &pb.GroupTemp{Uin: mr.PeerID, ToUin: mr.UserID},
			}
		}
		chatID := fmt.Sprintf("@%du%d", mr.PeerID, mr.UserID)
		seq := c.getNextMessageSeq(chatID)

		msg := pb.Message{}
		if err := mark.Unmarshal([]byte(mr.Text), &msg); err != nil {
			return nil, err
		}
		if _, err := c.MessageSendMessage(
			ctx, s2c.Username, NewMessageSendMessageRequest(
				routingHead,
				msg.GetContentHead(),
				msg.GetMessageBody(),
				seq,
				c.syncCookie,
			),
		); err != nil {
			return nil, err
		}
	}

	return NewOnlinePushMessageResponse(ctx, s2c.Username, []MessageDeleteInfo{{
		FromUin:     mr.FromID,
		MessageTime: mr.Time,
		MessageSeq:  mr.Seq,
	}}, push.GetServerIp(), int32(s2c.Seq))
}
