package rpc

import (
	"context"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthCheckCaptchaAndGetSessionTicketRequest struct {
	Seq      uint32
	Uin      uint64
	Username string
	Cookie   []byte
	CheckWeb bool

	T104         []byte
	Code         []byte
	Sign         []byte
	MiscBitmap   uint32
	SubSigMap    uint32
	SubAppIDList []uint64
	T547         []byte
}

func NewAuthCheckCaptchaAndGetSessionTicketRequest(uin uint64, code []byte) *AuthCheckCaptchaAndGetSessionTicketRequest {
	return &AuthCheckCaptchaAndGetSessionTicketRequest{
		Uin:      uin,
		Username: fmt.Sprintf("%d", uin),
		CheckWeb: true,

		T104:         nil,
		Code:         code,
		Sign:         nil,
		MiscBitmap:   clientMiscBitmap,
		SubSigMap:    defaultClientSubSigMap,
		SubAppIDList: defaultClientSubAppIDList,
		T547:         nil,
	}
}

func (req *AuthCheckCaptchaAndGetSessionTicketRequest) EncodeOICQMessage(ctx context.Context) (*oicq.Message, error) {
	tlvs := make(map[uint16]tlv.TLVCodec)
	if req.CheckWeb {
		tlvs[0x0193] = tlv.NewT193(req.Code)
	} else {
		tlvs[0x0002] = tlv.NewT2(req.Code, req.Sign)
	}
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	tlvs[0x0104] = tlv.NewT104(req.T104)
	tlvs[0x0116] = tlv.NewT116(req.MiscBitmap, req.SubSigMap, req.SubAppIDList)
	tlvs[0x0547] = tlv.NewT547(req.T547)

	return &oicq.Message{
		Version:       0x1f41,
		ServiceMethod: 0x0810,
		Uin:           req.Uin,
		EncryptMethod: 0x87,
		RandomKey:     clientRandomKey,
		KeyVersion:    ecdh.KeyVersion,
		PublicKey:     ecdh.PublicKey,
		ShareKey:      ecdh.ShareKey,
		Type:          0x0002,
		TLVs:          tlvs,
	}, nil
}

func (req *AuthCheckCaptchaAndGetSessionTicketRequest) Encode(ctx context.Context) (*ClientToServerMessage, error) {
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

func (c *Client) AuthCheckCaptchaAndGetSessionTicket(ctx context.Context, req *AuthCheckCaptchaAndGetSessionTicketRequest) (*AuthGetSessionTicketResponse, error) {
	req.Seq = c.getNextSeq()
	req.Cookie = c.cookie[:]
	req.T104 = c.t104
	req.T547 = c.t547
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
