package rpc

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/pb"
)

type PushMessageNotifyRequest struct {
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
	VRIP           uint16       `jce:",13" json:",omitempty"`

	Unknown14 []byte                             `jce:",14" json:",omitempty"`
	Unknown15 *PushMessageNotifyRequestUnknown15 `jce:",15" json:",omitempty"`
	Unknown16 *PushMessageNotifyRequestUnknown15 `jce:",16" json:",omitempty"`
	Unknown17 *PushMessageNotifyRequestUnknown17 `jce:",17" json:",omitempty"`
}

type PushMessageNotifyRequestUnknown15 struct {
	Unknown0 uint64 `jce:",0" json:",omitempty"`
	Unknown1 uint64 `jce:",1" json:",omitempty"`
	Unknown2 uint64 `jce:",2" json:",omitempty"`
	Unknown3 uint64 `jce:",3" json:",omitempty"`
	Unknown4 uint64 `jce:",4" json:",omitempty"`
}

type PushMessageNotifyRequestUnknown17 struct {
	Unknown0 string `jce:",0" json:",omitempty"`
	Unknown1 string `jce:",1" json:",omitempty"`
	Unknown2 string `jce:",2" json:",omitempty"`
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

func (c *Client) handlePushMessageNotify(
	ctx context.Context,
	s2c *ServerToClientMessage,
) (*ClientToServerMessage, error) {
	msg := new(uni.Message)
	req := new(PushMessageNotifyRequest)
	if err := uni.Unmarshal(ctx, s2c.Buffer, msg, map[string]interface{}{
		"req_PushNotify": req,
	}); err != nil {
		return nil, err
	}
	jreq, _ := json.MarshalIndent(&req, "", "  ")
	log.Printf("rpc.PushMessageNotifyRequest\n%s", jreq)
	log.Printf(
		"handlePushMessageNotify, userActive:%d bindedUin:%d command:%s generalFlag:0x%08x pingFlag:0x%016x",
		req.UserActive, req.BindedUin, req.Command, req.GeneralFlag, req.PingFlag,
	)
	resp, err := c.MessageGetMessage(ctx, NewMessageGetMessageRequest(
		s2c.Username, 0x00000000, 0x00000001, 0x00000000,
	))
	if err != nil {
		return nil, err
	}
	dataList := []struct {
		PeerUin uint64
		FromUin uint64
		Data    []byte
	}{}
	for {
		// TODO: logic
		for _, uinPairMessage := range resp.GetUinPairMessages() {
			log.Printf("==> [sync] peer %d:", uinPairMessage.GetPeerUin())
			for _, msg := range uinPairMessage.GetMessages() {
				data, err := mark.Marshal(msg)
				if err != nil {
					return nil, err
				}
				log.Printf(
					"==> [sync] peer %d from %d to %d:\n%s",
					msg.GetMessageHead().GetGroupInfo().GetGroupCode(),
					msg.GetMessageHead().GetFromUin(),
					msg.GetMessageHead().GetToUin(),
					string(data),
				)
				if s2c.Username == strconv.FormatInt(int64(msg.GetMessageHead().GetFromUin()), 10) {
					c.setSyncSeq(msg.GetMessageHead().GetGroupInfo().GetGroupCode(), msg.GetMessageHead().GetMessageSeq())
				} else {
					dataList = append(dataList, struct {
						PeerUin uint64
						FromUin uint64
						Data    []byte
					}{
						PeerUin: msg.GetMessageHead().GetGroupInfo().GetGroupCode(),
						FromUin: msg.GetMessageHead().GetFromUin(),
						Data:    data,
					})
				}
				switch msg.GetMessageHead().GetMessageType() {
				case 0, 26, 64, 38, 48, 53, 61, 63, 78, 81, 103, 107, 110, 111, 114, 118:
					_, _ = c.MessageDeleteMessage(ctx, &MessageDeleteMessageRequest{
						FromUin:     msg.GetMessageHead().GetFromUin(),
						ToUin:       msg.GetMessageHead().GetToUin(),
						MessageType: msg.GetMessageHead().GetMessageType(),
						MessageSeq:  msg.GetMessageHead().GetMessageSeq(),
						MessageUid:  msg.GetMessageHead().GetMessageUid(),
						Username:    s2c.Username,
					})
				}
			}
		}
		if resp.GetSyncFlag() == 0x00000001 {
			log.Printf("==> [sync] syncing, 0x01")
			resp, err = c.MessageGetMessage(ctx, NewMessageGetMessageRequest(
				s2c.Username, resp.GetSyncFlag(), 0x00000001, 0x00000000,
			))
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	for _, data := range dataList {
		log.Printf(">=< [echo] %d %d %s", data.PeerUin, data.FromUin, string(data.Data))
	}
	n := len(dataList)
	if n != 0 {
		subMsg := new(pb.Message)
		if err := mark.Unmarshal(dataList[n-1].Data, subMsg); err != nil {
			return nil, err
		}
		_, _ = c.MessageSendMessage(ctx, &MessageSendMessageRequest{
			SendMessageRequest: pb.SendMessageRequest{
				RoutingHead: &pb.RoutingHead{C2C: &pb.C2C{Uin: dataList[n-1].FromUin}},
				ContentHead: subMsg.GetContentHead(),
				MessageBody: subMsg.GetMessageBody(),
				MessageSeq:  c.getNextSyncSeq(dataList[n-1].PeerUin),
				MessageRand: 0x00000000,
				SyncCookie:  c.syncCookie,
			},
			Username: s2c.Username,
		})
	}
	return nil, nil
}
