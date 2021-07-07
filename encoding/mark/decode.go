package mark

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/encoding/mark/emoticon"
	"github.com/elap5e/go-mobileqq-api/pb"
)

func Unmarshal(v []byte, msg *pb.Message) error {
	// body := strings.SplitN(string(v), "\n", 2)[1]
	body := string(v)
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
		if uri.Scheme == "goqq" {
			switch uri.Hostname() {
			default:
				elems = append(elems, &pb.Element{
					Text: &pb.Text{
						Data: body[idx[2]:idx[3]],
					},
				})
			case "res":
				switch uri.Path {
				default:
					elems = append(elems, &pb.Element{
						Text: &pb.Text{
							Data: body[idx[2]:idx[3]],
						},
					})
				case "/face":
					id := uri.Query().Get("id")
					if id != "" {
						tmp, _ := strconv.ParseUint(id, 10, 32)
						if tmp < 260 {
							elems = append(elems, &pb.Element{
								Face: &pb.Face{
									Index: uint32(tmp),
								},
							})
						} else {
							des := []byte(emoticon.FaceType(tmp).String())
							buf, _ := proto.Marshal(&pb.MessageElementInfoServiceType33{
								Index:  uint32(tmp),
								Text:   des,
								Compat: des,
							})
							elems = append(elems, &pb.Element{
								CommonElement: &pb.CommonElement{
									ServiceType:  33,
									PbElement:    buf,
									BusinessType: 1,
								},
							})
						}
					}
				case "/marketFace":
					id, _ := base64.URLEncoding.DecodeString(uri.Query().Get("id"))
					tabId, _ := strconv.ParseUint(uri.Query().Get("tabId"), 10, 32)
					key, _ := base64.URLEncoding.DecodeString(uri.Query().Get("key"))
					h, _ := strconv.ParseUint(uri.Query().Get("h"), 10, 32)
					w, _ := strconv.ParseUint(uri.Query().Get("w"), 10, 32)
					if len(id) != 0 {
						elems = append(elems, &pb.Element{
							MarketFace: &pb.MarketFace{
								FaceName:    []byte(body[idx[2]:idx[3]]),
								FaceId:      id,
								TabId:       uint32(tabId),
								Key:         key,
								ImageHeight: uint32(h),
								ImageWidth:  uint32(w),
								ItemType:    6,
								FaceInfo:    1,
								SubType:     3,
							},
						})
						elems = append(elems, &pb.Element{
							Text: &pb.Text{
								Data: body[idx[2]:idx[3]],
							},
						})
					}
				case "/image":
					md5, _ := base64.URLEncoding.DecodeString(uri.Query().Get("md5"))
					typ, _ := strconv.ParseUint(uri.Query().Get("type"), 10, 32)
					uin := uri.Query().Get("uin")
					size, _ := strconv.ParseUint(uri.Query().Get("size"), 10, 32)
					h, _ := strconv.ParseUint(uri.Query().Get("h"), 10, 32)
					w, _ := strconv.ParseUint(uri.Query().Get("w"), 10, 32)
					path := fmt.Sprintf(
						"/%s//%s-%d-%s",
						uin, uin, rand.Intn(1e10), strings.ToUpper(hex.EncodeToString(md5)),
					)
					elems = append(elems, &pb.Element{
						NotOnlineImage: &pb.NotOnlineImage{
							PicMd5:       md5,
							BizType:      uint32(typ),
							FileLen:      uint32(size),
							PicHeight:    uint32(h),
							PicWidth:     uint32(w),
							FilePath:     []byte(body[idx[2]:idx[3]]),
							ResId:        []byte(path),
							DownloadPath: []byte(path),
							OrigUrl:      "/offpic_new" + path + "/0?term=2",
						},
					})
				}
			case "act":
				switch uri.Path {
				default:
					elems = append(elems, &pb.Element{
						Text: &pb.Text{
							Data: body[idx[2]:idx[3]],
						},
					})
				case "/at":
					attr6Buf := make([]byte, 13)
					attr6Buf[1] = 0x01
					binary.BigEndian.PutUint16(attr6Buf[4:], uint16(len([]rune(body[idx[2]:idx[3]]))))
					uin, _ := strconv.ParseUint(uri.Query().Get("uin"), 10, 32)
					if uin == 0 {
						attr6Buf[6] = 0x01
					} else {
						binary.BigEndian.PutUint32(attr6Buf[7:], uint32(uin))
					}
					elems = append(elems, &pb.Element{
						Text: &pb.Text{
							Data:        body[idx[2]:idx[3]],
							Attr6Buffer: attr6Buf,
						},
					})
				case "/shakeWindow":
					elems = append(elems, &pb.Element{
						ShakeWindow: &pb.ShakeWindow{
							Uin: 0,
						},
					})
				}
			}
		} else {
			// TODO: file
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
	// s = strings.ReplaceAll(s, "\\n", "\n")
	// s = strings.ReplaceAll(s, "%5D", "]")
	// s = strings.ReplaceAll(s, "%5C", "\\")
	// s = strings.ReplaceAll(s, "%5B", "[")
	// s = strings.ReplaceAll(s, "%21", "!")
	s = strings.ReplaceAll(s, "%21%5B", "![")
	s = strings.ReplaceAll(s, "%5D%28", "](")
	return strings.ReplaceAll(s, "%25", "%")
}
