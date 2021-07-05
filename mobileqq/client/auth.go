package client

import (
	"context"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthGetSessionTicketsRequest interface {
	GetTLVs(ctx context.Context) (map[uint16]tlv.TLVCodec, error)
	GetType() uint16
	SetType(typ uint16)
	GetSeq() uint32
	SetSeq(seq uint32)
	GetUin() uint64
	GetUsername() string
	SetUsername(username string)
	GetServiceMethod() string
}

type authGetSessionTicketsRequest struct {
	typ      uint16
	seq      uint32
	uin      uint64
	username string
	method   string
}

func (req *authGetSessionTicketsRequest) GetType() uint16 {
	return req.typ
}

func (req *authGetSessionTicketsRequest) SetType(typ uint16) {
	req.typ = typ
}

func (req *authGetSessionTicketsRequest) GetSeq() uint32 {
	return req.seq
}

func (req *authGetSessionTicketsRequest) SetSeq(seq uint32) {
	req.seq = seq
}

func (req *authGetSessionTicketsRequest) GetUin() uint64 {
	return req.uin
}

func (req *authGetSessionTicketsRequest) GetUsername() string {
	return req.username
}

func (req *authGetSessionTicketsRequest) SetUsername(username string) {
	req.username = username
	uin, _ := strconv.ParseInt(username, 10, 64)
	req.uin = uint64(uin)
}

func (req *authGetSessionTicketsRequest) GetServiceMethod() string {
	return req.method
}

func (req *authGetSessionTicketsRequest) SetServiceMethod(method string) {
	req.method = method
}

type AuthGetSessionTicketsResponse struct {
	Code     uint8
	Uin      uint64
	Username string

	PictureSign  []byte
	PictureData  []byte
	CaptchaSign  string
	ErrorCode    uint32
	ErrorTitle   string
	ErrorMessage string
	Message      string
	SMSMobile    string

	AuthSession []byte
	T119        []byte
	T150        []byte
	T161        []byte
	T174        []byte
	T17B        []byte
	T401        [16]byte
	T402        []byte
	T403        []byte
	T546        []byte

	LoginExtraData []byte
}

func (resp *AuthGetSessionTicketsResponse) SetTLVs(
	ctx context.Context,
	tlvs map[uint16]tlv.TLVCodec,
) error {
	if v, ok := tlvs[0x0104].(*tlv.TLV); ok {
		resp.AuthSession = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0105].(*tlv.TLV); ok {
		buf, _ := v.GetValue()
		resp.PictureSign, _ = buf.DecodeBytes()
		resp.PictureData, _ = buf.DecodeBytes()
	}
	if v, ok := tlvs[0x0119].(*tlv.TLV); ok {
		resp.T119 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0174].(*tlv.TLV); ok {
		resp.T174 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x017b].(*tlv.TLV); ok {
		resp.T17B = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0192].(*tlv.TLV); ok {
		resp.CaptchaSign = string(v.MustGetValue().Bytes())
	}
	if v, ok := tlvs[0x0146].(*tlv.TLV); ok {
		buf, _ := v.GetValue()
		resp.ErrorCode, _ = buf.DecodeUint32()
		resp.ErrorTitle, _ = buf.DecodeString()
		resp.ErrorMessage, _ = buf.DecodeString()
	}
	if v, ok := tlvs[0x0150].(*tlv.TLV); ok {
		resp.T150 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0161].(*tlv.TLV); ok {
		resp.T161 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x017e].(*tlv.TLV); ok {
		resp.Message = string(v.MustGetValue().Bytes())
	}
	if v, ok := tlvs[0x0178].(*tlv.TLV); ok {
		buf, _ := v.GetValue()
		_, _ = buf.DecodeString()
		mobile, _ := buf.DecodeString()
		resp.SMSMobile = mobile
	}
	if v, ok := tlvs[0x0402].(*tlv.TLV); ok {
		resp.T402 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0403].(*tlv.TLV); ok {
		resp.T403 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0537].(*tlv.TLV); ok {
		resp.LoginExtraData, _ = v.MustGetValue().DecodeBytes()
	}
	if v, ok := tlvs[0x0546].(*tlv.TLV); ok {
		resp.T546 = v.MustGetValue().Bytes()
	}
	return nil
}

type AuthContext struct {
}

var authCtxKey struct{}

func (c *Client) WithAuthContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, authCtxKey, c)
}

func ForAuthContext(ctx context.Context) *AuthContext {
	return ctx.Value(authCtxKey).(*AuthContext)
}
