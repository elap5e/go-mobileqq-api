package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
)

type MessagePushNotifyRequest struct {
	Uin         uint64 `jce:",0" json:",omitempty"`
	Type        uint8  `jce:",1" json:",omitempty"`
	Service     string `jce:",2" json:",omitempty"`
	Command     string `jce:",3" json:",omitempty"`
	Cookie      []byte `jce:",4" json:",omitempty"`
	MessageType uint16 `jce:",5" json:",omitempty"`
	UserActive  uint32 `jce:",6" json:",omitempty"`
	GeneralFlag uint32 `jce:",7" json:",omitempty"`
	BindedUin   uint64 `jce:",8" json:",omitempty"`

	MessageInfo       *MessageInfo `jce:",9" json:",omitempty"`
	MessageCtrlBuffer string       `jce:",10" json:",omitempty"`
	ServerBuffer      []byte       `jce:",11" json:",omitempty"`
	PingFlag          uint64       `jce:",12" json:",omitempty"`
	ServerIP          uint32       `jce:",13" json:",omitempty"`
}

type MessageInfo struct {
	FromUin         uint64           `jce:",0" json:",omitempty"`
	MessageTime     uint64           `jce:",1" json:",omitempty"`
	MessageType     uint16           `jce:",2" json:",omitempty"`
	MessageSeq      uint16           `jce:",3" json:",omitempty"`
	Message         string           `jce:",4" json:",omitempty"`
	RealMessageTime uint64           `jce:",5" json:",omitempty"`
	MessageBytes    []byte           `jce:",6" json:",omitempty"`
	AppShareID      uint64           `jce:",7" json:",omitempty"`
	MessageCookies  []byte           `jce:",8" json:",omitempty"`
	AppShareCookie  []byte           `jce:",9" json:",omitempty"`
	MessageUid      uint64           `jce:",10" json:",omitempty"`
	LastChangeTime  uint64           `jce:",11" json:",omitempty"`
	CPicInfo        []CPicInfo       `jce:",12" json:",omitempty"`
	ShareData       *ShareData       `jce:",13" json:",omitempty"`
	FromInstID      uint64           `jce:",14" json:",omitempty"`
	RemarkOfSender  []byte           `jce:",15" json:",omitempty"`
	FromMobile      string           `jce:",16" json:",omitempty"`
	FromName        string           `jce:",17" json:",omitempty"`
	Nickname        []string         `jce:",18" json:",omitempty"`
	TempMessageHead *TempMessageHead `jce:",19" json:",omitempty"`
}

type CPicInfo struct {
	Path []byte `jce:",0" json:",omitempty"`
	Host []byte `jce:",1" json:",omitempty"`
}

type ShareData struct {
	Pkgname     string `jce:",0" json:",omitempty"`
	Messagetail string `jce:",1" json:",omitempty"`
	PicURL      string `jce:",2" json:",omitempty"`
	URL         string `jce:",3" json:",omitempty"`
}

type TempMessageHead struct {
	C2CType     uint32 `jce:",0" json:",omitempty"`
	ServiceType uint32 `jce:",1" json:",omitempty"`
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
	dumpServerToClientMessage(s2c, &req)

	resp, err := c.MessageGetMessage(
		ctx, s2c.Username, NewMessageGetMessageRequest(
			0x00000000, c.syncCookie,
		),
	)
	if err != nil {
		return nil, err
	}
	uin, _ := strconv.ParseInt(s2c.Username, 10, 64)

	type Data struct {
		PeerID int64
		UserID int64
		FromID int64
		Text   []byte
	}
	dataList := []Data{}
	infos := []MessageDeleteInfo{}

	for {
		for _, uinPairMessage := range resp.GetUinPairMessages() {
			for _, msg := range uinPairMessage.GetMessages() {
				switch msg.GetMessageHead().GetMessageType() {
				case 9, 10, 31, 79, 97, 120, 132, 133, 141, 166, 167:
					switch msg.GetMessageHead().GetC2CCmd() {
					case 11, 175:
						if uin != msg.GetMessageHead().GetFromUin() {
							text, err := mark.Marshal(msg)
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
						infos = append(infos, MessageDeleteInfo{
							FromUin:     msg.GetMessageHead().GetFromUin(),
							MessageTime: msg.GetMessageHead().GetMessageTime(),
							MessageSeq:  msg.GetMessageHead().GetMessageSeq(),
						})
					}
				case 78, 81, 103, 107, 110, 111, 114, 118:
					_, _ = c.MessageDeleteMessage(ctx, s2c.Username, NewMessageDeleteMessageRequest(
						&pb.MessageDeleteMessageRequest_Item{
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
			dumpServerToClientMessage(s2c, &req)
			resp, err := c.MessageGetMessage(
				ctx, s2c.Username, NewMessageGetMessageRequest(
					resp.GetSyncFlag(), c.syncCookie,
				),
			)
			if err != nil {
				return nil, err
			}
			c.syncCookie = resp.GetSyncCookie()
		} else {
			break
		}
	}

	// echo message
	if l := len(dataList); l > 0 {
		item := dataList[l-1]

		routingHead := &pb.RoutingHead{}
		if item.PeerID == 0 {
			routingHead = &pb.RoutingHead{C2C: &pb.C2C{ToUin: item.UserID}}
		} else if item.UserID == 0 {
			routingHead = &pb.RoutingHead{Group: &pb.Group{Code: item.PeerID}}
		} else {
			routingHead = &pb.RoutingHead{
				GroupTemp: &pb.GroupTemp{Uin: item.PeerID, ToUin: item.UserID},
			}
		}
		chatID := fmt.Sprintf("@%du%d", item.PeerID, item.UserID)
		seq := c.getNextMessageSeq(chatID)

		msg := pb.Message{}
		if err := mark.Unmarshal(item.Text, &msg); err != nil {
			return nil, err
		}
		if _, err := c.MessageSendMessage(
			ctx, s2c.Username, NewMessageSendMessageRequest(
				routingHead,
				msg.GetContentHead(),
				msg.GetMessageBody(),
				seq,
				c.syncCookie,
			),
		); err != nil {
			return nil, err
		}
	}

	return NewOnlinePushMessageResponse(ctx, s2c.Username, infos, req.ServerIP, int32(s2c.Seq))
}
