package rpc

import (
	"context"
	"encoding/json"
	"log"

	"github.com/elap5e/go-mobileqq-api/pb"
	"google.golang.org/protobuf/proto"
)

type MessageSendMessageRequest struct {
	pb.SendMessageRequest

	Username string
}

func (c *Client) MessageSendMessage(
	ctx context.Context,
	req *MessageSendMessageRequest,
) (*pb.SendMessageResponse, error) {
	buf, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	s2c := new(ServerToClientMessage)
	if err := c.Call(ServiceMethodMessageSendMessage, &ClientToServerMessage{
		Username: req.Username,
		Seq:      c.getNextSeq(),
		Buffer:   buf,
		Simple:   true,
	}, s2c); err != nil {
		return nil, err
	}
	resp := pb.SendMessageResponse{}
	if err := proto.Unmarshal(s2c.Buffer, &resp); err != nil {
		return nil, err
	}
	jresp, _ := json.MarshalIndent(&resp, "", "  ")
	log.Printf("pb.SendMessageResponse\n%s", jresp)
	return &resp, nil
}
