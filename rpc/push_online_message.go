package rpc

import (
	"context"
	"encoding/json"
	"log"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/pb"
	"google.golang.org/protobuf/proto"
)

func (c *Client) handlePushOnlineMessage(
	ctx context.Context,
	s2c *ServerToClientMessage,
) (*ClientToServerMessage, error) {
	msg := pb.OnlinePushMessage{}
	if err := proto.Unmarshal(s2c.Buffer, &msg); err != nil {
		return nil, err
	}
	jmsg, _ := json.MarshalIndent(&msg, "", "  ")
	log.Printf("pb.OnlinePushMessage\n%s", jmsg)
	data, err := mark.Marshal(msg.GetMessage())
	if err != nil {
		return nil, err
	}
	log.Printf(
		"==> [sync] peer %d from %d to %d:\n%s",
		msg.GetMessage().GetMessageHead().GetGroupInfo().GetGroupCode(),
		msg.GetMessage().GetMessageHead().GetFromUin(),
		msg.GetMessage().GetMessageHead().GetToUin(),
		string(data),
	)
	return nil, nil
}
