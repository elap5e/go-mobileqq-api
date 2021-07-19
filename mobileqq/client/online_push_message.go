package client

import (
	"context"
	"strconv"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
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

func (c *Client) handleOnlinePushMessage(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error) {
	push := pb.OnlinePush{}
	if err := proto.Unmarshal(s2c.Buffer, &push); err != nil {
		return nil, err
	}
	dumpServerToClientMessage(s2c, &push)

	uin, _ := strconv.ParseUint(s2c.Username, 10, 64)
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
			c.setMessageSeq(mr.PeerID, mr.UserID, int64(uin), mr.Seq)
		}
	}
	text, _ := mark.NewEncoder(mr.PeerID, mr.UserID, mr.FromID).
		Encode(msg.GetMessageBody().GetRichText().GetElements())
	mr.Text = string(text)

	c.PrintMessageRecord(mr)
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

		elems, err := mark.NewDecoder(mr.PeerID, mr.UserID, mr.FromID).
			Decode([]byte(mr.Text))
		if err != nil {
			return nil, err
		}
		msg := pb.Message{
			MessageBody: &pb.MessageBody{
				RichText: &pb.RichText{
					Elements: elems,
				},
			},
		}
		if _, err := c.MessageSendMessage(
			ctx, s2c.Username, NewMessageSendMessageRequest(
				routingHead,
				msg.GetContentHead(),
				msg.GetMessageBody(),
				0,
				c.syncCookie[int64(uin)],
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
