package mark

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
	"path"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/encoding/mark/emoticon"
	"github.com/elap5e/go-mobileqq-api/pb"
	"github.com/elap5e/go-mobileqq-api/util"
)

type encoder struct {
	peerID, userID, fromID int64
}

func NewEncoder(peerID, userID, fromID int64) *encoder {
	return &encoder{peerID, userID, fromID}
}

func (enc encoder) Encode(elems []*pb.IMMessageBody_Element) ([]byte, error) {
	head, text := "", ""
	skip := new(int)
	for i, elem := range elems {
		if *skip > 0 {
			*skip--
			continue
		}
		if v := elem.GetLightApp(); v != nil {
			text += enc.encodeLightAppElement(v, skip)
		} else if v := elem.GetRichMessage(); v != nil {
			text += enc.encodeRichMessage(v)
		} else if v := elem.GetCommon(); v != nil {
			text += enc.encodeCommonElement(v, skip)
		} else if v := elem.GetText(); v != nil {
			text += enc.encodeTextMessage(v)
		} else if v := elem.GetFace(); v != nil {
			text += enc.encodeFaceElement(v)
		} else if v := elem.GetMarketFace(); v != nil {
			text += enc.encodeMarketFaceElement(v, elems[i+1].GetText(), skip)
		} else if v := elem.GetSmallEmoji(); v != nil {
			text += enc.encodeSmallEmojiElement(v, elems[i+1].GetText(), skip)
		} else if v := elem.GetCustomFace(); v != nil {
			text += enc.encodeCustomFaceElement(v)
		} else if v := elem.GetNotOnlineImage(); v != nil {
			text += enc.encodeNotOnlineImageElement(v)
		} else if v := elem.GetShakeWindow(); v != nil {
			text += enc.encodeShakeWindowElement(v)
		} else if v := elem.GetSourceMessage(); v != nil {
			text += enc.encodeSourceMessage(v)
		}
	}
	return []byte(head + text), nil
}

func (enc encoder) encodeCommonElement(elem *pb.IMMessageBody_CommonElement, skip *int) string {
	switch elem.GetServiceType() {
	case 33: // extra face
		info := pb.CommonElement_ServiceType33{}
		_ = proto.Unmarshal(elem.GetBuffer(), &info)
		id := emoticon.FaceType(info.GetIndex())
		return fmt.Sprintf(
			"![%s](goqq://res/face?id=%d)",
			id.String(),
			id,
		)
	case 37: // extra big face
		*skip++
		info := pb.CommonElement_ServiceType37{}
		_ = proto.Unmarshal(elem.GetBuffer(), &info)
		id := emoticon.FaceType(info.GetQSid())
		return fmt.Sprintf(
			"![%s](goqq://res/face?id=%d&pid=%s&sid=%s)",
			id.String(),
			id,
			base64.URLEncoding.EncodeToString(info.GetPackageId()),
			base64.URLEncoding.EncodeToString(info.GetStickerId()),
		)
	}
	return ""
}

func (enc encoder) encodeCustomFaceElement(elem *pb.IMMessageBody_CustomFace) string {
	hash := elem.GetFileMd5()
	return fmt.Sprintf(
		"![%s](goqq://res/image?md5=%s&type=5&size=%d&h=%d&w=%d)",
		util.HashToString(hash)+path.Ext(string(elem.GetFilePath())),
		base64.URLEncoding.EncodeToString(hash),
		elem.GetFileSize(),
		elem.GetHeight(),
		elem.GetWidth(),
	)
}

func (enc encoder) encodeFaceElement(elem *pb.IMMessageBody_Face) string {
	id := emoticon.FaceType(elem.GetIndex())
	return fmt.Sprintf(
		"![%s](goqq://res/face?id=%d)",
		id.String(),
		id,
	)
}

func (enc encoder) encodeLightAppElement(elem *pb.IMMessageBody_LightAppElement, skip *int) string {
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

func (enc encoder) encodeMarketFaceElement(elem *pb.IMMessageBody_MarketFace, text *pb.IMMessageBody_Text, skip *int) string {
	*skip++
	name := string(elem.GetFaceName())
	if name == "" {
		name = text.GetText()
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

func (enc encoder) encodeNotOnlineImageElement(elem *pb.IMMessageBody_NotOnlineImage) string {
	hash := elem.GetFileMd5()
	return fmt.Sprintf(
		"![%s](goqq://res/image?md5=%s&type=%d&size=%d&h=%d&w=%d)",
		util.HashToString(hash)+path.Ext(string(elem.GetFilePath())),
		base64.URLEncoding.EncodeToString(hash),
		elem.GetBizType(),
		elem.GetFileSize(),
		elem.GetHeight(),
		elem.GetWidth(),
	)
}

func (enc encoder) encodeRichMessage(elem *pb.IMMessageBody_RichMessage) string {
	data := elem.GetTemplate()[1:]
	if elem.GetTemplate()[0] == 1 {
		reader, _ := zlib.NewReader(bytes.NewReader(data))
		defer reader.Close()
		var buf bytes.Buffer
		io.Copy(&buf, reader)
		data = buf.Bytes()
	}
	return string(data)
}

func (enc encoder) encodeShakeWindowElement(elem *pb.IMMessageBody_ShakeWindow) string {
	return fmt.Sprintf(
		"![[shakeWindow]](goqq://act/shakeWindow?uin=%d&type=%d)",
		elem.GetUin(),
		elem.GetType(),
	)
}

func (enc encoder) encodeSmallEmojiElement(elem *pb.IMMessageBody_SmallEmoji, text *pb.IMMessageBody_Text, skip *int) string {
	*skip++
	return fmt.Sprintf(
		"![%s](goqq://res/smallEmoji?id=%d&type=%d)",
		text.GetText(),
		elem.GetPackIdSum(),
		elem.GetImageType(),
	)
}

func (enc encoder) encodeSourceMessage(elem *pb.IMMessageBody_SourceMessage) string {
	return fmt.Sprintf(
		"<!--goqq://msg/reply?time=%d&peer=%d&user=%d&from=%d&seq=%d-->\n",
		elem.GetTime(),
		enc.peerID,
		enc.userID,
		elem.GetFromUin(),
		elem.GetOrigSeqs()[0],
	)
}

func (enc encoder) encodeTextMessage(elem *pb.IMMessageBody_Text) string {
	attr6Buf := elem.GetAttribute6Buffer()
	if len(attr6Buf) < 13 {
		return escape(elem.GetText())
	} else {
		uin := uint64(attr6Buf[7])<<24 + uint64(attr6Buf[8])<<16 + uint64(attr6Buf[9])<<8 + uint64(attr6Buf[10])
		return fmt.Sprintf(
			"![%s](goqq://act/at?uin=%d)",
			escape(elem.GetText()),
			uin,
		)
	}
}

func escape(s string) string {
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
