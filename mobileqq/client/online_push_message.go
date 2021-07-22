package client

import (
	"context"
	"encoding/hex"
	"strconv"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
	"github.com/elap5e/go-mobileqq-api/util"
)

type OnlinePushMessageResponse struct {
	Uin      int64           `jce:",0" json:"uin,omitempty"`
	Items    []MessageDelete `jce:",1" json:"items,omitempty"`
	ServerIP Uint32IPType    `jce:",2" json:"server_ip,omitempty"`
}

func (c *Client) handleOnlinePushMessage(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error) {
	push := pb.OnlinePush_Message{}
	if err := proto.Unmarshal(s2c.Buffer, &push); err != nil {
		log.Debug().Msg(">>> [dump]\n" + hex.Dump(s2c.Buffer))
		return nil, err
	}
	util.DumpServerToClientMessage(s2c, &push)

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

	if c.db != nil {
		err := c.dbInsertMessageRecord(uin, mr)
		if err != nil {
			log.Error().Err(err).Msg(">>> [db  ] dbInsertMessageRecord")
		}
	} else {
		c.PrintMessageRecord(mr)
	}

	if uin != uint64(mr.FromID) {
		elems, err := mark.NewDecoder(mr.PeerID, mr.UserID, mr.FromID).
			Decode([]byte(mr.Text))
		if err != nil {
			return nil, err
		}
		msg := pb.MessageCommon_Message{
			MessageBody: &pb.IMMessageBody_MessageBody{
				RichText: &pb.IMMessageBody_RichText{
					Elements: elems,
				},
			},
		}
		if _, err := c.MessageSendMessage(
			ctx, s2c.Username, NewMessageSendMessageRequest(
				c.GetRoutingHead(mr.PeerID, mr.UserID),
				msg.GetContentHead(),
				msg.GetMessageBody(),
				0,
				c.syncCookie[int64(uin)],
			),
		); err != nil {
			return nil, err
		}
	}

	return NewOnlinePushMessageResponse(ctx, s2c.Username, []MessageDelete{{
		FromUin:     mr.FromID,
		MessageTime: mr.Time,
		MessageSeq:  mr.Seq,
	}}, Uint32IPType(push.GetServerIp()), int32(s2c.Seq))
}
