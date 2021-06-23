package rpc

import (
	"context"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/rpc/message"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthCheckSMSRequest struct {
	Seq      uint32
	Uin      uint64
	Username string
	LockType uint8

	T104         []byte
	T17C         []byte
	T174         []byte
	MiscBitmap   uint32
	SubSigMap    uint32
	SubAppIDList []uint64
	T401         [16]byte
	T547         []byte
	SMSExtraData []byte
}

func NewAuthCheckSMSRequest(uin uint64) *AuthCheckSMSRequest {
	return &AuthCheckSMSRequest{
		Uin:      uin,
		Username: fmt.Sprintf("%d", uin),
		LockType: 0x00,

		T104:         nil,
		MiscBitmap:   defaultClientMiscBitmap,
		SubSigMap:    defaultClientSubSigMap,
		SubAppIDList: defaultClientSubAppIDList,
		T547:         nil,
	}
}

func (req *AuthCheckSMSRequest) EncodeOICQMessage(ctx context.Context) (*message.OICQMessage, error) {
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	tlvs[0x0104] = tlv.NewT104(req.T104)
	tlvs[0x0116] = tlv.NewT116(req.MiscBitmap, req.SubSigMap, req.SubAppIDList)
	tlvs[0x0174] = tlv.NewT174(req.T174)
	tlvs[0x017c] = tlv.NewT17C(req.T17C)
	tlvs[0x0401] = tlv.NewT401(req.T401)
	tlvs[0x0197] = tlv.NewTLV(0x0198, 0x0000, bytes.NewBuffer([]byte{req.LockType}))
	tlvs[0x0542] = tlv.NewT542(req.SMSExtraData)

	return &message.OICQMessage{
		Version:       0x1f41,
		ServiceMethod: 0x0810,
		Uin:           req.Uin,
		EncryptMethod: 0x07,
		RandomKey:     defaultClientRandomKey,
		PublicKey:     ecdh.PublicKey,
		ShareKey:      ecdh.ShareKey,
		Type:          0x0007,
		TLVs:          tlvs,
	}, nil
}

func (req *AuthCheckSMSRequest) Encode(ctx context.Context) (*ClientToServerMessage, error) {
	msg, err := req.EncodeOICQMessage(ctx)
	if err != nil {
		return nil, err
	}
	buf, err := message.MarshalOICQMessage(ctx, msg)
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

func (c *Client) AuthCheckSMS(ctx context.Context, req *AuthCheckSMSRequest) (interface{}, error) {
	req.Seq = c.getNextSeq()
	req.T104 = []byte{}
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
