package client

import (
	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/encoding/jce"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
	"github.com/elap5e/go-mobileqq-api/pb"
)

type MessageType0210 struct {
	SubType int64  `jce:",0" json:",omitempty"`
	Buffer  []byte `jce:",10" json:",omitempty"`
}

func (c *Client) decodeMessageType0210(uin uint64, p []byte) (interface{}, error) {
	msg := MessageType0210{}
	if err := jce.Unmarshal(p, &msg, true); err != nil {
		return nil, err
	}
	log.Debug().Msgf(">>> [0210] subType:0x%x(%d)", msg.SubType, len(msg.Buffer))

	switch msg.SubType {
	case 0x0027:
		body := pb.MessageType0210_Type0027_Body{}
		if err := proto.Unmarshal(msg.Buffer, &body); err != nil {
			return nil, err
		}
		return &body, nil

	case 0x008A, 0x008B:
		body := pb.MessageType0210_Type008A_Body{}
		if err := proto.Unmarshal(msg.Buffer, &body); err != nil {
			return nil, err
		}

		for _, item := range body.GetItems() {
			mr := &db.MessageRecord{
				Time:   item.GetMessageTime(),
				Seq:    item.GetMessageSeq(),
				Uid:    int64(item.GetMessageRandom()) | 1<<56,
				PeerID: 0,
				UserID: item.GetFromUin(),
				FromID: item.GetFromUin(),
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
		return &body, nil

	case 0x00B3:
		body := pb.MessageType0210_Type00B3_Body{}
		if err := proto.Unmarshal(msg.Buffer, &body); err != nil {
			return nil, err
		}
		return &body, nil

	case 0x0087:
	case 0x008D:
	case 0x009B:
	case 0x00AA:
	case 0x00AE:
	case 197:
	case 199:
	case 203:
	case 215:
	case 220:
	case 232:
	case 238:
	case 244:
	case 249:
	case 251:
	case 253:
	case 254:
	case 256:
	case 258:
	case 260:
	case 264:
	case 273:
	case 278:
	case 281:
	case 286:
	case 287:
	case 290, 291:
	case 297:
	case 307:
	case 321:
	}
	return &msg, nil
}
