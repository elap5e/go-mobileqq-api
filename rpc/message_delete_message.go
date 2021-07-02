package rpc

import (
	"context"
	"encoding/json"
	"log"

	"github.com/elap5e/go-mobileqq-api/pb"
	"google.golang.org/protobuf/proto"
)

type MessageDeleteMessageRequest struct {
	FromUin     uint64
	ToUin       uint64
	MessageType uint32
	MessageSeq  uint32
	MessageUid  uint64

	Username string
}

func (c *Client) MessageDeleteMessage(
	ctx context.Context,
	req *MessageDeleteMessageRequest,
) (*pb.DeleteMessageResponse, error) {
	buf, err := proto.Marshal(&pb.DeleteMessageRequest{
		MessageItems: []*pb.DeleteMessageRequest_MessageItem{{
			FromUin:     req.FromUin,
			ToUin:       req.ToUin,
			MessageType: req.MessageType,
			MessageSeq:  req.MessageSeq,
			MessageUid:  req.MessageUid,
		}},
	})
	if err != nil {
		return nil, err
	}
	s2c := new(ServerToClientMessage)
	if err := c.Call(ServiceMethodMessageDeleteMessage, &ClientToServerMessage{
		Username: req.Username,
		Seq:      c.getNextSeq(),
		Buffer:   buf,
		Simple:   true,
	}, s2c); err != nil {
		return nil, err
	}
	resp := pb.DeleteMessageResponse{}
	if err := proto.Unmarshal(s2c.Buffer, &resp); err != nil {
		return nil, err
	}
	jresp, _ := json.MarshalIndent(&resp, "", "  ")
	log.Printf("pb.DeleteMessageResponse\n%s", jresp)
	return &resp, nil
}
