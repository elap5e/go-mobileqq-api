package client

import (
	"context"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
	"github.com/elap5e/go-mobileqq-api/util"
)

type MessagePushNotifyRequest struct {
	Uin         int64  `jce:",0" json:"uin,omitempty"`
	Type        uint8  `jce:",1" json:"type,omitempty"`
	Service     string `jce:",2" json:"service,omitempty"`
	Cmd         string `jce:",3" json:"cmd,omitempty"`
	Cookie      []byte `jce:",4" json:"cookie,omitempty"`
	MessageType uint16 `jce:",5" json:"message_type,omitempty"`
	UserActive  uint32 `jce:",6" json:"user_active,omitempty"`
	GeneralFlag uint32 `jce:",7" json:"general_flag,omitempty"`
	BindedUin   int64  `jce:",8" json:"binded_uin,omitempty"`

	Message       *Message     `jce:",9" json:"message,omitempty"`
	ControlBuffer string       `jce:",10" json:"control_buffer,omitempty"`
	ServerBuffer  []byte       `jce:",11" json:"server_buffer,omitempty"`
	PingFlag      uint64       `jce:",12" json:"ping_flag,omitempty"`
	ServerIP      Uint32IPType `jce:",13" json:"server_ip,omitempty"`
}

func (c *Client) handleMessagePushNotify(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error) {
	msg := uni.Message{}
	req := MessagePushNotifyRequest{}
	if err := uni.Unmarshal(ctx, s2c.Buffer, &msg, map[string]interface{}{
		"req_PushNotify": &req,
	}); err != nil {
		return nil, err
	}
	util.DumpServerToClientMessage(s2c, &req)

	uin, _ := strconv.ParseInt(s2c.Username, 10, 64)
	resp, err := c.MessageGetMessage(
		ctx, uin, NewMessageGetMessageRequest(
			0x00000000, c.syncCookie[uin],
		),
	)
	if err != nil {
		return nil, err
	}

	type Data struct {
		PeerID int64
		UserID int64
		FromID int64
		Text   []byte
	}
	dataList := []Data{}
	items := []MessageDelete{}

	for {
		for _, uinPairMessage := range resp.GetUinPairMessages() {
			for _, msg := range uinPairMessage.GetItems() {
				switch msg.GetMessageHead().GetMessageType() {
				case 9, 10, 31, 79, 97, 120, 132, 133, 141, 166, 167:
					switch msg.GetMessageHead().GetC2CCmd() {
					case 11, 175:
						if uin != msg.GetMessageHead().GetFromUin() {
							peerID := msg.GetMessageHead().GetC2CTempMessageHead().GetGroupCode()
							userID := uinPairMessage.GetPeerUin()
							fromID := msg.GetMessageHead().GetFromUin()
							text, err := mark.NewEncoder(peerID, userID, fromID).
								Encode(msg.GetMessageBody().GetRichText().GetElements())
							if err != nil {
								return nil, err
							}
							dataList = append(dataList, Data{
								PeerID: msg.GetMessageHead().GetC2CTempMessageHead().GetGroupCode(),
								UserID: uinPairMessage.GetPeerUin(),
								FromID: uin,
								Text:   text,
							})
						}
						items = append(items, MessageDelete{
							FromUin:     msg.GetMessageHead().GetFromUin(),
							MessageTime: msg.GetMessageHead().GetMessageTime(),
							MessageSeq:  msg.GetMessageHead().GetMessageSeq(),
						})
					}
				case 78, 81, 103, 107, 110, 111, 114, 118:
					_, _ = c.MessageDeleteMessage(ctx, s2c.Username, NewMessageDeleteMessageRequest(
						&pb.MessageService_DeleteRequest_MessageItem{
							FromUin:     msg.GetMessageHead().GetFromUin(),
							ToUin:       msg.GetMessageHead().GetToUin(),
							MessageType: msg.GetMessageHead().GetMessageType(),
							MessageSeq:  msg.GetMessageHead().GetMessageSeq(),
							MessageUid:  msg.GetMessageHead().GetMessageUid(),
						},
					))
				}
			}
		}
		if resp.GetSyncFlag() == 0x00000001 {
			util.DumpServerToClientMessage(s2c, &req)
			resp, err := c.MessageGetMessage(
				ctx, uin, NewMessageGetMessageRequest(
					resp.GetSyncFlag(), c.syncCookie[uin],
				),
			)
			if err != nil {
				return nil, err
			}
			c.syncCookie[uin] = resp.GetSyncCookie()
		} else {
			break
		}
	}

	// echo message
	if l := len(dataList); l > 0 {
		item := dataList[l-1]

		elems, err := mark.NewDecoder(item.PeerID, item.UserID, item.FromID).
			Decode(item.Text)
		if err != nil {
			return nil, err
		}
		msg := pb.MessageCommon_Message{
			MessageBody: &pb.IMMessageBody_MessageBody{
				RichText: &pb.IMMessageBody_RichText{
					Elements: elems,
				},
			},
		}
		if _, err := c.MessageSendMessage(
			ctx, s2c.Username, NewMessageSendMessageRequest(
				c.GetRoutingHead(item.PeerID, item.UserID),
				msg.GetContentHead(),
				msg.GetMessageBody(),
				0,
				c.syncCookie[uin],
			),
		); err != nil {
			return nil, err
		}
	}

	return NewOnlinePushMessageResponse(ctx, s2c.Username, items, req.ServerIP, int32(s2c.Seq))
}
