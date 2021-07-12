package client

import (
	"context"
	"fmt"
	"math/rand"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
)

func NewMessageSendMessageRequest(
	routingHead *pb.RoutingHead,
	contentHead *pb.ContentHead,
	messageBody *pb.MessageBody,
	seq uint32,
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
		chatID := fmt.Sprintf("@%d_%d", peerID, userID)
		req.MessageSeq = c.getNextMessageSeq(chatID)
	}
	if req.GetMessageRand() == 0 {
		req.MessageRand = rand.Uint32()
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

	dumpServerToClientMessage(&s2c, &resp)
	return &resp, nil
}
