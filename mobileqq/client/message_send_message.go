package client

import (
	"context"

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
		MessageRand: 0x00000000,
		SyncCookie:  cookie,
	}
}

func (c *Client) MessageSendMessage(
	ctx context.Context,
	username string,
	req *pb.MessageSendMessageRequest,
) (*pb.MessageSendMessageResponse, error) {
	if req.GetMessageSeq() == 0 {
		peerUin := req.GetRoutingHead().GetGroup().GetCode()
		req.MessageSeq = c.getNextSyncSeq(peerUin)
	}
	if len(req.GetSyncCookie()) == 0 {
		req.SyncCookie = c.syncCookie
	}
	buf, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	s2c := codec.ServerToClientMessage{}
	if err := c.rpc.Call(ServiceMethodMessageSendMessage, &codec.ClientToServerMessage{
		Username: username,
		Buffer:   buf,
		Simple:   true,
	}, &s2c); err != nil {
		return nil, err
	}
	resp := pb.MessageSendMessageResponse{}
	if err := proto.Unmarshal(s2c.Buffer, &resp); err != nil {
		return nil, err
	}
	c.dumpServerToClientMessage(&s2c, &resp)
	return &resp, nil
}
