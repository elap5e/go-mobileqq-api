package client

import (
	"encoding/hex"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
	"github.com/elap5e/go-mobileqq-api/pb"
)

func (c *Client) decodeMessageType02DC(uin uint64, p []byte) (interface{}, error) {
	if len(p) < 5 {
		return nil, nil
	}
	buf := bytes.NewBuffer(p)
	peerID, _ := buf.ReadUint32()
	subType, _ := buf.ReadUint8()
	log.Debug().
		Uint32("@peer", peerID).
		Msgf(">>> [02DC] subType:0x%x(%d)", subType, len(p))

	switch subType {
	case 0x03:
		buf.Skip(0x05)
		time, _ := buf.ReadUint32()
		l, _ := buf.ReadUint16()
		mrs := []*db.MessageRecord{}
		for i := 0; i < int(l); i++ {
			mr := db.MessageRecord{
				Time:   int64(time),
				Seq:    0,
				Uid:    0,
				PeerID: int64(peerID),
				UserID: 0,
				FromID: 0,
				Text:   "",
				Type:   0x02DC,
			}
			mr.Text, _ = buf.ReadString()

			c.PrintMessageRecord(&mr)
			mrs = append(mrs, &mr)
		}
		return &mrs, nil
	case 0x0C, 0x0E:
	case 0x10, 0x11, 0x14, 0x15:
		if len(p) < 8 {
			return nil, nil
		}
		body := pb.OIDB_Type0857_NotifyMessageBody{}
		if err := proto.Unmarshal(p[7:], &body); err != nil {
			return nil, err
		}

		if notify := body.GetRecall(); notify != nil {
			for _, msg := range notify.GetMessages() {
				mr := &db.MessageRecord{
					Time:   msg.GetMessageTime(),
					Seq:    msg.GetMessageSeq(),
					Uid:    int64(msg.GetMessageRandom()) | 1<<56,
					PeerID: body.GetGroupCode(),
					UserID: 0,
					FromID: notify.GetUin(),
					Text:   "",
					Type:   0x02DC,
				}
				mr.Text = fmt.Sprintf(
					"![%s](goqq://act/recall?time=%d&seq=%d&uid=%d&peer=%d&user=%d&from=%d)",
					notify.GetWording().GetItemName(),
					mr.Time,
					mr.Seq,
					mr.Uid,
					mr.PeerID,
					mr.UserID,
					mr.FromID,
				)

				c.PrintMessageRecord(mr)
				if c.db != nil {
					err := c.dbInsertMessageRecord(uin, mr)
					if err != nil {
						log.Error().Err(err).Msg(">>> [db  ] dbInsertMessageRecord")
					}
				}
			}
		}
		return &body, nil
	}

	log.Debug().Msg(">>> [dump]\n" + hex.Dump(p))
	return nil, nil
}
