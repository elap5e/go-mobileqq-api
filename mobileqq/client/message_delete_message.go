package client

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
)

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
	if len(req.GetMessageItems()) == 0 {
		return nil, fmt.Errorf("zero length")
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
	err = c.rpc.Call(ServiceMethodMessageDeleteMessage, &c2s, &s2c)
	if err != nil {
		return nil, err
	}
	resp := pb.MessageDeleteMessageResponse{}
	if err := proto.Unmarshal(s2c.Buffer, &resp); err != nil {
		return nil, err
	}

	c.dumpServerToClientMessage(&s2c, &resp)
	return &resp, nil
}
