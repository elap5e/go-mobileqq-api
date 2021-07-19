package mark

import (
	"encoding/base64"
	"encoding/binary"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/encoding/mark/emoticon"
	"github.com/elap5e/go-mobileqq-api/pb"
)

type decoder struct {
	peerID, userID, fromID int64
}

func NewDecoder(peerID, userID, fromID int64) *decoder {
	return &decoder{peerID, userID, fromID}
}

func (dec decoder) Decode(v []byte) ([]*pb.Element, error) {
	// body := strings.SplitN(string(v), "\n", 2)[1]
	body := string(v)
	idxes := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`).FindAllStringSubmatchIndex(body, -1)
	elems := []*pb.Element{}
	offset := 0
	for _, idx := range idxes {
		if offset < idx[0] {
			elems = append(elems,
				dec.decodeText(body[offset:idx[0]])...)
		}
		offset = idx[1]
		uri, _ := url.Parse(body[idx[4]:idx[5]])
		if uri.Scheme != "goqq" {
			continue
		}
		switch uri.Hostname() {
		default:
			elems = append(elems,
				dec.decodeText(body[idx[2]:idx[3]])...)
		case "act":
			switch uri.Path {
			default:
				elems = append(elems,
					dec.decodeText(body[idx[2]:idx[3]])...)
			case "/at":
				elems = append(elems,
					dec.decodeActionAt(uri, body[idx[2]:idx[3]])...)
			case "/shakeWindow":
				elems = append(elems,
					dec.decodeActionShakeWindow()...)
			}
		case "res":
			switch uri.Path {
			default:
				elems = append(elems,
					dec.decodeText(body[idx[2]:idx[3]])...)
			case "/face":
				elems = append(elems,
					dec.decodeResourceFace(uri, body[idx[2]:idx[3]])...)
			case "/marketFace":
				elems = append(elems,
					dec.decodeResourceMarketFace(uri, body[idx[2]:idx[3]])...)
			case "/smallEmoji":
				elems = append(elems,
					dec.decodeResourceSmallEmoji(uri, body[idx[2]:idx[3]])...)
			case "/image":
				elems = append(elems,
					dec.decodeResourceImage(uri, body[idx[2]:idx[3]])...)
			}
		}
	}
	if offset < len(body) {
		elems = append(elems, dec.decodeText(body[offset:])...)
	}
	if len(elems) == 0 {
		elems = append(elems, dec.decodeText("[blank]")...)
	}
	return elems, nil
}

func (dec decoder) decodeActionAt(uri *url.URL, text string) []*pb.Element {
	text = unscape(text)
	attr6Buf := make([]byte, 13)
	attr6Buf[1] = 0x01
	binary.BigEndian.PutUint16(attr6Buf[4:], uint16(len([]rune(text))))
	uin, _ := strconv.ParseUint(uri.Query().Get("uin"), 10, 32)
	if uin == 0 {
		attr6Buf[6] = 0x01
	} else {
		binary.BigEndian.PutUint32(attr6Buf[7:], uint32(uin))
	}
	return []*pb.Element{{
		Text: &pb.Text{
			Data:        text,
			Attr6Buffer: attr6Buf,
		},
	}}
}

func (dec decoder) decodeActionShakeWindow() []*pb.Element {
	return []*pb.Element{{
		ShakeWindow: &pb.ShakeWindow{
			Uin: dec.peerID,
		},
	}}
}

func (dec decoder) decodeResourceFace(uri *url.URL, text string) []*pb.Element {
	if id := uri.Query().Get("id"); id != "" {
		tmp, _ := strconv.ParseInt(id, 10, 16)
		if tmp < 260 {
			return []*pb.Element{{
				Face: &pb.Face{
					Index: uint32(tmp),
				},
			}}
		} else {
			pid, _ := base64.URLEncoding.DecodeString(uri.Query().Get("pid"))
			sid, _ := base64.URLEncoding.DecodeString(uri.Query().Get("sid"))
			des := []byte(emoticon.FaceType(tmp).String())
			if len(pid)+len(sid) == 0 {
				buf, _ := proto.Marshal(&pb.MessageElementInfoServiceType33{
					Index:  uint32(tmp),
					Text:   des,
					Compat: des,
				})
				return []*pb.Element{{
					Common: &pb.CommonElement{
						ServiceType:  33,
						Buffer:       buf,
						BusinessType: 1,
					},
				}}
			} else {
				buf, _ := proto.Marshal(&pb.MessageElementInfoServiceType37{
					PackId:      pid,
					StickerId:   sid,
					QsId:        uint32(tmp),
					SourceType:  1,
					StickerType: 1,
					Text:        des,
				})
				return []*pb.Element{{
					Common: &pb.CommonElement{
						ServiceType:  37,
						Buffer:       buf,
						BusinessType: 1,
					},
				}, {
					Text: &pb.Text{
						Data: text,
					},
				}}
			}
		}
	}
	return []*pb.Element{{
		Text: &pb.Text{
			Data: "[face]",
		},
	}}
}

func (dec decoder) decodeResourceImage(uri *url.URL, text string) []*pb.Element {
	typ, _ := strconv.ParseUint(uri.Query().Get("type"), 10, 32)
	md5, _ := base64.URLEncoding.DecodeString(uri.Query().Get("md5"))
	size, _ := strconv.ParseUint(uri.Query().Get("size"), 10, 32)
	h, _ := strconv.ParseUint(uri.Query().Get("h"), 10, 32)
	w, _ := strconv.ParseUint(uri.Query().Get("w"), 10, 32)
	if dec.peerID == 0 {
		return []*pb.Element{{
			NotOnlineImage: &pb.NotOnlineImage{
				BizType:  uint32(typ),
				FileMd5:  md5,
				FileSize: uint32(size),
				Height:   uint32(h),
				Width:    uint32(w),
				FilePath: text,
				FileId:   uint32(rand.Intn(1e10)),
			},
		}}
	} else {
		return []*pb.Element{{
			CustomFace: &pb.CustomFace{
				BizType:  uint32(typ),
				FileMd5:  md5,
				FileSize: uint32(size),
				Height:   uint32(h),
				Width:    uint32(w),
				FilePath: text,
				FileId:   uint32(rand.Intn(1e10)),
				Useful:   1,
			},
		}}
	}
}

func (dec decoder) decodeResourceMarketFace(uri *url.URL, text string) []*pb.Element {
	id, _ := base64.URLEncoding.DecodeString(uri.Query().Get("id"))
	if len(id) != 0 {
		tabId, _ := strconv.ParseUint(uri.Query().Get("tabId"), 10, 32)
		key, _ := base64.URLEncoding.DecodeString(uri.Query().Get("key"))
		h, _ := strconv.ParseUint(uri.Query().Get("h"), 10, 32)
		w, _ := strconv.ParseUint(uri.Query().Get("w"), 10, 32)
		p, _ := base64.URLEncoding.DecodeString(uri.Query().Get("p"))
		return []*pb.Element{{
			MarketFace: &pb.MarketFace{
				FaceName:    []byte(text),
				FaceId:      id,
				TabId:       uint32(tabId),
				Key:         key,
				ImageHeight: uint32(h),
				ImageWidth:  uint32(w),
				ItemType:    6,
				FaceInfo:    1,
				SubType:     3,
				MobileParam: p,
			},
		}, {
			Text: &pb.Text{
				Data: text,
			},
		}}
	}
	return []*pb.Element{{
		Text: &pb.Text{
			Data: "[marketFace]",
		},
	}}
}

func (dec decoder) decodeResourceSmallEmoji(uri *url.URL, text string) []*pb.Element {
	id, _ := strconv.ParseUint(uri.Query().Get("id"), 10, 32)
	if id != 0 {
		typ, _ := strconv.ParseUint(uri.Query().Get("type"), 10, 32)
		return []*pb.Element{{
			SmallEmoji: &pb.SmallEmoji{
				PackIdSum: uint32(id),
				ImageType: uint32(typ),
			},
		}, {
			Text: &pb.Text{
				Data: text,
			},
		}}
	}
	return []*pb.Element{{
		Text: &pb.Text{
			Data: "[smallEmoji]",
		},
	}}
}

func (dec decoder) decodeText(text string) []*pb.Element {
	return []*pb.Element{{
		Text: &pb.Text{
			Data: unscape(text),
		},
	}}
}

func unscape(s string) string {
	// s = strings.ReplaceAll(s, "\\n", "\n")
	// s = strings.ReplaceAll(s, "%5D", "]")
	// s = strings.ReplaceAll(s, "%5C", "\\")
	// s = strings.ReplaceAll(s, "%5B", "[")
	// s = strings.ReplaceAll(s, "%21", "!")
	s = strings.ReplaceAll(s, "%21%5B", "![")
	s = strings.ReplaceAll(s, "%5D%28", "](")
	return strings.ReplaceAll(s, "%25", "%")
}
