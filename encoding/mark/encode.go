package mark

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/elap5e/go-mobileqq-api/pb"
)

func Marshal(msg *pb.Message) ([]byte, error) {
	head := fmt.Sprintf(
		"<!--goqq://msg/info?time=%d&type=%d&peer=%d&seq=%d&uid=%d&from=%d&to=%d-->\n",
		msg.GetMessageHead().GetMessageTime(),
		msg.GetMessageHead().GetMessageType(),
		msg.GetMessageHead().GetGroupInfo().GetGroupCode(),
		msg.GetMessageHead().GetMessageSeq(),
		msg.GetMessageHead().GetMessageUid(),
		msg.GetMessageHead().GetFromUin(),
		msg.GetMessageHead().GetToUin(),
	)
	text := ""
	skip := 0
	for _, elem := range msg.GetMessageBody().GetRichText().GetElements() {
		if skip > 0 {
			skip--
			continue
		}
		if v := elem.GetText(); v != nil {
			attr6Buf := v.GetAttr6Buffer()
			if len(attr6Buf) < 13 {
				if id, ok := mapStringFaceID[v.GetData()]; !ok {
					text += EscapeString(v.GetData())
				} else {
					text += fmt.Sprintf(
						"![%s](goqq://res/face?id=%d)",
						mapFaceIDString[id],
						id,
					)
				}
			} else {
				uin := uint64(attr6Buf[7])<<24 + uint64(attr6Buf[8])<<16 + uint64(attr6Buf[9])<<8 + uint64(attr6Buf[10])
				text += fmt.Sprintf(
					"![%s](goqq://act/at?uin=%d)",
					EscapeString(v.GetData()),
					uin,
				)
			}
		} else if v := elem.GetFace(); v != nil {
			text += fmt.Sprintf(
				"![%s](goqq://res/face?id=%d)",
				mapFaceIDString[v.GetIndex()],
				v.GetIndex(),
			)
		} else if v := elem.GetMarketFace(); v != nil {
			text += fmt.Sprintf(
				"![%s](goqq://res/marketFace?id=%s&tabId=%d&key=%s)",
				string(v.GetFaceName()),
				base64.URLEncoding.EncodeToString(v.GetFaceId()),
				v.GetTabId(),
				base64.URLEncoding.EncodeToString(v.GetKey()),
			)
			skip++
		} else if v := elem.GetCustomFace(); v != nil {
			text += fmt.Sprintf(
				"![%s](goqq://res/image?md5=%s&type=%d&uin=%d&size=%d&h=%d&w=%d)",
				EscapeString(string(v.GetFilePath())),
				base64.URLEncoding.EncodeToString(v.GetMd5()),
				v.GetBizType(),
				msg.GetMessageHead().GetFromUin(),
				v.GetSize(),
				v.GetHeight(),
				v.GetWidth(),
			)
		} else if v := elem.GetNotOnlineImage(); v != nil {
			text += fmt.Sprintf(
				"![%s](goqq://res/image?md5=%s&type=%d&uin=%d&size=%d&h=%d&w=%d)",
				EscapeString(string(v.GetFilePath())),
				base64.URLEncoding.EncodeToString(v.GetPicMd5()),
				v.GetBizType(),
				msg.GetMessageHead().GetFromUin(),
				v.GetFileLen(),
				v.GetPicHeight(),
				v.GetPicWidth(),
			)
		} else if v := elem.GetShakeWindow(); v != nil {
			text += fmt.Sprintf(
				"![[shakeWindow]](goqq://act/shakeWindow?uin=%d&type=%d)",
				v.GetUin(),
				v.GetType(),
			)
		} else if v := elem.GetSourceMessage(); v != nil {
			head += fmt.Sprintf(
				"<!--goqq://msg/reply?time=%d&peer=%d&seq=%d&from=%d-->\n",
				v.GetTime(),
				msg.GetMessageHead().GetGroupInfo().GetGroupCode(),
				v.GetOrigSeqs()[0],
				v.GetFromUin(),
			)
		}
	}
	return []byte(head + text), nil
}

func EscapeString(s string) string {
	s = strings.ReplaceAll(s, "%", "%25")
	s = strings.ReplaceAll(s, "!", "%21")
	s = strings.ReplaceAll(s, "[", "%5B")
	s = strings.ReplaceAll(s, "]", "%5D")
	// s = strings.ReplaceAll(s, "\\", "%5C")
	// s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}
