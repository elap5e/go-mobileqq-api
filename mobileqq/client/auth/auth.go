package auth

import (
	"context"

	"github.com/elap5e/go-mobileqq-api/tlv"
)

type Request interface {
	GetSeq() uint32
	SetSeq(seq uint32)
	GetServiceMethod() string
	SetServiceMethod(service string)
	GetType() uint16
	SetType(typ uint16)
	GetUin() uint64
	SetUin(uin uint64)
	GetUsername() string
	SetUsername(username string)

	MustGetTLVs(ctx context.Context) map[uint16]tlv.TLVCodec
}

type request struct {
	seq      uint32
	service  string
	typ      uint16
	uin      uint64
	username string
}

func (req *request) GetSeq() uint32 { return req.seq }

func (req *request) SetSeq(seq uint32) { req.seq = seq }

func (req *request) GetServiceMethod() string { return req.service }

func (req *request) SetServiceMethod(service string) { req.service = service }

func (req *request) GetType() uint16 { return req.typ }

func (req *request) SetType(typ uint16) { req.typ = typ }

func (req *request) GetUin() uint64 { return req.uin }

func (req *request) SetUin(uin uint64) { req.uin = uin }

func (req *request) GetUsername() string { return req.username }

func (req *request) SetUsername(username string) { req.username = username }

type Response struct {
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

func (resp *Response) SetTLVs(
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
