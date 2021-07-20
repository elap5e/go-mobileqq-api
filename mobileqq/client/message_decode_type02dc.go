package client

import (
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
	_, _ = buf.ReadUint32()
	subType, _ := buf.ReadUint8()
	log.Debug().Msgf(">>> [02DC] subType:0x%x(%d)", subType, len(p))

	switch subType {
	case 0x03:
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
				mr.Text = "messageRecall: " + notify.GetWording().GetItemName()

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
	return nil, nil
}
