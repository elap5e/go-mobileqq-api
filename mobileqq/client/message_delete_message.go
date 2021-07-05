package client

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
)

// type MessageDeleteMessageRequest struct {
// 	FromUin     uint64
// 	ToUin       uint64
// 	MessageType uint32
// 	MessageSeq  uint32
// 	MessageUid  uint64

// 	Username string
// }

func NewMessageDeleteMessageRequest(
	items ...*pb.MessageDeleteMessageRequest_MessageItem,
) *pb.MessageDeleteMessageRequest {
	return &pb.MessageDeleteMessageRequest{
		MessageItems: items,
	}
}

func (c *Client) MessageDeleteMessage(
	ctx context.Context,
	username string,
	req *pb.MessageDeleteMessageRequest,
) (*pb.MessageDeleteMessageResponse, error) {
	buf, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	s2c := codec.ServerToClientMessage{}
	if err := c.rpc.Call(ServiceMethodMessageDeleteMessage, &codec.ClientToServerMessage{
		Username: username,
		Buffer:   buf,
		Simple:   true,
	}, &s2c); err != nil {
		return nil, err
	}
	resp := pb.MessageDeleteMessageResponse{}
	if err := proto.Unmarshal(s2c.Buffer, &resp); err != nil {
		return nil, err
	}
	c.dumpServerToClientMessage(&s2c, &resp)
	return &resp, nil
}
