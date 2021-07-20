package client

import (
	"context"
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
	FromUin         int64            `jce:",0" json:",omitempty"`
	MessageTime     int64            `jce:",1" json:",omitempty"`
	MessageType     int16            `jce:",2" json:",omitempty"`
	MessageSeq      int32            `jce:",3" json:",omitempty"`
	Message         string           `jce:",4" json:",omitempty"`
	RealMessageTime int64            `jce:",5" json:",omitempty"`
	MessageBytes    []byte           `jce:",6" json:",omitempty"`
	AppShareID      int64            `jce:",7" json:",omitempty"`
	MessageCookies  []byte           `jce:",8" json:",omitempty"`
	AppShareCookie  []byte           `jce:",9" json:",omitempty"`
	MessageUid      int64            `jce:",10" json:",omitempty"`
	LastChangeTime  int64            `jce:",11" json:",omitempty"`
	CPicInfo        []CPicInfo       `jce:",12" json:",omitempty"`
	ShareData       *ShareData       `jce:",13" json:",omitempty"`
	FromInstID      int64            `jce:",14" json:",omitempty"`
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

	uin, _ := strconv.ParseInt(s2c.Username, 10, 64)
	resp, err := c.MessageGetMessage(
		ctx, s2c.Username, NewMessageGetMessageRequest(
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
	infos := []MessageDeleteInfo{}

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
						infos = append(infos, MessageDeleteInfo{
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
			dumpServerToClientMessage(s2c, &req)
			resp, err := c.MessageGetMessage(
				ctx, s2c.Username, NewMessageGetMessageRequest(
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

	return NewOnlinePushMessageResponse(ctx, s2c.Username, infos, req.ServerIP, int32(s2c.Seq))
}
