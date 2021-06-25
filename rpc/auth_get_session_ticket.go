package rpc

import (
	"context"
	"crypto/md5"
	"log"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/rpc/message"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthGetSessionTicketResponse struct {
	Username string
	Uin      uint64
	Code     uint8

	PictureSign  []byte
	PictureData  []byte
	CaptchaSign  string
	ErrorCode    uint32
	ErrorTitle   string
	ErrorMessage string

	T104 []byte
	T119 []byte
	T150 []byte
	T161 []byte
	T401 [16]byte
	T402 []byte
	T403 []byte
	T546 []byte
}

func (resp *AuthGetSessionTicketResponse) Unmarshal(ctx context.Context, buf []byte) error {
	msg := &message.OICQMessage{
		RandomKey: clientRandomKey,
		PublicKey: ecdh.PublicKey,
		ShareKey:  ecdh.ShareKey,
	}
	if err := message.UnmarshalOICQMessage(ctx, buf, msg); err != nil {
		return err
	}
	resp.Username = strconv.Itoa(int(msg.Uin))
	resp.Uin = msg.Uin
	resp.Code = msg.Code
	if v, ok := msg.TLVs[0x0104].(*tlv.TLV); ok {
		resp.T104 = v.MustGetValue().Bytes()
	}
	if v, ok := msg.TLVs[0x0105].(*tlv.TLV); ok {
		buf, _ := v.GetValue()
		resp.PictureSign, _ = buf.DecodeBytes()
		resp.PictureData, _ = buf.DecodeBytes()
	}
	if v, ok := msg.TLVs[0x0119].(*tlv.TLV); ok {
		resp.T119 = v.MustGetValue().Bytes()
	}
	if v, ok := msg.TLVs[0x0192].(*tlv.TLV); ok {
		resp.CaptchaSign = string(v.MustGetValue().Bytes())
	}
	if v, ok := msg.TLVs[0x0146].(*tlv.TLV); ok {
		buf, _ := v.GetValue()
		resp.ErrorCode, _ = buf.DecodeUint32()
		resp.ErrorTitle, _ = buf.DecodeString()
		resp.ErrorMessage, _ = buf.DecodeString()
	}
	if v, ok := msg.TLVs[0x0150].(*tlv.TLV); ok {
		resp.T150 = v.MustGetValue().Bytes()
	}
	if v, ok := msg.TLVs[0x0161].(*tlv.TLV); ok {
		resp.T161 = v.MustGetValue().Bytes()
	}
	if v, ok := msg.TLVs[0x0402].(*tlv.TLV); ok {
		resp.T402 = v.MustGetValue().Bytes()
		resp.T401 = md5.Sum(append(append(deviceGUID[:], deviceDPWD...), resp.T402...))
	}
	if v, ok := msg.TLVs[0x0403].(*tlv.TLV); ok {
		resp.T403 = v.MustGetValue().Bytes()
	}
	if v, ok := msg.TLVs[0x0546].(*tlv.TLV); ok {
		resp.T546 = v.MustGetValue().Bytes()
	}
	return nil
}

func (c *Client) AuthGetSessionTicket(ctx context.Context, s2c *ServerToClientMessage) (*AuthGetSessionTicketResponse, error) {
	resp := new(AuthGetSessionTicketResponse)
	if err := resp.Unmarshal(ctx, s2c.Buffer); err != nil {
		return nil, err
	}
	switch resp.Code {
	case 0x00:
		// Success
		log.Printf("^_^ [info] login success, uin %s, code 0x00", resp.Username)
	case 0x02:
		// CAPTCHA
		c.t104 = resp.T104
		c.t547 = resp.T546 // TODO: check
		if resp.CaptchaSign != "" {
			log.Printf(">_x [warn] need captcha verify, uin %s, url %s, code 0x02", resp.Username, resp.CaptchaSign)
		} else {
			log.Printf(">_x [warn] need picture verify, uin %s, code 0x02", resp.Username)
		}
	case 0x01:
		log.Printf("x_x [fail] invalid login, uin %s, code 0x01, error %s: %s", resp.Username, resp.ErrorTitle, resp.ErrorMessage)
	case 0x06:
		log.Printf("x_x [fail] not implement, uin %s, code 0x06, error %s: %s", resp.Username, resp.ErrorTitle, resp.ErrorMessage)
	case 0x09:
		log.Printf("x_x [fail] invalid service, uin %s, code 0x09, error %s: %s", resp.Username, resp.ErrorTitle, resp.ErrorMessage)
	case 0xa0, 0xef:
		// DevLock
	case 0xed:
		log.Printf("x_x [fail] invalid device, uin %s, code 0x09, error %s: %s", resp.Username, resp.ErrorTitle, resp.ErrorMessage)
	case 0xcc:
		c.AuthRegisterDevice(ctx, NewAuthRegisterDeviceRequest(resp.Uin))
	}
	return resp, nil
}
