package rpc

import (
	"context"
	"encoding/json"
	"log"

	"github.com/elap5e/go-mobileqq-api/pb"
	"google.golang.org/protobuf/proto"
)

type MessageGetMessageRequest struct {
	pb.GetMessageRequest

	Username string
}

func NewMessageGetMessageRequest(
	username string,
	flag, onlineFlag uint32,
	reqType uint32,
) *MessageGetMessageRequest {
	return &MessageGetMessageRequest{
		GetMessageRequest: pb.GetMessageRequest{
			SyncFlag:            flag,
			SyncCookie:          []byte{},
			RambleFlag:          0x00000000,
			LatestRambleNumber:  0x00000014,
			OtherRambleNumber:   0x00000003,
			OnlineSyncFlag:      onlineFlag,
			ContextFlag:         0x00000001,
			WhisperSessionId:    0x00000000,
			RequestType:         reqType,
			PublicAccountCookie: nil,
			ControlBuffer:       nil,
			ServerBuffer:        nil,
		},
		Username: username,
	}
}

func (c *Client) MessageGetMessage(
	ctx context.Context,
	req *MessageGetMessageRequest,
) (*pb.GetMessageResponse, error) {
	req.SyncCookie = c.syncCookie
	buf, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	s2c := new(ServerToClientMessage)
	if err := c.Call(ServiceMethodMessageGetMessage, &ClientToServerMessage{
		Username: req.Username,
		Seq:      c.getNextSeq(),
		Buffer:   buf,
		Simple:   true,
	}, s2c); err != nil {
		return nil, err
	}
	resp := pb.GetMessageResponse{}
	if err := proto.Unmarshal(s2c.Buffer, &resp); err != nil {
		return nil, err
	}
	c.syncCookie = resp.GetSyncCookie()
	jresp, _ := json.MarshalIndent(&resp, "", "  ")
	log.Printf("pb.GetMessageResponse\n%s", jresp)
	return &resp, nil
}
