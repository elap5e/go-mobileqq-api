package client

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
)

func (c *Client) handleMessageGetMessageResponse(s2c *codec.ServerToClientMessage, resp *pb.MessageGetMessageResponse) {
	dumpServerToClientMessage(s2c, resp)

	for _, uinPairMessage := range resp.GetUinPairMessages() {
		syncUinPairMessage(uinPairMessage)

		for _, msg := range uinPairMessage.GetMessages() {
			skip := true
			switch msg.GetMessageHead().GetMessageType() {
			case 9, 10, 31, 79, 97, 120, 132, 133, 141, 166, 167:
				skip = msg.GetMessageHead().GetC2CCmd() == 0
			case 208:
				skip = msg.GetMessageHead().GetC2CCmd() == 0
			case 529:
				skip = msg.GetMessageHead().GetC2CCmd() == 0
			case 43, 82:
				skip = msg.GetMessageHead().GetGroupInfo() == nil
			case 42, 83:
				skip = msg.GetMessageHead().GetDiscussInfo() == nil
			}
			if !skip {
				mr := &db.MessageRecord{
					Time:   msg.GetMessageHead().GetMessageTime(),
					Seq:    msg.GetMessageHead().GetMessageSeq(),
					Uid:    int64(msg.GetMessageBody().GetRichText().GetAttribute().GetRandom()) | 1<<56,
					PeerID: msg.GetMessageHead().GetC2CTempMessageHead().GetGroupCode(),
					UserID: uinPairMessage.GetPeerUin(),
					FromID: msg.GetMessageHead().GetFromUin(),
					Text:   "",
					Type:   msg.GetMessageHead().GetMessageType(),
				}
				if msg.GetMessageHead().GetC2CCmd() == 0 {
					mr.PeerID = uinPairMessage.GetPeerUin()
					mr.UserID = 0
					if mr.Type == 82 {
						c.setMessageSeq(fmt.Sprintf("@%du%d", mr.PeerID, mr.UserID), mr.Seq)
					}
				}
				text, _ := mark.NewMarshaler(mr.PeerID, mr.UserID, mr.FromID).
					Marshal(msg.GetMessageBody().GetRichText().GetElements())
				mr.Text = string(text)

				c.PrintMessageRecord(mr)
				if c.db != nil {
					uin, _ := strconv.ParseUint(s2c.Username, 10, 64)
					err := c.dbInsertMessageRecord(uin, mr)
					if err != nil {
						log.Error().Err(err).Msg(">>> [db  ] dbInsertMessageRecord")
					}
				}
			}
		}
	}
}

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
	if len(req.SyncCookie) == 0 {
		req.SyncCookie = c.syncCookie
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
	err = c.rpc.Call(ServiceMethodMessageGetMessage, &c2s, &s2c)
	if err != nil {
		return nil, err
	}
	resp := pb.MessageGetMessageResponse{}
	if err := proto.Unmarshal(s2c.Buffer, &resp); err != nil {
		return nil, err
	}

	c.syncCookie = resp.GetSyncCookie()

	c.handleMessageGetMessageResponse(&s2c, &resp)
	return &resp, nil
}
