package mark

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/encoding/mark/emoticon"
	"github.com/elap5e/go-mobileqq-api/pb"
)

type marshaler struct {
	peerID, userID, fromID int64
}

func NewMarshaler(peerID, userID, fromID int64) *marshaler {
	return &marshaler{peerID, userID, fromID}
}

func (m marshaler) Marshal(elems []*pb.Element) ([]byte, error) {
	head, text := "", ""
	skip := new(int)
	for i, elem := range elems {
		if *skip > 0 {
			*skip--
			continue
		}
		if v := elem.GetLightApp(); v != nil {
			text += m.marshalLightAppElement(v, skip)
		} else if v := elem.GetRichMessage(); v != nil {
			text += m.marshalRichMessage(v)
		} else if v := elem.GetCommon(); v != nil {
			text += m.marshalCommonElement(v, skip)
		} else if v := elem.GetText(); v != nil {
			text += m.marshalTextMessage(v)
		} else if v := elem.GetFace(); v != nil {
			text += m.marshalFaceElement(v)
		} else if v := elem.GetMarketFace(); v != nil {
			text += m.marshalMarketFaceElement(v, elems[i+1].GetText(), skip)
		} else if v := elem.GetSmallEmoji(); v != nil {
			text += m.marshalSmallEmojiElement(v, elems[i+1].GetText(), skip)
		} else if v := elem.GetCustomFace(); v != nil {
			text += m.marshalCustomFaceElement(v)
		} else if v := elem.GetNotOnlineImage(); v != nil {
			text += m.marshalNotOnlineImage(v)
		} else if v := elem.GetShakeWindow(); v != nil {
			text += m.marshalShakeWindowElement(v)
		} else if v := elem.GetSourceMessage(); v != nil {
			text += m.marshalSourceMessage(v)
		}
	}
	return []byte(head + text), nil
}

func (m marshaler) marshalCommonElement(elem *pb.CommonElement, skip *int) string {
	switch elem.GetServiceType() {
	case 33: // extra face
		info := pb.MessageElementInfoServiceType33{}
		_ = proto.Unmarshal(elem.GetBuffer(), &info)
		id := emoticon.FaceType(info.GetIndex())
		return fmt.Sprintf(
			"![%s](goqq://res/face?id=%d)",
			id.String(),
			id,
		)
	case 37: // extra big face
		*skip++
		info := pb.MessageElementInfoServiceType37{}
		_ = proto.Unmarshal(elem.GetBuffer(), &info)
		id := emoticon.FaceType(info.GetQsId())
		return fmt.Sprintf(
			"![%s](goqq://res/face?id=%d&pid=%s&sid=%s)",
			id.String(),
			id,
			base64.URLEncoding.EncodeToString(info.GetPackId()),
			base64.URLEncoding.EncodeToString(info.GetStickerId()),
		)
	}
	return ""
}

func (m marshaler) marshalCustomFaceElement(elem *pb.CustomFace) string {
	return fmt.Sprintf(
		"![%s](goqq://res/image?md5=%s&type=%d&uin=%d&size=%d&h=%d&w=%d)",
		EscapeString(string(elem.GetFilePath())),
		base64.URLEncoding.EncodeToString(elem.GetMd5()),
		elem.GetBizType(),
		m.fromID,
		elem.GetSize(),
		elem.GetHeight(),
		elem.GetWidth(),
	)
}

func (m marshaler) marshalFaceElement(elem *pb.Face) string {
	id := emoticon.FaceType(elem.GetIndex())
	return fmt.Sprintf(
		"![%s](goqq://res/face?id=%d)",
		id.String(),
		id,
	)
}

func (m marshaler) marshalLightAppElement(elem *pb.LightAppElement, skip *int) string {
	*skip++
	data := elem.GetData()[1:]
	if elem.GetData()[0] == 1 {
		reader, _ := zlib.NewReader(bytes.NewReader(data))
		defer reader.Close()
		var buf bytes.Buffer
		io.Copy(&buf, reader)
		data = buf.Bytes()
	}
	return string(data)
}

func (m marshaler) marshalMarketFaceElement(elem *pb.MarketFace, text *pb.Text, skip *int) string {
	*skip++
	name := string(elem.GetFaceName())
	if name == "" {
		name = text.GetData()
	}
	return fmt.Sprintf(
		"![%s](goqq://res/marketFace?id=%s&tabId=%d&key=%s&h=%d&w=%d&p=%s)",
		name,
		base64.URLEncoding.EncodeToString(elem.GetFaceId()),
		elem.GetTabId(),
		base64.URLEncoding.EncodeToString(elem.GetKey()),
		elem.GetImageHeight(),
		elem.GetImageWidth(),
		base64.URLEncoding.EncodeToString(elem.GetMobileParam()),
	)
}

func (m marshaler) marshalNotOnlineImage(elem *pb.NotOnlineImage) string {
	return fmt.Sprintf(
		"![%s](goqq://res/image?md5=%s&type=%d&uin=%d&size=%d&h=%d&w=%d)",
		EscapeString(string(elem.GetFilePath())),
		base64.URLEncoding.EncodeToString(elem.GetPictureMd5()),
		elem.GetBizType(),
		m.fromID,
		elem.GetFileSize(),
		elem.GetPictureHeight(),
		elem.GetPictureWidth(),
	)
}

func (m marshaler) marshalRichMessage(elem *pb.RichMessage) string {
	data := elem.GetTemplate1()[1:]
	if elem.GetTemplate1()[0] == 1 {
		reader, _ := zlib.NewReader(bytes.NewReader(data))
		defer reader.Close()
		var buf bytes.Buffer
		io.Copy(&buf, reader)
		data = buf.Bytes()
	}
	return string(data)
}

func (m marshaler) marshalShakeWindowElement(elem *pb.ShakeWindow) string {
	return fmt.Sprintf(
		"![[shakeWindow]](goqq://act/shakeWindow?uin=%d&type=%d)",
		elem.GetUin(),
		elem.GetType(),
	)
}

func (m marshaler) marshalSmallEmojiElement(elem *pb.SmallEmoji, text *pb.Text, skip *int) string {
	*skip++
	return fmt.Sprintf(
		"![%s](goqq://res/smallEmoji?id=%d&type=%d)",
		text.GetData(),
		elem.GetPackIdSum(),
		elem.GetImageType(),
	)
}

func (m marshaler) marshalSourceMessage(elem *pb.SourceMessage) string {
	return fmt.Sprintf(
		"<!--goqq://msg/reply?time=%d&peer=%d&user=%d&from=%d&seq=%d-->\n",
		elem.GetTime(),
		m.peerID,
		m.userID,
		elem.GetFromUin(),
		elem.GetOrigSeqs()[0],
	)
}

func (m marshaler) marshalTextMessage(elem *pb.Text) string {
	attr6Buf := elem.GetAttr6Buffer()
	if len(attr6Buf) < 13 {
		return EscapeString(elem.GetData())
	} else {
		uin := uint64(attr6Buf[7])<<24 + uint64(attr6Buf[8])<<16 + uint64(attr6Buf[9])<<8 + uint64(attr6Buf[10])
		return fmt.Sprintf(
			"![%s](goqq://act/at?uin=%d)",
			EscapeString(elem.GetData()),
			uin,
		)
	}
}

func Marshal(msg *pb.Message) ([]byte, error) {
	// head := fmt.Sprintf(
	// 	"<!--goqq://msg/info?time=%d&type=%d&peer=%d&seq=%d&uid=%d&from=%d&to=%d-->\n",
	// 	msg.GetMessageHead().GetMessageTime(),
	// 	msg.GetMessageHead().GetMessageType(),
	// 	msg.GetMessageHead().GetGroupInfo().GetGroupCode(),
	// 	msg.GetMessageHead().GetMessageSeq(),
	// 	msg.GetMessageHead().GetMessageUid(),
	// 	msg.GetMessageHead().GetFromUin(),
	// 	msg.GetMessageHead().GetToUin(),
	// )
	peerID := msg.GetMessageHead().GetGroupInfo().GetGroupCode()
	userID := msg.GetMessageHead().GetGroupInfo().GetGroupCode()
	fromID := msg.GetMessageHead().GetFromUin()
	m := NewMarshaler(peerID, userID, fromID)
	return m.Marshal(msg.GetMessageBody().GetRichText().GetElements())
}

func EscapeString(s string) string {
	s = strings.ReplaceAll(s, "%", "%25")
	s = strings.ReplaceAll(s, "![", "%21%5B")
	s = strings.ReplaceAll(s, "](", "%5D%28")
	// s = strings.ReplaceAll(s, "!", "%21")
	// s = strings.ReplaceAll(s, "[", "%5B")
	// s = strings.ReplaceAll(s, "\\", "%5C")
	// s = strings.ReplaceAll(s, "]", "%5D")
	// s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}
