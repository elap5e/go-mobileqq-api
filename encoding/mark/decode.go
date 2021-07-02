package mark

import (
	"encoding/base64"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/elap5e/go-mobileqq-api/pb"
)

func Unmarshal(v []byte, msg *pb.Message) error {
	body := strings.SplitN(string(v), "\n", 2)[1]
	idxes := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`).FindAllStringSubmatchIndex(body, -1)
	elems := []*pb.Element{}
	off := 0
	for _, idx := range idxes {
		if off < idx[0] {
			elems = append(elems, &pb.Element{
				Text: &pb.Text{
					Data: UnscapeString(body[off:idx[0]]),
				},
			})
		}
		off = idx[1]
		uri, _ := url.Parse(body[idx[4]:idx[5]])
		switch uri.Hostname() {
		case "sticker":
			if uri.Path == "/face" {
				if id := uri.Query().Get("id"); id != "" {
					tmp, _ := strconv.Atoi(id)
					elems = append(elems, &pb.Element{
						Face: &pb.Face{
							Index: uint32(tmp),
						},
					})
				}
			}
		case "photo":
			height, _ := strconv.Atoi(uri.Query().Get("height"))
			width, _ := strconv.Atoi(uri.Query().Get("width"))
			size, _ := strconv.Atoi(uri.Query().Get("size"))
			md5 := make([]byte, 16)
			base64.URLEncoding.Decode(md5, []byte(uri.Query().Get("md5")))
			elems = append(elems, &pb.Element{
				NotOnlineImage: &pb.NotOnlineImage{
					FilePath:  []byte(uri.Path[1:]),
					PicMd5:    md5,
					PicHeight: uint32(height),
					PicWidth:  uint32(width),
					OrigUrl:   uri.Query().Get("url"),
					FileLen:   uint32(size),
				},
			})
		}
	}
	if off < len(body) {
		elems = append(elems, &pb.Element{
			Text: &pb.Text{
				Data: UnscapeString(body[off:]),
			},
		})
	}
	if len(elems) == 0 {
		elems = append(elems, &pb.Element{
			Text: &pb.Text{
				Data: "blank",
			},
		})
	}
	msg.MessageBody = &pb.MessageBody{
		RichText: &pb.RichText{
			Elements: elems,
		},
	}
	return nil
}

func UnscapeString(s string) string {
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "%5C", "\\")
	s = strings.ReplaceAll(s, "%5D", "]")
	s = strings.ReplaceAll(s, "%21", "!")
	s = strings.ReplaceAll(s, "%5B", "[")
	return strings.ReplaceAll(s, "%25", "%")
}
