package rpc

import (
	"context"
	"encoding/binary"

	"github.com/elap5e/go-mobileqq-api/pb"
	"google.golang.org/protobuf/proto"
)

func (c *Client) MessageGetMessage(
	ctx context.Context,
	s2c *ServerToClientMessage,
) (*ClientToServerMessage, error) {
	buf, err := proto.Marshal(&pb.GetMessageRequest{
		SyncFlag:           0x00000001,
		SyncCookie:         []byte{},
		RambleFlag:         0x00000000,
		LatestRambleNumber: 0x00000014,
		OtherRambleNumber:  0x00000003,
		OnlineSyncFlag:     0x00000000,
		ContextFlag:        0x00000001,
		WhisperSessionId:   0x00000000,
	})
	if err != nil {
		return nil, err
	}
	data := append(make([]byte, 4), buf...)
	binary.BigEndian.PutUint32(data[0:], uint32(len(data)))
	if err := c.Call(ServiceMethodMessageGetMessage, &ClientToServerMessage{
		Username: s2c.Username,
		Seq:      s2c.Seq,
		Buffer:   data,
		Simple:   true,
	}, s2c); err != nil {
		return nil, err
	}
	return nil, nil
}
