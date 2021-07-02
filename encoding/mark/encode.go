package mark

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/elap5e/go-mobileqq-api/pb"
)

func Marshal(msg *pb.Message) ([]byte, error) {
	text := fmt.Sprintf(
		"<!--mark://message?time=%d&type=%d&peer=%d&seq=%d&uid=%d&from=%d&to=%d-->\n",
		msg.GetMessageHead().GetMessageTime(),
		msg.GetMessageHead().GetMessageType(),
		msg.GetMessageHead().GetGroupInfo().GetGroupCode(),
		msg.GetMessageHead().GetMessageSeq(),
		msg.GetMessageHead().GetMessageUid(),
		msg.GetMessageHead().GetFromUin(),
		msg.GetMessageHead().GetToUin(),
	)
	for _, elem := range msg.GetMessageBody().GetRichText().GetElements() {
		if v := elem.GetText(); v != nil {
			text += EscapeString(v.GetData())
		}
		if v := elem.GetFace(); v != nil {
			text += fmt.Sprintf(
				"![sticker](mark://sticker/face?id=%d)",
				v.GetIndex(),
			)
		}
		if v := elem.GetNotOnlineImage(); v != nil {
			text += fmt.Sprintf(
				"![photo](mqqapi://photo/%s?uin=%d&type=%d&url=%s&md5=%s&height=%d&width=%d&size=%d)",
				EscapeString(string(v.GetFilePath())),
				msg.GetMessageHead().GetFromUin(),
				v.GetBizType(),
				url.QueryEscape(v.GetOrigUrl()),
				base64.URLEncoding.EncodeToString(v.GetPicMd5()),
				v.GetPicHeight(),
				v.GetPicWidth(),
				v.GetFileLen(),
			)
		}
		if v := elem.GetCustomFace(); v != nil {
			text += fmt.Sprintf(
				"![photo](mqqapi://photo/%s?uin=%d&type=5)",
				EscapeString(v.GetFilePath()),
				msg.GetMessageHead().GetFromUin(),
			)
		}
	}
	return []byte(text), nil
}

func EscapeString(s string) string {
	s = strings.ReplaceAll(s, "%", "%25")
	s = strings.ReplaceAll(s, "!", "%21")
	s = strings.ReplaceAll(s, "[", "%5B")
	s = strings.ReplaceAll(s, "\\", "%5C")
	s = strings.ReplaceAll(s, "]", "%5D")
	return strings.ReplaceAll(s, "\n", "\\n")
}
