package rpc

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/pb"
)

func NewMessageGetMessageRequest(
	flag uint32,
	cookie []byte,
) *pb.MessageGetMessageRequest {
	return &pb.MessageGetMessageRequest{
		SyncFlag:            flag,
		SyncCookie:          cookie,
		RambleFlag:          0x00000000,
		LatestRambleNumber:  0x00000014,
		OtherRambleNumber:   0x00000003,
		OnlineSyncFlag:      0x00000001, // fix
		ContextFlag:         0x00000001,
		WhisperSessionId:    0x00000000,
		RequestType:         0x00000000, // fix
		PublicAccountCookie: nil,
		ControlBuffer:       nil,
		ServerBuffer:        nil,
	}
}

func (c *Client) MessageGetMessage(
	ctx context.Context,
	username string,
	req *pb.MessageGetMessageRequest,
) (*pb.MessageGetMessageResponse, error) {
	buf, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	s2c := ServerToClientMessage{}
	if err := c.Call(ServiceMethodMessageGetMessage, &ClientToServerMessage{
		Username: username,
		Seq:      c.getNextSeq(),
		Buffer:   buf,
		Simple:   true,
	}, &s2c); err != nil {
		return nil, err
	}
	resp := pb.MessageGetMessageResponse{}
	if err := proto.Unmarshal(s2c.Buffer, &resp); err != nil {
		return nil, err
	}
	c.dumpServerToClientMessage(&s2c, &resp)
	return &resp, nil
}