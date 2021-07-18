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
)

func (c *Client) handleMessageSendMessageResponse(s2c *codec.ServerToClientMessage, req *pb.MessageSendMessageRequest, resp *pb.MessageSendMessageResponse) {
	dumpServerToClientMessage(s2c, resp)

	if resp.Result == 0 && c.db != nil {
		uin, _ := strconv.ParseUint(s2c.Username, 10, 64)
		mr := &db.MessageRecord{
			Time:   resp.GetSendTime(),
			Seq:    req.GetMessageSeq(),
			Uid:    int64(req.GetMessageRand()) | 1<<56,
			PeerID: req.GetRoutingHead().GetGroup().GetCode(),
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
			mr.PeerID = req.GetRoutingHead().GetGroupTemp().GetUin()
			mr.UserID = req.GetRoutingHead().GetGroupTemp().GetToUin()
			mr.Type = 141
		}
		text, _ := mark.NewMarshaler(mr.PeerID, mr.UserID, mr.FromID).
			Marshal(req.GetMessageBody().GetRichText().GetElements())
		mr.Text = string(text)

		c.PrintMessageRecord(mr)
		if c.db != nil {
			err := c.dbInsertMessageRecord(uin, mr)
			if err != nil {
				log.Error().Err(err).Msg(">>> [db  ] dbInsertMessageRecord")
			}
		}
	}
}

func NewMessageSendMessageRequest(
	routingHead *pb.RoutingHead,
	contentHead *pb.ContentHead,
	messageBody *pb.MessageBody,
	seq int32,
	cookie []byte,
) *pb.MessageSendMessageRequest {
	return &pb.MessageSendMessageRequest{
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
	req *pb.MessageSendMessageRequest,
) (*pb.MessageSendMessageResponse, error) {
	if req.GetMessageSeq() == 0 {
		peerID := req.GetRoutingHead().GetGroup().GetCode()
		userID := req.GetRoutingHead().GetC2C().GetToUin()
		chatID := fmt.Sprintf("@%du%d", peerID, userID)
		req.MessageSeq = c.getNextMessageSeq(chatID)
	}
	if req.GetMessageRand() == 0 {
		req.MessageRand = rand.Int31()
	}
	if len(req.GetSyncCookie()) == 0 {
		req.SyncCookie = c.syncCookie
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
	resp := pb.MessageSendMessageResponse{}
	if err := proto.Unmarshal(s2c.Buffer, &resp); err != nil {
		return nil, err
	}
	// resp.Result
	//     0: success
	//     1: ???
	//    16: elements (notFriend)
	//   241: ???
	// -3902: marketFace (vip/svip)
	// -4902: marketFace magic (vip/svip)

	c.handleMessageSendMessageResponse(&s2c, req, &resp)
	return &resp, nil
}
