package client

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
	"github.com/elap5e/go-mobileqq-api/pb"
)

type Message struct {
	FromUin         int64            `jce:",0" json:"from_uin,omitempty"`
	MessageTime     int64            `jce:",1" json:"message_time,omitempty"`
	MessageType     int16            `jce:",2" json:"message_type,omitempty"`
	MessageSeq      int32            `jce:",3" json:"message_seq,omitempty"`
	Message         string           `jce:",4" json:"message,omitempty"`
	RealMessageTime int64            `jce:",5" json:"real_message_time,omitempty"`
	MessageBytes    []byte           `jce:",6" json:"message_bytes,omitempty"`
	AppShareID      int64            `jce:",7" json:"app_share_id,omitempty"`
	MessageCookie   []byte           `jce:",8" json:"message_cookie,omitempty"`
	AppShareCookie  []byte           `jce:",9" json:"app_share_cookie,omitempty"`
	MessageUid      int64            `jce:",10" json:"message_uid,omitempty"`
	LastChangeTime  int64            `jce:",11" json:"last_change_time,omitempty"`
	Pictures        []Picture        `jce:",12" json:"pictures,omitempty"`
	ShareData       *ShareData       `jce:",13" json:"share_data,omitempty"`
	FromInstanceID  int64            `jce:",14" json:"from_instance_id,omitempty"`
	FromRemark      []byte           `jce:",15" json:"from_remark,omitempty"`
	FromMobile      string           `jce:",16" json:"from_mobile,omitempty"`
	FromName        string           `jce:",17" json:"from_name,omitempty"`
	FromNick        []string         `jce:",18" json:"from_nick,omitempty"`
	TempMessageHead *TempMessageHead `jce:",19" json:"temp_message_head,omitempty"`
}

type MessageDelete struct {
	FromUin       int64  `jce:",0" json:"from_uin,omitempty"`
	MessageTime   int64  `jce:",1" json:"message_time,omitempty"`
	MessageSeq    int32  `jce:",2" json:"message_seq,omitempty"`
	MessageCookie []byte `jce:",3" json:"message_cookie,omitempty"`
}

type MessageReadedC2C struct {
	Uin          int64 `jce:",0" json:"uin,omitempty"`
	LastReadTime int64 `jce:",1" json:"last_read_time,omitempty"`
}

type MessageReadedGroup struct {
	Uin        int64 `jce:",0" json:"uin,omitempty"`
	Type       int64 `jce:",1" json:"type,omitempty"`
	MemberSeq  int32 `jce:",2" json:"member_seq,omitempty"`
	MessageSeq int32 `jce:",3" json:"message_seq,omitempty"`
}

type MessageReadedDiscuss struct {
	Uin        int64 `jce:",0" json:"uin,omitempty"`
	Type       int64 `jce:",1" json:"type,omitempty"`
	MemberSeq  int32 `jce:",2" json:"member_seq,omitempty"`
	MessageSeq int32 `jce:",3" json:"message_seq,omitempty"`
}

type Picture struct {
	Path []byte `jce:",0" json:"path,omitempty"`
	Host []byte `jce:",1" json:"host,omitempty"`
}

type ShareData struct {
	PackageName string `jce:",0" json:"package_name,omitempty"`
	MessageTail string `jce:",1" json:"message_tail,omitempty"`
	PictureURL  string `jce:",2" json:"picture_url,omitempty"`
	URL         string `jce:",3" json:"url,omitempty"`
}

type TempMessageHead struct {
	C2CType     int32 `jce:",0" json:"c2c_type,omitempty"`
	ServiceType int32 `jce:",1" json:"service_type,omitempty"`
}

type UinPairMessage struct {
	LastReadTime     int64      `jce:",1" json:"last_read_time,omitempty"`
	PeerUin          int64      `jce:",2" json:"peer_uin,omitempty"`
	MessageCompleted int64      `jce:",3" json:"message_completed,omitempty"`
	Messages         []*Message `jce:",4" json:"messages,omitempty"`
}

func (c *Client) syncUinPairMessage(uinPairMessage *pb.MessageCommon_UinPairMessage) {
	log.Info().
		Int64("@uin", uinPairMessage.GetPeerUin()).
		Int64("last_read_time", uinPairMessage.GetLastReadTime()).
		Msgf("<-> [sync] %d message(s)", len(uinPairMessage.GetItems()))
	for _, item := range uinPairMessage.GetItems() {
		_ = item
	}
}

func (c *Client) GetRoutingHead(peerID, userID int64) *pb.MessageService_RoutingHead {
	routingHead := pb.MessageService_RoutingHead{}
	if peerID == 0 {
		routingHead.C2C = &pb.MessageService_C2C{ToUin: userID}
	} else if userID == 0 {
		routingHead.Group = &pb.MessageService_Group{GroupCode: peerID}
	} else {
		routingHead.GroupTemp = &pb.MessageService_GroupTemp{
			GroupUin: peerID, ToUin: userID,
		}
	}
	return &routingHead
}

func (c *Client) PrintMessageRecord(mr *db.MessageRecord) {
	peerName := strconv.FormatInt(mr.PeerID, 10)
	userName := strconv.FormatInt(mr.UserID, 10)
	fromName := strconv.FormatInt(mr.FromID, 10)
	if mr.PeerID == 0 {
		if contact, ok := c.contacts[mr.UserID]; ok {
			userName = contact.Remark
		}
		if contact, ok := c.contacts[mr.FromID]; ok {
			fromName = contact.Remark
		}
	} else {
		if channel, ok := c.channels[mr.PeerID]; ok {
			peerName = channel.Name
		}
		if cmember, ok := c.cmembers[mr.PeerID][mr.UserID]; ok {
			userName = cmember.Remark
			if userName == "" {
				userName = cmember.Nick
			}
		}
		if cmember, ok := c.cmembers[mr.PeerID][mr.FromID]; ok {
			fromName = cmember.Remark
			if fromName == "" {
				fromName = cmember.Nick
			}
		}
	}
	if log.GetLevel() > zerolog.DebugLevel {
		if mr.PeerID == 0 {
			log.PrintMessageSimple(mr.Time, userName, fromName, mr.Seq, mr.Text)
		} else if mr.UserID == 0 {
			log.PrintMessageSimple(mr.Time, peerName, fromName, mr.Seq, mr.Text)
		} else {
			log.PrintMessageSimple(mr.Time, userName+"("+peerName+")", fromName, mr.Seq, mr.Text)
		}
	} else {
		log.PrintMessage(mr.Time, peerName, userName, fromName, mr.PeerID, mr.UserID, mr.FromID, mr.Seq, mr.Text)
	}
}

func syncUinPairMessage(uinPairMessage *pb.MessageCommon_UinPairMessage) {
	log.Info().
		Int64("@uin", uinPairMessage.GetPeerUin()).
		Int64("last_read_time", uinPairMessage.GetLastReadTime()).
		Msgf("<-> [sync] %d message(s)", len(uinPairMessage.GetItems()))

	for _, msg := range uinPairMessage.GetItems() {
		peerID := msg.GetMessageHead().GetC2CTempMessageHead().GetGroupCode()
		userID := uinPairMessage.GetPeerUin()
		if msg.GetMessageHead().GetC2CCmd() == 0 {
			peerID = uinPairMessage.GetPeerUin()
			userID = 0
		}

		log.Debug().
			Str("@chat", fmt.Sprintf("@%du%d", peerID, userID)).
			Int32("@seq", msg.GetMessageHead().GetMessageSeq()).
			Int64("from", msg.GetMessageHead().GetFromUin()).
			Int64("time", msg.GetMessageHead().GetMessageTime()).
			Int32("type", msg.GetMessageHead().GetMessageType()).
			Int64("uid", msg.GetMessageHead().GetMessageUid()).
			Msg("--> [recv]")
	}
}
