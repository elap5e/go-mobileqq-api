package client

import (
	"context"
	"encoding/hex"
	"strconv"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
)

func (c *Client) handleMessageGetMessageResponse(s2c *codec.ServerToClientMessage, resp *pb.MessageService_GetResponse) {
	uin, _ := strconv.ParseUint(s2c.Username, 10, 64)
	dumpServerToClientMessage(s2c, resp)

	for _, uinPairMessage := range resp.GetUinPairMessages() {
		syncUinPairMessage(uinPairMessage)

		for _, msg := range uinPairMessage.GetItems() {
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
			case 0x0210:
				body, err := c.decodeMessageType0210Pb(uin, msg.GetMessageBody().GetContent())
				if err != nil {
					log.Error().Err(err).Msg(">>x [0210] failed to decode")
				} else if body != nil {
					dumpServerToClientMessage(s2c, &body)
				}
			case 0x02DC:
				body, err := c.decodeMessageType02DC(uin, msg.GetMessageBody().GetContent())
				if err != nil {
					log.Error().Err(err).Msg(">>x [02dC] failed to decode")
				} else if body != nil {
					dumpServerToClientMessage(s2c, &body)
				}
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
						c.setMessageSeq(mr.PeerID, mr.UserID, int64(uin), mr.Seq)
					}
				}
				text, _ := mark.NewEncoder(mr.PeerID, mr.UserID, mr.FromID).
					Encode(msg.GetMessageBody().GetRichText().GetElements())
				mr.Text = string(text)

				c.PrintMessageRecord(mr)
				if c.db != nil {
					err := c.dbInsertMessageRecord(uin, mr)
					if err != nil {
						log.Error().Err(err).Msg(">>x [db  ] dbInsertMessageRecord")
					}
				}
			}
		}
	}
}

func NewMessageGetMessageRequest(
	flag uint32,
	cookie []byte,
) *pb.MessageService_GetRequest {
	return &pb.MessageService_GetRequest{
		SyncFlag:             flag,
		SyncCookie:           cookie,
		RambleFlag:           0x00000000,
		LatestRambleNumber:   0x00000014,
		OtherRambleNumber:    0x00000003,
		OnlineSyncFlag:       0x00000001, // fix
		ContextFlag:          0x00000001,
		WhisperSessionId:     0x00000000,
		MessageRequestType:   0x00000000, // fix
		PublicAccountCookie:  nil,
		MessageControlBuffer: nil,
		ServerBuffer:         nil,
	}
}

func (c *Client) MessageGetMessage(
	ctx context.Context,
	username string,
	req *pb.MessageService_GetRequest,
) (*pb.MessageService_GetResponse, error) {
	uin, _ := strconv.ParseInt(username, 10, 64)
	if len(req.SyncCookie) == 0 {
		req.SyncCookie = c.syncCookie[uin]
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
	resp := pb.MessageService_GetResponse{}
	if err := proto.Unmarshal(s2c.Buffer, &resp); err != nil {
		log.Debug().Msg(">>> [dump]\n" + hex.Dump(s2c.Buffer))
		return nil, err
	}

	if c.db != nil {
		uin, _ := strconv.ParseInt(s2c.Username, 10, 64)
		if err := c.dbUpdateAccount(&db.Account{
			Uin:        uin,
			SyncCookie: resp.GetSyncCookie(),
		}); err != nil {
			log.Fatal().Err(err).
				Msg("failed to operate database")
		}
	}
	c.syncCookie[uin] = resp.GetSyncCookie()

	c.handleMessageGetMessageResponse(&s2c, &resp)
	return &resp, nil
}
