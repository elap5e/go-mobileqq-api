package client

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
	"github.com/elap5e/go-mobileqq-api/util"
)

func (c *Client) handleMessageSendMessageResponse(s2c *codec.ServerToClientMessage, req *pb.MessageService_SendRequest, resp *pb.MessageService_SendResponse) {
	util.DumpServerToClientMessage(s2c, resp)

	uin, _ := strconv.ParseUint(s2c.Username, 10, 64)
	mr := &db.MessageRecord{
		Time:   resp.GetSendTime(),
		Seq:    req.GetMessageSeq(),
		Uid:    int64(req.GetMessageRand()) | 1<<56,
		PeerID: req.GetRoutingHead().GetGroup().GetGroupCode(),
		UserID: req.GetRoutingHead().GetC2C().GetToUin(),
		FromID: int64(uin),
		Text:   string(""),
		Type:   0,
	}
	if mr.UserID != 0 {
		mr.Type = 166
	} else if mr.PeerID != 0 {
		mr.Type = 82
	} else {
		mr.PeerID = req.GetRoutingHead().GetGroupTemp().GetGroupUin()
		mr.UserID = req.GetRoutingHead().GetGroupTemp().GetToUin()
		mr.Type = 141
	}
	text, _ := mark.NewEncoder(mr.PeerID, mr.UserID, mr.FromID).
		Encode(req.GetMessageBody().GetRichText().GetElements())
	mr.Text = string(text)

	if resp.Result == 0 {
		if c.db != nil {
			err := c.dbInsertMessageRecord(uin, mr)
			if err != nil {
				log.Error().Err(err).Msg(">>> [db  ] dbInsertMessageRecord")
			}
		} else {
			c.PrintMessageRecord(mr)
		}
	} else {
		mr.Text += " " + log.Colorize(fmt.Sprintf("(%d)", resp.Result), log.ColorBrightRed, false)
		c.PrintMessageRecord(mr)
	}
}

func NewMessageSendMessageRequest(
	routingHead *pb.MessageService_RoutingHead,
	contentHead *pb.MessageCommon_ContentHead,
	messageBody *pb.IMMessageBody_MessageBody,
	seq int32,
	cookie []byte,
) *pb.MessageService_SendRequest {
	return &pb.MessageService_SendRequest{
		RoutingHead: routingHead,
		ContentHead: contentHead,
		MessageBody: messageBody,
		MessageSeq:  seq,
		MessageRand: 0,
		SyncCookie:  cookie,
	}
}

func (c *Client) MessageSendMessage(
	ctx context.Context,
	username string,
	req *pb.MessageService_SendRequest,
) (*pb.MessageService_SendResponse, error) {
	uin, _ := strconv.ParseInt(username, 10, 64)
	if req.GetMessageSeq() == 0 {
		var peerID, userID int64
		if req.GetRoutingHead().GetC2C() != nil {
			userID = req.GetRoutingHead().GetC2C().GetToUin()
		} else if req.GetRoutingHead().GetGroup() != nil {
			peerID = req.GetRoutingHead().GetGroup().GetGroupCode()
		} else if req.GetRoutingHead().GetC2C() != nil {
			peerID = req.GetRoutingHead().GetGroupTemp().GetGroupUin()
			userID = req.GetRoutingHead().GetGroupTemp().GetToUin()
		}
		req.MessageSeq = c.getNextMessageSeq(peerID, userID, uin)
	}
	if req.GetMessageRand() == 0 {
		req.MessageRand = rand.Int31()
	}
	if len(req.GetSyncCookie()) == 0 {
		req.SyncCookie = c.syncCookie[uin]
	}

	buf, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	c2s, s2c := codec.ClientToServerMessage{
		Username: username,
		Buffer:   buf,
		Simple:   true,
	}, codec.ServerToClientMessage{}
	err = c.rpc.Call(ServiceMethodMessageSendMessage, &c2s, &s2c)
	if err != nil {
		return nil, err
	}
	resp := pb.MessageService_SendResponse{}
	if err := proto.Unmarshal(s2c.Buffer, &resp); err != nil {
		return nil, err
	}
	// resp.Result
	//     0: success
	//     1: ???
	//    16: elements (notFriend)
	//   120: elements (groupMute)
	//   241: ???
	// -3902: marketFace (vip/svip)
	// -4902: marketFace magic (vip/svip)
	//  5002: poke (vip/svip)

	c.handleMessageSendMessageResponse(&s2c, req, &resp)
	return &resp, nil
}
