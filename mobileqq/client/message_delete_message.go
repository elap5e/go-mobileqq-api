package client

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
)

func NewMessageDeleteMessageRequest(
	items ...*pb.MessageService_DeleteRequest_MessageItem,
) *pb.MessageService_DeleteRequest {
	return &pb.MessageService_DeleteRequest{
		Items: items,
	}
}

func (c *Client) MessageDeleteMessage(
	ctx context.Context,
	username string,
	req *pb.MessageService_DeleteRequest,
) (*pb.MessageService_DeleteResponse, error) {
	if len(req.GetItems()) == 0 {
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
	resp := pb.MessageService_DeleteResponse{}
	if err := proto.Unmarshal(s2c.Buffer, &resp); err != nil {
		return nil, err
	}

	dumpServerToClientMessage(&s2c, &resp)
	return &resp, nil
}
