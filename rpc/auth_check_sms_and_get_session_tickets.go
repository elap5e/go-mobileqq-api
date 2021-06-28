package rpc

import (
	"context"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/encoding/oicq"
	"github.com/elap5e/go-mobileqq-api/tlv"
)

type AuthCheckSMSAndGetSessionTicketsRequest struct {
	Seq    uint32
	Cookie []byte

	Uin      uint64
	Username string

	Session      []byte
	Code         []byte
	T174         []byte
	MiscBitmap   uint32
	SubSigMap    uint32
	SubAppIDList []uint64
	HashedGUID   [16]byte
	SMSExtraData []byte

	lockType uint8
}

func NewAuthCheckSMSAndGetSessionTicketsRequest(uin uint64, code []byte) *AuthCheckSMSAndGetSessionTicketsRequest {
	return &AuthCheckSMSAndGetSessionTicketsRequest{
		Uin:      uin,
		Username: fmt.Sprintf("%d", uin),

		Session:      nil,
		Code:         code,
		T174:         nil,
		MiscBitmap:   clientMiscBitmap,
		SubSigMap:    defaultClientSubSigMap,
		SubAppIDList: defaultClientSubAppIDList,
		HashedGUID:   [16]byte{},
		SMSExtraData: nil,

		lockType: 0x00,
	}
}

func (req *AuthCheckSMSAndGetSessionTicketsRequest) GetTLVs(ctx context.Context) (map[uint16]tlv.TLVCodec, error) {
	tlvs := make(map[uint16]tlv.TLVCodec)
	tlvs[0x0008] = tlv.NewT8(0x0000, defaultClientLocaleID, 0x0000)
	tlvs[0x0104] = tlv.NewT104(req.Session)
	tlvs[0x0116] = tlv.NewT116(req.MiscBitmap, req.SubSigMap, req.SubAppIDList)
	tlvs[0x0174] = tlv.NewT174(req.T174)
	tlvs[0x017c] = tlv.NewT17C(req.Code)
	tlvs[0x0401] = tlv.NewT401(req.HashedGUID)
	tlvs[0x0197] = tlv.NewTLV(0x0198, 0x0000, bytes.NewBuffer([]byte{req.lockType}))
	tlvs[0x0542] = tlv.NewT542(req.SMSExtraData)
	return tlvs, nil
}

func (c *Client) AuthCheckSMSAndGetSessionTickets(ctx context.Context, req *AuthCheckSMSAndGetSessionTicketsRequest) (*AuthGetSessionTicketsResponse, error) {
	req.Seq = c.getNextSeq()
	req.Cookie = c.cookie[:]
	req.Session = c.session
	req.T174 = c.t174
	req.HashedGUID = c.hashedGUID
	tlvs, err := req.GetTLVs(ctx)
	if err != nil {
		return nil, err
	}
	buf, err := oicq.Marshal(ctx, &oicq.Message{
		Version:       0x1f41,
		ServiceMethod: 0x0810,
		Uin:           req.Uin,
		EncryptMethod: oicq.EncryptMethodECDH,
		RandomKey:     c.randomKey,
		KeyVersion:    c.serverPublicKeyVersion,
		PublicKey:     c.privateKey.Public().Bytes(),
		ShareKey:      c.privateKey.ShareKey(c.serverPublicKey),
		Type:          0x0007,
		TLVs:          tlvs,
	})
	if err != nil {
		return nil, err
	}
	s2c := new(ServerToClientMessage)
	if err := c.Call(ServiceMethodAuthLogin, &ClientToServerMessage{
		Username: req.Username,
		Seq:      req.Seq,
		Buffer:   buf,
		Simple:   false,
	}, s2c); err != nil {
		return nil, err
	}
	return c.AuthGetSessionTickets(ctx, s2c)
}
