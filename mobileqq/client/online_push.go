package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/encoding/jce"
	"github.com/elap5e/go-mobileqq-api/encoding/uni"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
	"google.golang.org/protobuf/proto"
)

type OnlinePushRequest struct {
	Uin             int64             `jce:",0" json:",omitempty"`
	MessageTime     int64             `jce:",1" json:",omitempty"`
	MessageInfos    []*MessageInfo    `jce:",2" json:",omitempty"`
	ServerIP        uint32            `jce:",3" json:",omitempty"`
	SyncCookie      []byte            `jce:",4" json:",omitempty"`
	UinPairMessages []*UinPairMessage `jce:",5" json:",omitempty"`
	Previews        map[string][]byte `jce:",6" json:",omitempty"`
	UserActive      int32             `jce:",7" json:",omitempty"`
	GeneralFlag     int32             `jce:",12" json:",omitempty"`
}

type UinPairMessage struct {
	LastReadTime     int64          `jce:",1" json:",omitempty"`
	PeerUin          int64          `jce:",2" json:",omitempty"`
	MessageCompleted int64          `jce:",3" json:",omitempty"`
	MessageInfos     []*MessageInfo `jce:",4" json:",omitempty"`
}

type OnlinePushResponse struct {
	Type   int32  `jce:",1" json:",omitempty"`
	Seq    int64  `jce:",2" json:",omitempty"`
	Buffer []byte `jce:",3" json:",omitempty"`
}

type MessageType0x0210 struct {
	SubType int64  `jce:",0" json:",omitempty"`
	Buffer  []byte `jce:",10" json:",omitempty"`
}

func NewOnlinePushMessageResponse(
	ctx context.Context,
	username string,
	infos []MessageDeleteInfo,
	serverIP uint32,
	seq int32,
) (*codec.ClientToServerMessage, error) {
	if len(infos) == 0 {
		return nil, fmt.Errorf("zero length")
	}

	uin, err := strconv.ParseInt(username, 10, 64)
	if err != nil {
		return nil, err
	}
	resp := OnlinePushMessageResponse{
		Uin:      uin,
		Infos:    infos,
		ServerIP: serverIP,
	}
	buf, err := uni.Marshal(ctx, &uni.Message{
		Version:     0x0003,
		PacketType:  0x00,
		MessageType: 0x00000000,
		RequestID:   seq,
		ServantName: "OnlinePush",
		FuncName:    "SvcRespPushMsg",
		Buffer:      []byte{},
		Timeout:     0x00000000,
		Context:     map[string]string{},
		Status:      map[string]string{},
	}, map[string]interface{}{
		"resp": resp,
	})
	if err != nil {
		return nil, err
	}
	return &codec.ClientToServerMessage{
		Username:      username,
		Seq:           uint32(seq),
		ServiceMethod: ServiceMethodOnlinePushResponse,
		Buffer:        buf,
		Simple:        false,
	}, nil
}

func (c *Client) handleOnlinePushRequest(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error) {
	msg := uni.Message{}
	push := OnlinePushRequest{}
	if err := uni.Unmarshal(ctx, s2c.Buffer, &msg, map[string]interface{}{
		"req": &push,
	}); err != nil {
		return nil, err
	}
	dumpServerToClientMessage(s2c, &push)

	uin, _ := strconv.ParseUint(s2c.Username, 10, 64)
	infos := []MessageDeleteInfo{}
	for _, msg := range push.MessageInfos {
		switch msg.MessageType {
		case 0x0210:
			v := MessageType0x0210{}
			if err := jce.Unmarshal(msg.MessageBytes, &v, true); err != nil {
				return nil, err
			}
			log.Debug().Msgf(">>> [0x210] subType:%x(%d)", v.SubType, len(v.Buffer))
			switch v.SubType {
			case 0x8a:
				notify := pb.MessageType0X0210SubType0X8A{}
				if err := proto.Unmarshal(v.Buffer, &notify); err != nil {
					return nil, err
				}
				dumpServerToClientMessage(s2c, &notify)

				for _, subMsg := range notify.GetInfo() {
					mr := &db.MessageRecord{
						Time:   subMsg.GetTime(),
						Seq:    subMsg.GetSeq(),
						Uid:    int64(subMsg.GetRandom()) | 1<<56,
						PeerID: 0,
						UserID: subMsg.GetFromUin(),
						FromID: subMsg.GetFromUin(),
						Text:   "",
						Type:   0x0210,
					}
					mr.Text = "messageRecall"

					c.PrintMessageRecord(mr)
					if c.db != nil {
						err := c.dbInsertMessageRecord(uin, mr)
						if err != nil {
							log.Error().Err(err).Msg(">>> [db  ] dbInsertMessageRecord")
						}
					}
				}
			}
		case 0x02DC:
			if len(msg.MessageBytes) > 4 {
				buf := bytes.NewBuffer(msg.MessageBytes)
				_, _ = buf.ReadUint32()
				subType, _ := buf.ReadUint8()
				log.Debug().Msgf(">>> [0x2DC] subType:%x(%d)", subType, len(msg.MessageBytes))
				switch subType {
				case 0x03:
				case 0x0c, 0x0e:
				case 0x10, 0x11, 0x14, 0x15:
					if len(msg.MessageBytes) > 7 {
						notify := pb.NotifyMessageBody{}
						if err := proto.Unmarshal(msg.MessageBytes[7:], &notify); err != nil {
							return nil, err
						}
						dumpServerToClientMessage(s2c, &notify)
						if v := notify.GetMessageRecall(); v != nil {
							for _, msg := range v.GetRecalledMessageList() {
								mr := &db.MessageRecord{
									Time:   msg.GetTime(),
									Seq:    msg.GetSeq(),
									Uid:    int64(msg.GetRandom()) | 1<<56,
									PeerID: notify.GetGroupCode(),
									UserID: 0,
									FromID: v.GetUin(),
									Text:   "",
									Type:   0x02DC,
								}
								mr.Text = "messageRecall: " + v.GetMessageWordingInfo().GetName()

								c.PrintMessageRecord(mr)
								if c.db != nil {
									err := c.dbInsertMessageRecord(uin, mr)
									if err != nil {
										log.Error().Err(err).Msg(">>> [db  ] dbInsertMessageRecord")
									}
								}
							}
						}
					}
				}
			}
		}
		infos = append(infos, MessageDeleteInfo{
			FromUin:       msg.FromUin,
			MessageTime:   msg.MessageTime,
			MessageSeq:    msg.MessageSeq,
			MessageCookie: msg.MessageCookies,
		})
	}

	return NewOnlinePushMessageResponse(ctx, s2c.Username, infos, push.ServerIP, int32(s2c.Seq))
}
