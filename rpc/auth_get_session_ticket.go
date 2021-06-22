package rpc

import (
	"context"
	"crypto/md5"

	"github.com/elap5e/go-mobileqq-api/rpc/message"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthGetSessionTicketResponse struct {
	Username string
	Uin      uint64
	Code     uint8
	T104     []byte
	T119     []byte
	T150     []byte
	T161     []byte
	T401     [16]byte
	T402     []byte
	T403     []byte
}

func (resp *AuthGetSessionTicketResponse) Unmarshal(ctx context.Context, buf []byte) error {
	msg := &message.OICQMessage{
		RandomKey: defaultClientRandomKey,
		PublicKey: ecdh.PublicKey,
		ShareKey:  ecdh.ShareKey,
	}
	if err := message.UnmarshalOICQMessage(ctx, buf, msg); err != nil {
		return err
	}
	if v, ok := msg.TLVs[0x0104].(*tlv.TLV); ok {
		resp.T104 = v.MustGetValue().Bytes()
	}
	if v, ok := msg.TLVs[0x0119].(*tlv.TLV); ok {
		resp.T119 = v.MustGetValue().Bytes()
	}
	if v, ok := msg.TLVs[0x0150].(*tlv.TLV); ok {
		resp.T150 = v.MustGetValue().Bytes()
	}
	if v, ok := msg.TLVs[0x0161].(*tlv.TLV); ok {
		resp.T161 = v.MustGetValue().Bytes()
	}
	if v, ok := msg.TLVs[0x0402].(*tlv.TLV); ok {
		resp.T402 = v.MustGetValue().Bytes()
		resp.T401 = md5.Sum(append(append(defaultDeviceGUID, defaultDeviceDPWD...), resp.T402...))
	}
	if v, ok := msg.TLVs[0x0403].(*tlv.TLV); ok {
		resp.T403 = v.MustGetValue().Bytes()
	}
	return nil
}

func (c *Client) AuthGetSessionTicket(ctx context.Context, s2c *ServerToClientMessage) (interface{}, error) {
	resp := new(AuthGetSessionTicketResponse)
	if err := resp.Unmarshal(ctx, s2c.Buffer); err != nil {
		return nil, err
	}
	switch resp.Code {
	case 0x00:
	case 0xcc:
		return c.AuthRegisterDevice(ctx, NewAuthRegisterDeviceRequest(resp.Uin, resp.T104, resp.T401))
	}
	return resp.Code, nil
}
