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
	head, body := "", ""
	skip := new(int)
	for i, elem := range elems {
		if *skip > 0 {
			*skip--
			continue
		}
		if v := elem.GetLightApp(); v != nil {
			body += enc.encodeLightAppElement(v, skip)
		} else if v := elem.GetRichMessage(); v != nil {
			body += enc.encodeRichMessage(v)
		} else if v := elem.GetCommon(); v != nil {
			text := new(pb.IMMessageBody_Text)
			if len(elems) > i+1 {
				text = elems[i+1].GetText()
			}
			body += enc.encodeCommonElement(v, text, skip)
		} else if v := elem.GetText(); v != nil {
			body += enc.encodeTextMessage(v)
		} else if v := elem.GetFace(); v != nil {
			body += enc.encodeFaceElement(v)
		} else if v := elem.GetMarketFace(); v != nil {
			body += enc.encodeMarketFaceElement(v, elems[i+1].GetText(), skip)
		} else if v := elem.GetSmallEmoji(); v != nil {
			body += enc.encodeSmallEmojiElement(v, elems[i+1].GetText(), skip)
		} else if v := elem.GetCustomFace(); v != nil {
			body += enc.encodeCustomFaceElement(v)
		} else if v := elem.GetNotOnlineImage(); v != nil {
			body += enc.encodeNotOnlineImageElement(v)
		} else if v := elem.GetShakeWindow(); v != nil {
			body += enc.encodeShakeWindowElement(v)
		} else if v := elem.GetSourceMessage(); v != nil {
			head += enc.encodeSourceMessage(v)
		} else if v := elem.GetAnonymousGroupMessage(); v != nil {
			head += enc.encodeAnonymousGroupMessage(v)
		}
	}
	return []byte(head + body), nil
}

func (enc encoder) encodeAnonymousGroupMessage(elem *pb.IMMessageBody_AnonymousGroupMessage) string {
	return fmt.Sprintf(
		"![[anonymous(%s)]](goqq://act/anonymous?id=%s&bid=%d&exp=%d)",
		elem.GetAnonymousNick(),
		base64.URLEncoding.EncodeToString(elem.AnonymousId),
		elem.GetBubbleId(),
		elem.GetExpireTime(),
	)
}

func (enc encoder) encodeCommonElement(elem *pb.IMMessageBody_CommonElement, text *pb.IMMessageBody_Text, skip *int) string {
	switch elem.GetServiceType() {
	case 2: // poke
		if id := elem.GetBusinessType(); id == 0 {
			return "![[shakeWindow]](goqq://act/shakeWindow)"
		} else {
			*skip++
			return fmt.Sprintf(
				"![%s](goqq://act/poke?id=%d&buf=%s)",
				text.GetText(),
				id,
				base64.URLEncoding.EncodeToString(elem.GetBuffer()),
			)
		}
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
			"![%s](goqq://res/face?id=%d&pid=%s&sid=%s&rsv=%s)",
			id.String(),
			id,
			base64.URLEncoding.EncodeToString(info.GetPackageId()),
			base64.URLEncoding.EncodeToString(info.GetStickerId()),
			base64.URLEncoding.EncodeToString(text.GetPbReserve()),
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
	return "![[shakeWindow]](goqq://act/shakeWindow)"
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
		"![[reply]](goqq://act/reply?time=%d&peer=%d&user=%d&from=%d&seq=%d)",
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
