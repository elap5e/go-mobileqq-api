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

	MessageInfo       *MessageInfo `jce:",9" json:",omitempty"`
	MessageCtrlBuffer string       `jce:",10" json:",omitempty"`
	ServerBuffer      []byte       `jce:",11" json:",omitempty"`
	PingFlag          uint64       `jce:",12" json:",omitempty"`
	ServerIP          int32        `jce:",13" json:",omitempty"`
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
	c.dumpServerToClientMessage(s2c, &req)

	resp, err := c.MessageGetMessage(
		ctx, s2c.Username, NewMessageGetMessageRequest(
			0x00000000, c.syncCookie,
		),
	)
	if err != nil {
		return nil, err
	}

	type Data struct {
		ChatID uint64
		PeerID uint64
		FromID uint64
		Text   []byte
	}
	dataList := []Data{}
	infos := []MessageDeleteInfo{}
	for {
		for _, uinPairMessage := range resp.GetUinPairMessages() {
			syncUinPairMessage(uinPairMessage)

			for _, msg := range uinPairMessage.GetMessages() {
				data, err := mark.Marshal(msg)
				if err != nil {
					return nil, err
				}

				chatID := msg.GetMessageHead().GetGroupInfo().GetGroupCode()
				peerID := uint64(0)
				fromID := msg.GetMessageHead().GetFromUin()
				if msg.GetMessageHead().GetC2CCmd() != 0 {
					chatID = msg.GetMessageHead().GetC2CTempMessageHead().GetGroupCode()
					peerID = uinPairMessage.GetPeerUin()
				}
				chatName := strconv.FormatUint(chatID, 10)
				peerName := strconv.FormatUint(peerID, 10)
				fromName := strconv.FormatUint(fromID, 10)
				seq := msg.GetMessageHead().GetMessageSeq()
				text := string(data)

				log.PrintMessage(
					time.Unix(int64(msg.GetMessageHead().GetMessageTime()), 0),
					chatName, peerName, fromName, chatID, peerID, fromID, seq, text,
				)

				// message processed
				switch msg.GetMessageHead().GetMessageType() {
				case 9, 10, 31, 79, 97, 120, 132, 133, 166, 167:
					switch msg.GetMessageHead().GetC2CCmd() {
					case 11, 175:
						if s2c.Username != strconv.FormatInt(int64(fromID), 10) {
							// add to data list
							toID, _ := strconv.ParseUint(s2c.Username, 10, 64)
							dataList = append(dataList, Data{
								ChatID: chatID,
								PeerID: peerID,
								FromID: toID,
								Text:   data,
							})
						}
						infos = append(infos, MessageDeleteInfo{
							FromUin:     fromID,
							MessageTime: msg.GetMessageHead().GetMessageTime(),
							MessageSeq:  uint16(msg.GetMessageHead().GetMessageSeq()),
						})
					case 129, 131, 133:
					case 169, 241, 242, 243:
					}
				case 141:
					if msg.GetMessageHead().GetC2CCmd() == 11 {
						if s2c.Username != strconv.FormatInt(int64(fromID), 10) {
							// add to data list
							toID, _ := strconv.ParseUint(s2c.Username, 10, 64)
							dataList = append(dataList, Data{
								ChatID: chatID,
								PeerID: peerID,
								FromID: toID,
								Text:   data,
							})
						}
						infos = append(infos, MessageDeleteInfo{
							FromUin:     fromID,
							MessageTime: msg.GetMessageHead().GetMessageTime(),
							MessageSeq:  uint16(msg.GetMessageHead().GetMessageSeq()),
						})
					}
				case 208:
				case 193:
				case 734:
				case 0x0210:
					// subMsg := pb.MessageType0X210{}
					// if err = proto.Unmarshal(msg.GetMessageBody().GetMessageContent(), &subMsg); err != nil {
					// 	log.Error().Err(err).Msg("--> [0210] unmarshal 0x0210")
					// } else {
					// 	switch subMsg.GetType() {
					// 	case 138, 139:
					// 		subSubMsg := pb.MessageSubType0X8ARequest{}
					// 		if err = proto.Unmarshal(subMsg.GetContent(), &subSubMsg); err != nil {
					// 			log.Error().Err(err).Msg("--> [0210] unmarshal 0x0210_0x8a")
					// 		}
					// 		subData, _ := json.Marshal(&subSubMsg)
					// 		fmt.Println(string(subData))
					// 	}
					// }
				case 0x0211:
				case 0, 26, 64, 38, 48, 53, 61, 63:
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
	if l := len(dataList); l > 0 {
		item := dataList[l-1]
		seq := c.getNextMessageSeq(
			fmt.Sprintf("%d:%d", item.ChatID, item.PeerID),
		)
		routingHead := &pb.RoutingHead{}
		if item.ChatID == 0 {
			routingHead = &pb.RoutingHead{C2C: &pb.C2C{Uin: item.PeerID}}
		} else {
			routingHead = &pb.RoutingHead{
				GroupTemp: &pb.GroupTemp{Code: item.ChatID, ToUin: item.PeerID},
			}
		}

		msg := pb.Message{}
		if err := mark.Unmarshal(item.Text, &msg); err != nil {
			return nil, err
		}
		resp, err := c.MessageSendMessage(
			ctx, s2c.Username, NewMessageSendMessageRequest(
				routingHead,
				msg.GetContentHead(),
				msg.GetMessageBody(),
				seq,
				c.syncCookie,
			),
		)
		if err != nil {
			return nil, err
		}

		data, err := mark.Marshal(&msg)
		if err != nil {
			return nil, err
		}
		chatID := item.ChatID
		peerID := item.PeerID
		fromID := item.FromID
		chatName := strconv.FormatUint(chatID, 10)
		peerName := strconv.FormatUint(peerID, 10)
		fromName := strconv.FormatUint(fromID, 10)
		text := string(data)
		log.PrintMessage(
			time.Unix(resp.GetSendTime(), 0),
			chatName, peerName, fromName, chatID, peerID, uint64(fromID), seq, text,
		)
	}

	return NewOnlinePushMessageResponse(ctx, s2c.Username, infos, req.ServerIP, s2c.Seq)
}
