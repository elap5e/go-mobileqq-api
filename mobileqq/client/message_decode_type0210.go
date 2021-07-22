package client

import (
	"encoding/hex"
	"fmt"

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

func (c *Client) decodeMessageType0210Pb(uin uint64, buf []byte) (interface{}, error) {
	msg := pb.MessageType0210{}
	if err := proto.Unmarshal(buf, &msg); err != nil {
		return nil, err
	}

	body, err := c.decodeMessageType0210(uin, msg.SubType, msg.Buffer)
	if err != nil {
		return nil, err
	} else if body == nil {
		return &msg, nil
	}
	return body, nil
}

func (c *Client) decodeMessageType0210Jce(uin uint64, buf []byte) (interface{}, error) {
	msg := MessageType0210{}
	if err := jce.Unmarshal(buf, &msg, true); err != nil {
		return nil, err
	}

	body, err := c.decodeMessageType0210(uin, msg.SubType, msg.Buffer)
	if err != nil {
		return nil, err
	} else if body == nil {
		return &msg, nil
	}
	return body, nil
}

func (c *Client) decodeMessageType0210(uin uint64, typ int64, buf []byte) (interface{}, error) {
	log.Debug().
		Msgf(">>> [0210] decode message type:0x0210 sub_type:0x%04x length:%d", typ, len(buf))
	switch typ {
	case 0x0026:
		body := pb.MessageType0210_Type0026_Body{}
		if err := proto.Unmarshal(buf, &body); err != nil {
			return nil, err
		}
		return &body, nil

	case 0x0027:
		body := pb.MessageType0210_Type0027_Body{}
		if err := proto.Unmarshal(buf, &body); err != nil {
			return nil, err
		}
		return &body, nil

	case 0x0044:
		body := pb.MessageType0210_Type0044_Body{}
		if err := proto.Unmarshal(buf, &body); err != nil {
			return nil, err
		}
		return &body, nil

	case 0x004E: // GroupBulletinNotify
		body := pb.MessageType0210_Type004E_Body{}
		if err := proto.Unmarshal(buf, &body); err != nil {
			return nil, err
		}
		return &body, nil

	case 0x008A, 0x008B: // RecallMessageNotify
		body := pb.MessageType0210_Type008A_Body{}
		if err := proto.Unmarshal(buf, &body); err != nil {
			return nil, err
		}

		for _, msg := range body.GetMessages() {
			mr := &db.MessageRecord{
				Time:   msg.GetMessageTime(),
				Seq:    msg.GetMessageSeq(),
				Uid:    int64(msg.GetMessageRandom()) | 1<<56,
				PeerID: 0,
				UserID: msg.GetFromUin(),
				FromID: msg.GetFromUin(),
				Text:   "",
				Type:   0x0210,
			}
			mr.Text = fmt.Sprintf(
				"![%s](goqq://act/recall?time=%d&seq=%d&uid=%d&peer=%d&user=%d&from=%d)",
				msg.GetWording().GetItemName(),
				mr.Time,
				mr.Seq,
				mr.Uid,
				mr.PeerID,
				mr.UserID,
				mr.FromID,
			)

			if c.db != nil {
				err := c.dbInsertMessageRecord(uin, mr)
				if err != nil {
					log.Error().Err(err).Msg(">>> [db  ] dbInsertMessageRecord")
				}
			} else {
				c.PrintMessageRecord(mr)
			}
		}
		return &body, nil

	case 0x00B3: // AddFriendNotify
		body := pb.MessageType0210_Type00B3_Body{}
		if err := proto.Unmarshal(buf, &body); err != nil {
			return nil, err
		}
		return &body, nil

	case 0x00C7: // HotFriendNotify
		body := pb.MessageType0210_Type00C7_Body{}
		if err := proto.Unmarshal(buf, &body); err != nil {
			return nil, err
		}
		return &body, nil

	case 0x0087:
	case 0x008D:
	case 0x009B:
	case 0x00AA:
	case 0x00AE:
	case 197:
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

	log.Debug().Msg(">>> [dump]\n" + hex.Dump(buf))
	return nil, nil
}
