package rpc

import (
	"context"
	"crypto/md5"
	"log"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/crypto"
	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthGetSessionTicketResponse struct {
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

	T104 []byte
	T119 []byte
	T150 []byte
	T161 []byte
	T174 []byte
	T17B []byte
	T401 [16]byte
	T402 []byte
	T403 []byte
	T546 []byte
}

func (resp *AuthGetSessionTicketResponse) SetTLVs(ctx context.Context, tlvs map[uint16]tlv.TLVCodec) error {
	if v, ok := tlvs[0x0104].(*tlv.TLV); ok {
		resp.T104 = v.MustGetValue().Bytes()
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
	if v, ok := tlvs[0x0546].(*tlv.TLV); ok {
		resp.T546 = v.MustGetValue().Bytes()
	}
	return nil
}

func (c *Client) AuthGetSessionTicket(ctx context.Context, s2c *ServerToClientMessage) (*AuthGetSessionTicketResponse, error) {
	resp := new(AuthGetSessionTicketResponse)
	msg := &oicq.Message{
		RandomKey:  c.randomKey,
		KeyVersion: c.serverPublicKeyVersion,
		PublicKey:  c.privateKey.Public().Bytes(),
		ShareKey:   c.privateKey.ShareKey(c.serverPublicKey),
	}
	if err := oicq.Unmarshal(ctx, s2c.Buffer, msg); err != nil {
		return nil, err
	}
	resp.Code = msg.Code
	resp.Uin = msg.Uin
	resp.Username = strconv.Itoa(int(msg.Uin))
	if err := resp.SetTLVs(ctx, msg.TLVs); err != nil {
		return nil, err
	}
	switch resp.Code {
	case 0x00:
		// success
		c.t119 = crypto.NewCipher(c.tgtgtKey).Decrypt(resp.T119)
		buf := bytes.NewBuffer(c.t119)
		l, _ := buf.DecodeUint16()
		tlvs := map[uint16]tlv.TLVCodec{}
		for i := 0; i < int(l); i++ {
			tlv := tlv.TLV{}
			tlv.Decode(buf)
			tlvs[tlv.GetType()] = &tlv
		}
		tlv.DumpTLVs(ctx, tlvs)
		log.Printf("^_^ [info] login success, uin %s, code 0x00", resp.Username)
	case 0x02:
		// captcha
		c.t104 = resp.T104
		c.t547 = resp.T546 // TODO: check
		if resp.CaptchaSign != "" {
			log.Printf(">_x [warn] need captcha verify, uin %s, url %s, code 0x02", resp.Username, resp.CaptchaSign)
		} else {
			log.Printf(">_x [warn] need picture verify, uin %s, code 0x02", resp.Username)
		}
	case 0xa0:
		// device lock
		c.t104 = resp.T104
		c.t17b = resp.T17B
		log.Printf(">_x [warn] need sms mobile verify response, uin %s, code 0xa0", resp.Username)
	case 0xef:
		// device lock
		c.t104 = resp.T104
		c.t174 = resp.T174
		c.t402 = resp.T402
		c.t403 = resp.T403
		c.t401 = md5.Sum(append(append(deviceGUID[:], deviceDPWD...), c.t402...))
		log.Printf(">_x [warn] need sms mobile verify, uin %s, mobile %s, code 0x%02x, message %s, code 0xef", resp.Username, resp.SMSMobile, resp.Code, resp.Message)
	case 0x01:
		log.Printf("x_x [fail] invalid login, uin %s, code 0x01, error %s: %s", resp.Username, resp.ErrorTitle, resp.ErrorMessage)
	case 0x06:
		log.Printf("x_x [fail] not implement, uin %s, code 0x06, error %s: %s", resp.Username, resp.ErrorTitle, resp.ErrorMessage)
	case 0x09:
		log.Printf("x_x [fail] invalid service, uin %s, code 0x09, error %s: %s", resp.Username, resp.ErrorTitle, resp.ErrorMessage)
	case 0xed:
		log.Printf("x_x [fail] invalid device, uin %s, code 0xed, error %s: %s", resp.Username, resp.ErrorTitle, resp.ErrorMessage)
	case 0xcc:
		c.t104 = resp.T104
		c.t402 = resp.T402
		c.t403 = resp.T403
		c.t401 = md5.Sum(append(append(deviceGUID[:], deviceDPWD...), c.t402...))
		return c.AuthUnlockDevice(ctx, NewAuthUnlockDeviceRequest(resp.Uin))
	}
	return resp, nil
}
