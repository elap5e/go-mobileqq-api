package rpc

import (
	"context"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthCheckSMSAndGetSessionTicketRequest struct {
	Seq      uint32
	Uin      uint64
	Username string
	Cookie   []byte
	LockType uint8

	T104         []byte
	Code         []byte
	T174         []byte
	MiscBitmap   uint32
	SubSigMap    uint32
	SubAppIDList []uint64
	T401         [16]byte
	SMSExtraData []byte
}

func NewAuthCheckSMSAndGetSessionTicketRequest(uin uint64, code []byte) *AuthCheckSMSAndGetSessionTicketRequest {
	return &AuthCheckSMSAndGetSessionTicketRequest{
		Uin:      uin,
		Username: fmt.Sprintf("%d", uin),
		LockType: 0x00,

		T104:         nil,
		Code:         code,
		T174:         nil,
		MiscBitmap:   clientMiscBitmap,
		SubSigMap:    defaultClientSubSigMap,
		SubAppIDList: defaultClientSubAppIDList,
		T401:         [16]byte{},
		SMSExtraData: nil,
	}
}

func (req *AuthCheckSMSAndGetSessionTicketRequest) EncodeOICQMessage(ctx context.Context) (*oicq.Message, error) {
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	tlvs[0x0104] = tlv.NewT104(req.T104)
	tlvs[0x0116] = tlv.NewT116(req.MiscBitmap, req.SubSigMap, req.SubAppIDList)
	tlvs[0x0174] = tlv.NewT174(req.T174)
	tlvs[0x017c] = tlv.NewT17C(req.Code)
	tlvs[0x0401] = tlv.NewT401(req.T401)
	tlvs[0x0197] = tlv.NewTLV(0x0198, 0x0000, bytes.NewBuffer([]byte{req.LockType}))
	tlvs[0x0542] = tlv.NewT542(req.SMSExtraData)

	return &oicq.Message{
		Version:       0x1f41,
		ServiceMethod: 0x0810,
		Uin:           req.Uin,
		EncryptMethod: 0x87,
		RandomKey:     clientRandomKey,
		PublicKey:     ecdh.PublicKey,
		ShareKey:      ecdh.ShareKey,
		Type:          0x0007,
		TLVs:          tlvs,
	}, nil
}

func (req *AuthCheckSMSAndGetSessionTicketRequest) Encode(ctx context.Context) (*ClientToServerMessage, error) {
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
		AppID:    clientAppID,
		Cookie:   req.Cookie,
		Buffer:   buf,
		Simple:   false,
	}, nil
}

func (c *Client) AuthCheckSMSAndGetSessionTicket(ctx context.Context, req *AuthCheckSMSAndGetSessionTicketRequest) (*AuthGetSessionTicketResponse, error) {
	req.Seq = c.getNextSeq()
	req.Cookie = c.cookie[:]
	req.T104 = c.t104
	req.T174 = c.t174
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
