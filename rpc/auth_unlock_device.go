package rpc

import (
	"context"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthUnlockDeviceRequest struct {
	Seq      uint32
	Uin      uint64
	Username string

	T104         []byte
	MiscBitmap   uint32
	SubSigMap    uint32
	SubAppIDList []uint64
	T401         [16]byte
}

func NewAuthUnlockDeviceRequest(uin uint64) *AuthUnlockDeviceRequest {
	return &AuthUnlockDeviceRequest{
		Uin:      uin,
		Username: fmt.Sprintf("%d", uin),

		T104:         nil,
		MiscBitmap:   clientMiscBitmap,
		SubSigMap:    defaultClientSubSigMap,
		SubAppIDList: defaultClientSubAppIDList,
		T401:         [16]byte{},
	}
}

func (req *AuthUnlockDeviceRequest) EncodeOICQMessage(ctx context.Context) (*oicq.Message, error) {
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	tlvs[0x0104] = tlv.NewT104(req.T104)
	tlvs[0x0116] = tlv.NewT116(req.MiscBitmap, req.SubSigMap, req.SubAppIDList)
	tlvs[0x0401] = tlv.NewT401(req.T401)

	return &oicq.Message{
		Version:       0x1f41,
		ServiceMethod: 0x0810,
		Uin:           req.Uin,
		EncryptMethod: 0x87,
		RandomKey:     clientRandomKey,
		KeyVersion:    ecdh.KeyVersion,
		PublicKey:     ecdh.PublicKey,
		ShareKey:      ecdh.ShareKey,
		Type:          0x0014,
		TLVs:          tlvs,
	}, nil
}

func (req *AuthUnlockDeviceRequest) Encode(ctx context.Context) (*ClientToServerMessage, error) {
	msg, err := req.EncodeOICQMessage(ctx)
	if err != nil {
		return nil, err
	}
	buf, err := oicq.Marshal(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &ClientToServerMessage{
		Username: req.Username,
		Seq:      req.Seq,
		Buffer:   buf,
		Simple:   false,
	}, nil
}

func (c *Client) AuthUnlockDevice(ctx context.Context, req *AuthUnlockDeviceRequest) (*AuthGetSessionTicketResponse, error) {
	req.Seq = c.getNextSeq()
	req.T104 = c.t104
	req.T401 = c.t401
	c2s, err := req.Encode(ctx)
	if err != nil {
		return nil, err
	}
	s2c := new(ServerToClientMessage)
	if err := c.Call("wtlogin.login", c2s, s2c); err != nil {
		return nil, err
	}
	return c.AuthGetSessionTicket(ctx, s2c)
}
