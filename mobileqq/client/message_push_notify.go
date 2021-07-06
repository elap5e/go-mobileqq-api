package client

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/log"
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

	MessageInfo    *MessageInfo `jce:",9" json:",omitempty"`
	MessageCtrlBuf string       `jce:",10" json:",omitempty"`
	ServerBuf      []byte       `jce:",11" json:",omitempty"`
	PingFlag       uint64       `jce:",12" json:",omitempty"`
	Svrip          int32        `jce:",13" json:",omitempty"`
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
	NickName        []string         `jce:",18" json:",omitempty"`
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
	c.dumpServerToClientMessage(s2c, &req)
	resp, err := c.MessageGetMessage(
		ctx, s2c.Username, NewMessageGetMessageRequest(
			0x00000000, c.syncCookie,
		),
	)
	if err != nil {
		return nil, err
	}
	// TODO: logic
	type Data struct {
		PeerUin uint64
		FromUin uint64
		Data    []byte
	}
	dataList := []Data{}
	infoList := []MessageDeleteInfo{}
	for {
		for _, uinPairMessage := range resp.GetUinPairMessages() {
			log.Info().
				Uint32("@peer", uinPairMessage.GetPeerUin()).
				Uint32("lastReadTime", uinPairMessage.GetLastReadTime()).
				Msgf("<-> [sync] %d message(s)", len(uinPairMessage.GetMessages()))
			for _, msg := range uinPairMessage.GetMessages() {
				data, err := c.marshalMessage(msg)
				if err != nil {
					return nil, err
				}
				peerUin := msg.GetMessageHead().GetGroupInfo().GetGroupCode()
				fromUin := msg.GetMessageHead().GetFromUin()
				if s2c.Username == strconv.FormatInt(int64(fromUin), 10) {
					c.setSyncSeq(peerUin, msg.GetMessageHead().GetMessageSeq())
				} else if msg.GetMessageHead().GetMessageType() == 166 {
					// add to data list
					dataList = append(dataList, Data{
						PeerUin: peerUin,
						FromUin: fromUin,
						Data:    data,
					})
				}

				// message processed
				switch msg.GetMessageHead().GetMessageType() {
				case 9, 10, 31, 79, 97, 120, 132, 133, 166, 167:
					switch msg.GetMessageHead().GetC2CCmd() {
					case 11, 175:
						infoList = append(infoList, MessageDeleteInfo{
							FromUin:     fromUin,
							MessageTime: uint64(msg.GetMessageHead().GetMessageTime()),
							MessageSeq:  uint16(msg.GetMessageHead().GetMessageSeq()),
						})
					}
				case 0, 26, 64, 38, 48, 53, 61, 63:
				case 78, 81, 103, 107, 110, 111, 114, 118:
					_, _ = c.MessageDeleteMessage(ctx, s2c.Username, NewMessageDeleteMessageRequest(
						&pb.MessageDeleteMessageRequest_MessageItem{
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
			c.dumpServerToClientMessage(s2c, &req)
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
	for i, data := range dataList {
		msg := pb.Message{}
		if err := mark.Unmarshal(data.Data, &msg); err != nil {
			return nil, err
		}
		seq := c.getNextSyncSeq(data.PeerUin)
		if i == len(dataList)-1 {
			log.Info().
				Str("@peer", fmt.Sprintf("%d:%s:%d", data.PeerUin, s2c.Username, data.FromUin)).
				Uint32("@seq", seq).
				Int64("@time", time.Now().Unix()).
				Str("mark", string(data.Data)).
				Msg("<-- [send] message")
			_, _ = c.MessageSendMessage(
				ctx, s2c.Username, NewMessageSendMessageRequest(
					&pb.RoutingHead{C2C: &pb.C2C{Uin: data.FromUin}},
					msg.GetContentHead(),
					msg.GetMessageBody(),
					seq,
					c.syncCookie,
				),
			)
		}
	}
	return NewOnlinePushMessageResponse(ctx, s2c.Username, infoList, req.Svrip, s2c.Seq)
}
